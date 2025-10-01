package webserver

import (
	"net/http"
	"os"
	"strings"
	"vServer/Backend/config"
	tools "vServer/Backend/tools"
)

func StartHandler() {
	http.HandleFunc("/", handler)
}

func Alias_check(r *http.Request) (alias_found bool, host string) {

	alias_found = false

	for _, site := range config.ConfigData.Site_www {

		for _, alias := range site.Alias {

			if alias == r.Host {
				alias_found = true
				return alias_found, site.Host

			} else {
				alias_found = false
			}
		}
	}

	return alias_found, ""

}

func Alias_Run(r *http.Request) (rhost string) {

	var host string
	host = r.Host

	alias_check, alias := Alias_check(r)

	if alias_check {
		host = alias
	}

	return host
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
	// По умолчанию роутинг выключен
	return false
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

// Обработчик запросов
func handler(w http.ResponseWriter, r *http.Request) {

	host := Alias_Run(r)           // Получаем хост из запроса
	https_check := !(r.TLS == nil) // Проверяем, по HTTPS ли запрос
	root_url := r.URL.Path == "/"  // Проверяем, является ли запрос корневым URL

	// Проверяем, обработал ли прокси запрос
	if StartHandlerProxy(w, r) {
		return // Если прокси обработал запрос, прерываем выполнение
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
