package webserver

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"vServer/Backend/config"
	tools "vServer/Backend/tools"
	"vServer/Backend/WebServer/acme"
)

var (
	siteStatusCache map[string]bool
	statusMutex     sync.RWMutex
)

func StartHandler() {
	http.HandleFunc("/", handler)
	updateSiteStatusCache()
}

func updateSiteStatusCache() {
	statusMutex.Lock()
	defer statusMutex.Unlock()

	siteStatusCache = make(map[string]bool)
	for _, site := range config.ConfigData.Site_www {
		siteStatusCache[site.Host] = site.Status == "active"
	}
}

// UpdateSiteStatusCache - экспортируемая функция для обновления кэша
func UpdateSiteStatusCache() {
	updateSiteStatusCache()
}

// Проверка wildcard паттерна для alias
func matchWildcardAlias(pattern, host string) bool {
	// Если нет звёздочки - точное совпадение
	if !strings.Contains(pattern, "*") {
		return pattern == host
	}

	// Поддержка wildcard: *.example.com, example.*, *example*, *
	// Паттерн *.example.com -> должен совпадать с sub.example.com

	// Если паттерн = * (любой хост)
	if pattern == "*" {
		return true
	}

	// Разбиваем паттерн на части по звёздочке
	parts := strings.Split(pattern, "*")

	// Проверяем каждую часть
	currentPos := 0
	for i, part := range parts {
		if part == "" {
			continue
		}

		// Ищем часть в хосте начиная с текущей позиции
		idx := strings.Index(host[currentPos:], part)
		if idx == -1 {
			return false
		}

		// Для первой части проверяем что она в начале (если паттерн не начинается с *)
		if i == 0 && !strings.HasPrefix(pattern, "*") {
			if idx != 0 {
				return false
			}
		}

		// Для последней части проверяем что она в конце (если паттерн не кончается на *)
		if i == len(parts)-1 && !strings.HasSuffix(pattern, "*") {
			if currentPos+idx+len(part) != len(host) {
				return false
			}
		}

		currentPos += idx + len(part)
	}

	return true
}

func Alias_Run(r *http.Request) (rhost string) {
	requestHost := r.Host

	// Убираем порт если есть (например :80 или :443)
	if colonIndex := strings.Index(requestHost, ":"); colonIndex != -1 {
		requestHost = requestHost[:colonIndex]
	}

	// Приоритет 1: Проверяем точное совпадение с site.Host
	for _, site := range config.ConfigData.Site_www {
		if site.Host == requestHost {
			return site.Host
		}
	}

	// Приоритет 2: Проверяем точные alias (без wildcard)
	for _, site := range config.ConfigData.Site_www {
		for _, alias := range site.Alias {
			if !strings.Contains(alias, "*") && alias == requestHost {
				return site.Host
			}
		}
	}

	// Приоритет 3: Проверяем wildcard alias
	for _, site := range config.ConfigData.Site_www {
		for _, alias := range site.Alias {
			if strings.Contains(alias, "*") && matchWildcardAlias(alias, requestHost) {
				return site.Host
			}
		}
	}

	// Не нашли совпадений - возвращаем как есть
	return requestHost
}

// Получает список root_file для сайта из конфигурации
func getRootFiles(host string) []string {
	for _, site := range config.ConfigData.Site_www {
		if site.Host == host {
			if site.Root_file != "" {
				// Разделяем по запятой и убираем пробелы
				files := strings.Split(site.Root_file, ",")
				var cleanFiles []string
				for _, file := range files {
					cleanFile := strings.TrimSpace(file)
					if cleanFile != "" {
						cleanFiles = append(cleanFiles, cleanFile)
					}
				}
				if len(cleanFiles) > 0 {
					return cleanFiles
				}
			}
			// Если не указан, используем index.html как fallback
			return []string{"index.html"}
		}
	}
	// Если сайт не найден в конфиге, используем index.html
	return []string{"index.html"}
}

// Находит первый существующий root файл из списка
func findExistingRootFile(host string, dirPath string) (string, bool) {
	rootFiles := getRootFiles(host)
	basePath := "WebServer/www/" + host + "/public_www" + dirPath

	for _, rootFile := range rootFiles {
		fullPath := basePath + rootFile
		if _, err := os.Stat(fullPath); err == nil {
			return rootFile, true
		}
	}
	return "", false
}

// Проверяет включен ли роутинг через root файл для сайта
func isRootFileRoutingEnabled(host string) bool {
	for _, site := range config.ConfigData.Site_www {
		if site.Host == host {
			return site.Root_file_routing
		}
	}
	return false
}

// Проверяет включено ли сжатие для сайта
func isSiteCompressionEnabled(host string) bool {
	for _, site := range config.ConfigData.Site_www {
		if site.Host == host {
			return site.IsCompressionEnabled()
		}
	}
	return true
}

// Проверка vAccess с обработкой ошибки
// Возвращает true если доступ разрешён, false если заблокирован
func checkVAccessAndHandle(w http.ResponseWriter, r *http.Request, filePath string, host string) bool {
	accessAllowed, errorPage := CheckVAccess(filePath, host, r)
	if !accessAllowed {
		HandleVAccessError(w, r, errorPage, host)
		tools.Logs_file(2, "vAccess", "🚫 Доступ запрещён vAccess: "+r.RemoteAddr+" → "+r.Host+filePath+" (error: "+errorPage+")", "logs_vaccess.log", false)
		return false
	}
	return true
}

// Проверяет включен ли сайт (оптимизировано через кэш)
func isSiteActive(host string) bool {
	statusMutex.RLock()
	defer statusMutex.RUnlock()

	if status, exists := siteStatusCache[host]; exists {
		return status
	}
	return false
}

// Обработчик запросов
func handler(w http.ResponseWriter, r *http.Request) {

	// ACME HTTP-01 Challenge (для Let's Encrypt)
	if strings.HasPrefix(r.URL.Path, "/.well-known/acme-challenge/") {
		if acme.DefaultManager != nil && acme.DefaultManager.HandleChallenge(w, r) {
			return
		}
	}

	host := Alias_Run(r)           // Получаем хост из запроса
	https_check := !(r.TLS == nil) // Проверяем, по HTTPS ли запрос
	root_url := r.URL.Path == "/"  // Проверяем, является ли запрос корневым URL

	// Проверяем, обработал ли прокси запрос
	if StartHandlerProxy(w, r) {
		return // Если прокси обработал запрос, прерываем выполнение
	}

	// Проверяем статус сайта
	if !isSiteActive(host) {
		http.ServeFile(w, r, "WebServer/tools/error_page/index.html")
		tools.Logs_file(2, "H503", "🚫 Сайт отключен: "+host, "logs_http.log", false)
		return
	}

	// ЕДИНСТВЕННАЯ ПРОВЕРКА vAccess - простая проверка запрошенного пути
	if !checkVAccessAndHandle(w, r, r.URL.Path, host) {
		return
	}

	if https_check {

		tools.Logs_file(0, "HTTPS", "🔍 IP клиента: "+r.RemoteAddr+" Обработка запроса: https://"+r.Host+r.URL.Path, "logs_https.log", false)

	} else {

		tools.Logs_file(0, "HTTP", "🔍 IP клиента: "+r.RemoteAddr+" Обработка запроса: http://"+r.Host+r.URL.Path, "logs_http.log", false)

		// Если сертификат для домена существует в папке cert, перенаправляем на HTTPS
		if checkHostCert(r) {
			// Если запрос не по HTTPS, перенаправляем на HTTPS
			httpsURL := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
			return // Прерываем выполнение после редиректа
		}

	}

	// Сжатие ответа (gzip)
	if isSiteCompressionEnabled(host) && clientAcceptsGzip(r) {
		gzw := newGzipResponseWriter(w)
		defer gzw.close()
		w = gzw
	}

	// Проверяем существование директории сайта
	if _, err := os.Stat("WebServer/www/" + host + "/public_www"); err != nil {
		http.ServeFile(w, r, "WebServer/tools/error_page/index.html")
		tools.Logs_file(2, "H404", "🔍 IP клиента: "+r.RemoteAddr+" Директория сайта не найдена: "+host, "logs_http.log", false)
		return
	}

	if root_url {
		// Если корневой URL, то ищем первый существующий root файл
		if rootFile, found := findExistingRootFile(host, "/"); found {
			// Обрабатываем найденный root файл (статический или PHP)
			HandlePHPRequest(w, r, host, "/"+rootFile, r.URL.RequestURI(), r.URL.Path)
		} else {
			// Ни один root файл не найден - показываем ошибку
			rootFiles := getRootFiles(host)
			tools.Logs_file(2, "H404", "🔍 IP клиента: "+r.RemoteAddr+" Root файлы не найдены: "+strings.Join(rootFiles, ", "), "logs_http.log", false)
			http.ServeFile(w, r, "WebServer/tools/error_page/index.html")
		}
	}

	if !root_url {

		// Проверяем существование запрашиваемого файла
		filePath := "WebServer/www/" + host + "/public_www" + r.URL.Path

		if fileInfo, err := os.Stat(filePath); err == nil {
			// Путь существует - проверяем что это
			if fileInfo.IsDir() {
				// Это директория - ищем индексные файлы
				// Убираем слэш в конце если есть, и добавляем обратно для единообразия
				dirPath := r.URL.Path
				if !strings.HasSuffix(dirPath, "/") {
					dirPath += "/"
				}

				// Ищем первый существующий root файл в директории
				if rootFile, found := findExistingRootFile(host, dirPath); found {
					// Обрабатываем найденный индексный файл в директории
					HandlePHPRequest(w, r, host, dirPath+rootFile, r.URL.RequestURI(), r.URL.Path)
					return
				}

				// Если никаких индексных файлов нет - показываем ошибку (запрещаем листинг)
				rootFiles := getRootFiles(host)
				tools.Logs_file(2, "H404", "🔍 IP клиента: "+r.RemoteAddr+" Индексные файлы не найдены в директории "+r.Host+r.URL.Path+": "+strings.Join(rootFiles, ", "), "logs_http.log", false)
				http.ServeFile(w, r, "WebServer/tools/error_page/index.html")

			} else {
				// Это файл - обрабатываем через HandlePHPRequest
				HandlePHPRequest(w, r, host, r.URL.Path, "", "")
			}

		} else {
			// Файл не найден - проверяем нужен ли роутинг через root файл
			if isRootFileRoutingEnabled(host) {
				// Ищем первый существующий root файл для роутинга
				if rootFile, found := findExistingRootFile(host, "/"); found {
					// Root файл существует - используем для роутинга
					HandlePHPRequest(w, r, host, "/"+rootFile, r.URL.RequestURI(), r.URL.Path)
				} else {
					// Root файлы не найдены
					rootFiles := getRootFiles(host)
					tools.Logs_file(2, "H404", "🔍 IP клиента: "+r.RemoteAddr+" Root файлы не найдены для роутинга: "+strings.Join(rootFiles, ", "), "logs_http.log", false)
					http.ServeFile(w, r, "WebServer/tools/error_page/index.html")
				}
			} else {
				// Роутинг отключен - показываем обычную 404
				http.ServeFile(w, r, "WebServer/tools/error_page/index.html")
				tools.Logs_file(2, "H404", "🔍 IP клиента: "+r.RemoteAddr+" Файл не найден: "+r.Host+r.URL.Path, "logs_http.log", false)
			}
		}
	}
}
