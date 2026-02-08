package webserver

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"vServer/Backend/config"
	tools "vServer/Backend/tools"
)

var (
	configMutex sync.RWMutex
)

// Проверяет, является ли запрос WebSocket
func isWebSocketRequest(r *http.Request) bool {
	connection := strings.ToLower(r.Header.Get("Connection"))
	upgrade := strings.ToLower(r.Header.Get("Upgrade"))
	return strings.Contains(connection, "upgrade") && upgrade == "websocket"
}

// Проксирует WebSocket соединение
func handleWebSocketProxy(w http.ResponseWriter, r *http.Request, proxyConfig config.Proxy_Service) bool {
	networkProtocol := "tcp"
	targetAddr := proxyConfig.LocalAddress + ":" + proxyConfig.LocalPort


	// Захватываем клиентское соединение через Hijacker
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "WebSocket не поддерживается", http.StatusInternalServerError)
		tools.Logs_file(1, "WS-PROXY", "❌ Hijacker не поддерживается", "logs_proxy.log", false)
		return true
	}

	// Подключаемся к бэкенд серверу
	var backendConn net.Conn
	var err error

	if proxyConfig.ServiceHTTPSuse {
		// TLS соединение для wss://
		backendConn, err = tls.Dial(networkProtocol, targetAddr, &tls.Config{
			InsecureSkipVerify: true,
		})
	} else {
		// Обычное TCP соединение для ws://
		backendConn, err = net.DialTimeout(networkProtocol, targetAddr, 10*time.Second)
	}

	if err != nil {
		http.Error(w, "Ошибка подключения к бэкенду", http.StatusBadGateway)
		tools.Logs_file(1, "WS-PROXY", "❌ Ошибка подключения к бэкенду: "+err.Error(), "logs_proxy.log", false)
		return true
	}
	defer backendConn.Close()

	// Формируем WebSocket handshake запрос к бэкенду
	var handshakeReq bytes.Buffer
	handshakeReq.WriteString(r.Method + " " + r.URL.RequestURI() + " HTTP/1.1\r\n")
	handshakeReq.WriteString("Host: " + r.Host + "\r\n")

	// Копируем все заголовки от клиента (включая WebSocket заголовки!)
	for name, values := range r.Header {
		for _, value := range values {
			handshakeReq.WriteString(name + ": " + value + "\r\n")
		}
	}

	// Добавляем X-Forwarded заголовки
	clientIP := r.RemoteAddr
	if colonIndex := strings.LastIndex(clientIP, ":"); colonIndex != -1 {
		clientIP = clientIP[:colonIndex]
	}
	handshakeReq.WriteString("X-Real-IP: " + clientIP + "\r\n")
	handshakeReq.WriteString("X-Forwarded-For: " + clientIP + "\r\n")
	handshakeReq.WriteString("\r\n")

	// Отправляем handshake на бэкенд
	if _, err := backendConn.Write(handshakeReq.Bytes()); err != nil {
		http.Error(w, "Ошибка отправки handshake", http.StatusBadGateway)
		tools.Logs_file(1, "WS-PROXY", "❌ Ошибка отправки handshake: "+err.Error(), "logs_proxy.log", false)
		return true
	}

	// Читаем ответ от бэкенда
	backendReader := bufio.NewReader(backendConn)
	resp, err := http.ReadResponse(backendReader, r)
	if err != nil {
		http.Error(w, "Ошибка чтения ответа бэкенда", http.StatusBadGateway)
		tools.Logs_file(1, "WS-PROXY", "❌ Ошибка чтения ответа: "+err.Error(), "logs_proxy.log", false)
		return true
	}

	// Проверяем, что бэкенд принял WebSocket (101 Switching Protocols)
	if resp.StatusCode != http.StatusSwitchingProtocols {
		http.Error(w, "Бэкенд отклонил WebSocket", resp.StatusCode)
		tools.Logs_file(1, "WS-PROXY", "❌ Бэкенд отклонил WebSocket: "+resp.Status, "logs_proxy.log", false)
		return true
	}

	// Захватываем клиентское соединение
	clientConn, clientBuf, err := hijacker.Hijack()
	if err != nil {
		tools.Logs_file(1, "WS-PROXY", "❌ Ошибка Hijack: "+err.Error(), "logs_proxy.log", false)
		return true
	}
	defer clientConn.Close()

	// Отправляем ответ handshake клиенту
	var handshakeResp bytes.Buffer
	handshakeResp.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	for name, values := range resp.Header {
		for _, value := range values {
			handshakeResp.WriteString(name + ": " + value + "\r\n")
		}
	}
	handshakeResp.WriteString("\r\n")

	if _, err := clientConn.Write(handshakeResp.Bytes()); err != nil {
		tools.Logs_file(1, "WS-PROXY", "❌ Ошибка отправки ответа клиенту: "+err.Error(), "logs_proxy.log", false)
		return true
	}


	// Двунаправленное проксирование данных
	done := make(chan struct{}, 2)

	// Клиент → Бэкенд
	go func() {
		defer func() { done <- struct{}{} }()

		// Сначала отправляем буферизированные данные если есть
		if clientBuf.Reader.Buffered() > 0 {
			buffered := make([]byte, clientBuf.Reader.Buffered())
			clientBuf.Read(buffered)
			backendConn.Write(buffered)
		}

		io.Copy(backendConn, clientConn)
	}()

	// Бэкенд → Клиент
	go func() {
		defer func() { done <- struct{}{} }()

		// Сначала отправляем буферизированные данные если есть
		if backendReader.Buffered() > 0 {
			buffered := make([]byte, backendReader.Buffered())
			backendReader.Read(buffered)
			clientConn.Write(buffered)
		}

		io.Copy(clientConn, backendConn)
	}()

	// Ждём завершения одного из направлений
	<-done

	return true
}

func StartHandlerProxy(w http.ResponseWriter, r *http.Request) (valid bool) {
	valid = false

	configMutex.RLock()
	defer configMutex.RUnlock()

	// Проверяем глобальный флаг прокси
	if !config.ConfigData.Soft_Settings.Proxy_enabled {
		return false
	}

	// Проходим по всем прокси конфигурациям
	for _, proxyConfig := range config.ConfigData.Proxy_Service {
		// Пропускаем отключенные прокси
		if !proxyConfig.Enable {
			continue
		}

		// Проверяем совпадение домена
		if r.Host != proxyConfig.ExternalDomain {
			continue
		}

		valid = true

		// Проверяем vAccess для прокси
		accessAllowed, errorPage := CheckProxyVAccess(r.URL.Path, proxyConfig.ExternalDomain, r)
		if !accessAllowed {
			// Доступ запрещён - обрабатываем страницу ошибки
			HandleProxyVAccessError(w, r, errorPage)
			return valid
		}

		// Проверяем WebSocket запрос — обрабатываем отдельно
		if isWebSocketRequest(r) {
			return handleWebSocketProxy(w, r, proxyConfig)
		}

		// Проверяем AutoHTTPS - редирект с HTTP на HTTPS
		https_check := !(r.TLS == nil)
		if !https_check && proxyConfig.AutoHTTPS {
			// Перенаправляем на HTTPS
			httpsURL := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
			tools.Logs_file(0, "P-HTTP", "🔀 IP клиента: "+r.RemoteAddr+" Редирект HTTP → HTTPS: "+r.Host+r.URL.Path, "logs_http.log", false)
			return valid
		}

		// Логирование прокси-запроса
		if https_check {
			tools.Logs_file(0, "P-HTTPS", "🔍 IP клиента: "+r.RemoteAddr+" Обработка запроса: https://"+r.Host+r.URL.Path+" → "+proxyConfig.LocalAddress+":"+proxyConfig.LocalPort, "logs_https.log", false)
		} else {
			tools.Logs_file(0, "P-HTTP", "🔍 IP клиента: "+r.RemoteAddr+" Обработка запроса: http://"+r.Host+r.URL.Path+" → "+proxyConfig.LocalAddress+":"+proxyConfig.LocalPort, "logs_http.log", false)
		}

		// Определяем протокол для локального соединения
		protocol := "http"
		if proxyConfig.ServiceHTTPSuse {
			protocol = "https"
		}

		// Читаем тело запроса в буфер для корректной передачи POST данных
		var bodyBuffer bytes.Buffer
		if r.Body != nil {
			if _, err := io.Copy(&bodyBuffer, r.Body); err != nil {
				http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
				return valid
			}
			r.Body.Close()
		}

		// Проксирование на локальный адрес
		proxyURL := protocol + "://" + proxyConfig.LocalAddress + ":" + proxyConfig.LocalPort + r.URL.RequestURI()
		proxyReq, err := http.NewRequest(r.Method, proxyURL, &bodyBuffer)
		if err != nil {
			http.Error(w, "Ошибка создания прокси-запроса", http.StatusInternalServerError)
			return valid
		}

		// Копируем ВСЕ заголовки без изменений (кроме технических)
		for name, values := range r.Header {
			// Пропускаем только технические заголовки HTTP/1.1
			lowerName := strings.ToLower(name)
			if lowerName == "connection" || lowerName == "upgrade" ||
				lowerName == "proxy-connection" || lowerName == "te" ||
				lowerName == "trailers" || lowerName == "transfer-encoding" {
				continue
			}

			// Копируем заголовок как есть
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}

		// Добавляем заголовки для передачи реального IP клиента
		clientIP := r.RemoteAddr
		if colonIndex := strings.LastIndex(clientIP, ":"); colonIndex != -1 {
			clientIP = clientIP[:colonIndex]
		}
		proxyReq.Header.Set("X-Real-IP", clientIP)
		proxyReq.Header.Set("X-Forwarded-For", clientIP)
		proxyReq.Header.Set("X-Forwarded-Proto", protocol)

		// Устанавливаем правильный Content-Length для POST/PUT запросов
		if bodyBuffer.Len() > 0 {
			proxyReq.ContentLength = int64(bodyBuffer.Len())
		}

		// Выполняем прокси-запрос
		client := &http.Client{
			// Отключаем автоматическое следование редиректам для корректной работы с авторизацией
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		// Для HTTPS соединений настраиваем TLS (если понадобится)
		if proxyConfig.ServiceHTTPSuse {
			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // Простая настройка для внутренних соединений
				},
			}
		}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, "Ошибка прокси-запроса", http.StatusBadGateway)
			tools.Logs_file(1, "PROXY", "Ошибка прокси-запроса: "+err.Error(), "logs_proxy.log", false)
			return valid
		}
		defer resp.Body.Close()

		// Прозрачно копируем ВСЕ заголовки ответа без изменений
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// Сжатие ответа (gzip) — только если бэкенд не сжал сам
		var gzw *gzipResponseWriter
		if proxyConfig.IsCompressionEnabled() && clientAcceptsGzip(r) && !isAlreadyCompressed(resp.Header) {
			gzw = newGzipResponseWriter(w)
			defer gzw.close()
			w = gzw
		}

		// Устанавливаем статус код
		w.WriteHeader(resp.StatusCode)

		// Копируем тело ответа с поддержкой streaming (SSE, chunked responses)
		flusher, canFlush := w.(http.Flusher)
		
		buffer := make([]byte, 32*1024)
		
		for {
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				if _, writeErr := w.Write(buffer[:n]); writeErr != nil {
					log.Printf("Ошибка записи тела ответа: %v", writeErr)
					break
				}
				
				if canFlush {
					flusher.Flush()
				}
			}
			
			if err != nil {
				if err != io.EOF {
					log.Printf("Ошибка чтения тела ответа: %v", err)
				}
				break
			}
		}

		return valid
	}

	return valid
}
