package webserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	config "vServer/Backend/config"
	tools "vServer/Backend/tools"
)

var (
	phpProcesses []*exec.Cmd
	fcgiPorts    []int
	portIndex    int
	portMutex    sync.Mutex
	maxWorkers   = 4
	stopping     = false // Флаг остановки
)

var address_php string
var Сonsole_php bool = false

// GetPHPStatus возвращает статус PHP сервера
func GetPHPStatus() bool {
	return len(phpProcesses) > 0 && !stopping
}

// FastCGI константы
const (
	FCGI_VERSION_1         = 1
	FCGI_BEGIN_REQUEST     = 1
	FCGI_ABORT_REQUEST     = 2
	FCGI_END_REQUEST       = 3
	FCGI_PARAMS            = 4
	FCGI_STDIN             = 5
	FCGI_STDOUT            = 6
	FCGI_STDERR            = 7
	FCGI_DATA              = 8
	FCGI_GET_VALUES        = 9
	FCGI_GET_VALUES_RESULT = 10
	FCGI_UNKNOWN_TYPE      = 11
	FCGI_MAXTYPE           = FCGI_UNKNOWN_TYPE

	FCGI_NULL_REQUEST_ID = 0

	FCGI_KEEP_CONN = 1

	FCGI_RESPONDER  = 1
	FCGI_AUTHORIZER = 2
	FCGI_FILTER     = 3
)

// FastCGI заголовок
type FCGIHeader struct {
	Version       byte
	Type          byte
	RequestID     uint16
	ContentLength uint16
	PaddingLength byte
	Reserved      byte
}

// FastCGI BeginRequest body
type FCGIBeginRequestBody struct {
	Role     uint16
	Flags    byte
	Reserved [5]byte
}

func PHP_Start() {
	// Сбрасываем флаг остановки
	stopping = false

	// Читаем настройки из конфига
	address_php = config.ConfigData.Soft_Settings.Php_host

	// Запускаем FastCGI процессы
	for i := 0; i < maxWorkers; i++ {
		port := config.ConfigData.Soft_Settings.Php_port + i
		fcgiPorts = append(fcgiPorts, port)
		go startFastCGIWorker(port, i)
		time.Sleep(200 * time.Millisecond) // Задержка между запусками
	}

	tools.Logs_file(0, "PHP  ", fmt.Sprintf("💻 PHP FastCGI пул запущен (%d процессов на портах %d-%d)", maxWorkers, config.ConfigData.Soft_Settings.Php_port, config.ConfigData.Soft_Settings.Php_port+maxWorkers-1), "logs_php.log", true)
}

func startFastCGIWorker(port int, workerID int) {
	phpPath := "WebServer/soft/PHP/php_v_8/php-cgi.exe"

	cmd := exec.Command(phpPath, "-b", fmt.Sprintf("%s:%d", address_php, port))
	cmd.Env = append(os.Environ(),
		"PHP_FCGI_CHILDREN=0",        // Один процесс на порт
		"PHP_FCGI_MAX_REQUESTS=1000", // Перезапуск после 1000 запросов
	)

	// Скрываем консольное окно
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

	if !Сonsole_php {
		cmd.Stdout = nil
		cmd.Stderr = nil
	}

	err := cmd.Start()
	if err != nil {
		tools.Logs_file(1, "PHP", fmt.Sprintf("❌ Ошибка запуска FastCGI worker %d на порту %d: %v", workerID, port, err), "logs_php.log", true)
		return
	}

	phpProcesses = append(phpProcesses, cmd)
	tools.Logs_file(0, "PHP", fmt.Sprintf("✅ PHP FastCGI %d запущен на %s:%d", workerID, address_php, port), "logs_php.log", false)

	// Ждём завершения процесса и перезапускаем
	go func() {
		cmd.Wait()

		// Проверяем, не останавливается ли сервер
		if stopping {
			return // Не перезапускаем если сервер останавливается
		}

		tools.Logs_file(1, "PHP", fmt.Sprintf("⚠️ FastCGI worker %d завершился, перезапускаем...", workerID), "logs_php.log", true)
		time.Sleep(1 * time.Second)
		startFastCGIWorker(port, workerID) // Перезапуск
	}()
}

// Получение следующего порта из пула (round-robin)
func getNextFCGIPort() int {
	portMutex.Lock()
	defer portMutex.Unlock()

	port := fcgiPorts[portIndex]
	portIndex = (portIndex + 1) % len(fcgiPorts)
	return port
}

// Создание FastCGI пакета
func createFCGIPacket(requestType byte, requestID uint16, content []byte) []byte {
	contentLength := len(content)
	paddingLength := 8 - (contentLength % 8)
	if paddingLength == 8 {
		paddingLength = 0
	}

	header := FCGIHeader{
		Version:       FCGI_VERSION_1,
		Type:          requestType,
		RequestID:     requestID,
		ContentLength: uint16(contentLength),
		PaddingLength: byte(paddingLength),
		Reserved:      0,
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, header)
	buf.Write(content)
	buf.Write(make([]byte, paddingLength)) // Padding

	return buf.Bytes()
}

// Кодирование FastCGI параметров
func encodeFCGIParams(params map[string]string) []byte {
	var buf bytes.Buffer

	for key, value := range params {
		keyLen := len(key)
		valueLen := len(value)

		// Длина ключа
		if keyLen < 128 {
			buf.WriteByte(byte(keyLen))
		} else {
			binary.Write(&buf, binary.BigEndian, uint32(keyLen)|0x80000000)
		}

		// Длина значения
		if valueLen < 128 {
			buf.WriteByte(byte(valueLen))
		} else {
			binary.Write(&buf, binary.BigEndian, uint32(valueLen)|0x80000000)
		}

		// Ключ и значение
		buf.WriteString(key)
		buf.WriteString(value)
	}

	return buf.Bytes()
}

// HandlePHPRequest - универсальная функция для обработки файлов
// Проверяет является ли файл PHP и обрабатывает соответственно
// Возвращает true если файл был обработан (PHP или статический), false если нужна обработка ошибки
func HandlePHPRequest(w http.ResponseWriter, r *http.Request, host string, filePath string, originalURI string, originalPath string) bool {
	// Импортируем path/filepath для проверки расширения
	if filepath.Ext(filePath) == ".php" {
		// Сохраняем оригинальные значения URL
		originalURL := r.URL.Path
		originalRawQuery := r.URL.RawQuery

		// Устанавливаем путь к PHP файлу
		r.URL.Path = filePath

		// Вызываем существующий PHPHandler
		PHPHandler(w, r, host, originalURI, originalPath)

		// Восстанавливаем оригинальные значения
		r.URL.Path = originalURL
		r.URL.RawQuery = originalRawQuery
		return true
	} else {
		// Это не PHP файл - обрабатываем как статический
		fullPath := "WebServer/www/" + host + "/public_www" + filePath
		http.ServeFile(w, r, fullPath)
		return true
	}
}

// PHPHandler с FastCGI
func PHPHandler(w http.ResponseWriter, r *http.Request, host string, originalURI string, originalPath string) {
	phpPath := "WebServer/www/" + host + "/public_www" + r.URL.Path

	// Проверяем существование файла
	if _, err := os.Stat(phpPath); os.IsNotExist(err) {
		http.ServeFile(w, r, "WebServer/tools/error_page/index.html")
		tools.Logs_file(2, "PHP_404", "🔍 PHP файл не найден: "+phpPath, "logs_php.log", false)
		return
	}

	// Получаем абсолютный путь для SCRIPT_FILENAME
	absPath, err := filepath.Abs(phpPath)
	if err != nil {
		tools.Logs_file(1, "PHP", "❌ Ошибка получения абсолютного пути: "+err.Error(), "logs_php.log", false)
		absPath = phpPath
	}

	// Получаем порт FastCGI
	port := getNextFCGIPort()

	// Подключаемся к FastCGI процессу
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", address_php, port), 5*time.Second)
	if err != nil {
		tools.Logs_file(1, "PHP", fmt.Sprintf("❌ Ошибка подключения к FastCGI порт %d: %v", port, err), "logs_php.log", false)
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}
	defer conn.Close()

	// Читаем POST данные
	var postData []byte
	if r.Method == "POST" {
		postData, _ = io.ReadAll(r.Body)
		r.Body.Close()
	}

	// Формируем параметры FastCGI
	serverPort := "80"
	if r.TLS != nil {
		serverPort = "443"
	}

	// Используем переданные оригинальные значения или текущие если не переданы
	requestURI := r.URL.RequestURI()
	if originalURI != "" {
		requestURI = originalURI
	}

	pathInfo := r.URL.Path
	if originalPath != "" {
		pathInfo = originalPath
	}

	params := map[string]string{
		"REQUEST_METHOD":    r.Method,
		"REQUEST_URI":       requestURI,
		"QUERY_STRING":      r.URL.RawQuery,
		"CONTENT_TYPE":      r.Header.Get("Content-Type"),
		"CONTENT_LENGTH":    fmt.Sprintf("%d", len(postData)),
		"SCRIPT_FILENAME":   absPath,
		"SCRIPT_NAME":       r.URL.Path,
		"DOCUMENT_ROOT":     "WebServer/www/" + host + "/public_www",
		"SERVER_NAME":       host,
		"HTTP_HOST":         host,
		"SERVER_PORT":       serverPort,
		"SERVER_PROTOCOL":   "HTTP/1.1",
		"GATEWAY_INTERFACE": "CGI/1.1",
		"REDIRECT_STATUS":   "200",
		"REMOTE_ADDR":       strings.Split(r.RemoteAddr, ":")[0],
		"REMOTE_HOST":       strings.Split(r.RemoteAddr, ":")[0],
		"PATH_INFO":         pathInfo,
		"PATH_TRANSLATED":   absPath,
	}

	// Добавляем HTTP заголовки
	for name, values := range r.Header {
		if len(values) > 0 {
			httpName := "HTTP_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
			params[httpName] = values[0]
		}
	}

	requestID := uint16(1)

	// 1. Отправляем BEGIN_REQUEST
	beginRequest := FCGIBeginRequestBody{
		Role:  FCGI_RESPONDER,
		Flags: 0,
	}
	var beginBuf bytes.Buffer
	binary.Write(&beginBuf, binary.BigEndian, beginRequest)
	packet := createFCGIPacket(FCGI_BEGIN_REQUEST, requestID, beginBuf.Bytes())
	conn.Write(packet)

	// 2. Отправляем PARAMS с разбивкой на чанки
	paramsData := encodeFCGIParams(params)
	if len(paramsData) > 0 {
		const maxChunkSize = 65535 // Максимальный размер FastCGI пакета

		for offset := 0; offset < len(paramsData); offset += maxChunkSize {
			end := offset + maxChunkSize
			if end > len(paramsData) {
				end = len(paramsData)
			}

			chunk := paramsData[offset:end]
			packet = createFCGIPacket(FCGI_PARAMS, requestID, chunk)
			conn.Write(packet)
		}
	}

	// 3. Пустой PARAMS (конец параметров)
	packet = createFCGIPacket(FCGI_PARAMS, requestID, []byte{})
	conn.Write(packet)

	// 4. Отправляем STDIN (POST данные) с разбивкой на чанки
	if len(postData) > 0 {
		const maxChunkSize = 65535 // Максимальный размер FastCGI пакета

		for offset := 0; offset < len(postData); offset += maxChunkSize {
			end := offset + maxChunkSize
			if end > len(postData) {
				end = len(postData)
			}

			chunk := postData[offset:end]
			packet = createFCGIPacket(FCGI_STDIN, requestID, chunk)
			conn.Write(packet)
		}
	}

	// 5. Пустой STDIN (конец данных)
	packet = createFCGIPacket(FCGI_STDIN, requestID, []byte{})
	conn.Write(packet)

	// Читаем и стримим ответ (с поддержкой SSE и chunked transfer)
	err = streamFastCGIResponse(conn, requestID, w)
	if err != nil {
		tools.Logs_file(1, "PHP", "❌ Ошибка чтения FastCGI ответа: "+err.Error(), "logs_php.log", false)
		// Не вызываем http.Error здесь, т.к. заголовки уже могли быть отправлены
		return
	}

}

// Streaming чтение FastCGI ответа с поддержкой SSE и chunked transfer
func streamFastCGIResponse(conn net.Conn, requestID uint16, w http.ResponseWriter) error {
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	var stderr bytes.Buffer
	var headerBuffer bytes.Buffer
	headersWritten := false
	flusher, canFlush := w.(http.Flusher)

	for {
		// Читаем заголовок FastCGI
		headerBuf := make([]byte, 8)
		_, err := io.ReadFull(conn, headerBuf)
		if err != nil {
			return err
		}

		var header FCGIHeader
		buf := bytes.NewReader(headerBuf)
		binary.Read(buf, binary.BigEndian, &header)

		// Читаем содержимое
		content := make([]byte, header.ContentLength)
		if header.ContentLength > 0 {
			_, err = io.ReadFull(conn, content)
			if err != nil {
				return err
			}
		}

		// Читаем padding
		if header.PaddingLength > 0 {
			padding := make([]byte, header.PaddingLength)
			io.ReadFull(conn, padding)
		}

		// Обрабатываем пакет
		switch header.Type {
		case FCGI_STDOUT:
			if header.ContentLength > 0 {
				if !headersWritten {
					// Накапливаем данные до разделителя заголовков
					headerBuffer.Write(content)

					// Ищем разделитель между заголовками и телом
					headerStr := headerBuffer.String()
					sepIndex := strings.Index(headerStr, "\r\n\r\n")
					if sepIndex == -1 {
						sepIndex = strings.Index(headerStr, "\n\n")
					}

					if sepIndex != -1 {
						// Нашли разделитель - обрабатываем заголовки
						var sepLen int
						if strings.Contains(headerStr[:sepIndex+4], "\r\n\r\n") {
							sepLen = 4
						} else {
							sepLen = 2
						}

						headersPart := headerStr[:sepIndex]
						bodyPart := headerStr[sepIndex+sepLen:]

						// Парсим и устанавливаем заголовки
						processStreamingHeaders(w, headersPart)
						headersWritten = true

						// Отправляем первую часть тела
						if len(bodyPart) > 0 {
							w.Write([]byte(bodyPart))
							if canFlush {
								flusher.Flush()
							}
						}
					}
				} else {
					// Заголовки уже отправлены - стримим тело
					w.Write(content)
					// Принудительно отправляем данные (критично для SSE)
					if canFlush {
						flusher.Flush()
					}
				}
			} else {
				// Пустой STDOUT - конец данных, если остались заголовки без тела
				if !headersWritten && headerBuffer.Len() > 0 {
					processStreamingHeaders(w, headerBuffer.String())
					headersWritten = true
				}
			}
		case FCGI_STDERR:
			if header.ContentLength > 0 {
				stderr.Write(content)
			}
		case FCGI_END_REQUEST:
			// Завершение запроса
			if stderr.Len() > 0 {
				tools.Logs_file(1, "PHP", "FastCGI stderr: "+stderr.String(), "logs_php.log", false)
			}
			// Если заголовки так и не были записаны (пустой ответ)
			if !headersWritten {
				w.WriteHeader(http.StatusOK)
			}
			return nil
		}
	}
}

// Обработка заголовков для streaming ответа
func processStreamingHeaders(w http.ResponseWriter, headersPart string) {
	headers := strings.Split(headersPart, "\n")
	statusCode := 200

	for _, header := range headers {
		header = strings.TrimSpace(header)
		if header == "" {
			continue
		}

		if strings.HasPrefix(strings.ToLower(header), "content-type:") {
			contentType := strings.TrimSpace(strings.SplitN(header, ":", 2)[1])
			w.Header().Set("Content-Type", contentType)
		} else if strings.HasPrefix(strings.ToLower(header), "set-cookie:") {
			cookie := strings.TrimSpace(strings.SplitN(header, ":", 2)[1])
			w.Header().Add("Set-Cookie", cookie)
		} else if strings.HasPrefix(strings.ToLower(header), "location:") {
			location := strings.TrimSpace(strings.SplitN(header, ":", 2)[1])
			w.Header().Set("Location", location)
			statusCode = http.StatusFound
		} else if strings.HasPrefix(strings.ToLower(header), "status:") {
			status := strings.TrimSpace(strings.SplitN(header, ":", 2)[1])
			if code, err := strconv.Atoi(strings.Split(status, " ")[0]); err == nil {
				statusCode = code
			}
		} else if strings.Contains(header, ":") {
			headerParts := strings.SplitN(header, ":", 2)
			if len(headerParts) == 2 {
				w.Header().Set(strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1]))
			}
		}
	}

	w.WriteHeader(statusCode)
}

// PHP_Stop останавливает все FastCGI процессы
func PHP_Stop() {
	// Устанавливаем флаг остановки
	stopping = true

	for i, cmd := range phpProcesses {
		if cmd != nil && cmd.Process != nil {
			err := cmd.Process.Kill()
			if err != nil {
				tools.Logs_file(1, "PHP", fmt.Sprintf("❌ Ошибка остановки FastCGI процесса %d: %v", i, err), "logs_php.log", true)
			} else {
				tools.Logs_file(0, "PHP", fmt.Sprintf("✅ FastCGI процесс %d остановлен", i), "logs_php.log", false)
			}
		}
	}

	phpProcesses = nil
	fcgiPorts = nil

	// Дополнительно убиваем все процессы php-cgi.exe
	cmd := exec.Command("taskkill", "/F", "/IM", "php-cgi.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	cmd.Run()

	tools.Logs_file(0, "PHP", "🛑 Все FastCGI процессы остановлены", "logs_php.log", true)
}
