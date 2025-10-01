package webserver

import (
	"bufio"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	tools "vServer/Backend/tools"
)

// Структура для правила vAccess
type VAccessRule struct {
	Type          string   // "Allow" или "Disable"
	TypeFile      []string // Список расширений файлов
	PathAccess    []string // Список путей для применения правила
	IPList        []string // Список IP адресов для фильтрации
	ExceptionsDir []string // Список путей-исключений (не применять правило к этим путям)
	UrlError      string   // Страница ошибки: "404", внешний URL или локальный путь
}

// Структура для конфигурации vAccess
type VAccessConfig struct {
	Rules []VAccessRule
}

// Проверка валидности правила
func isValidRule(rule *VAccessRule) bool {
	// Минимум нужен Type
	if rule.Type == "" {
		return false
	}

	// Должно быть хотя бы одно условие: type_file, path_access или ip_list
	hasCondition := len(rule.TypeFile) > 0 || len(rule.PathAccess) > 0 || len(rule.IPList) > 0

	return hasCondition
}

// Парсинг vAccess.conf файла
func parseVAccessFile(filePath string) (*VAccessConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &VAccessConfig{}
	scanner := bufio.NewScanner(file)

	var currentRule *VAccessRule

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Пропускаем пустые строки
		if line == "" {
			continue
		}

		// Комментарии разделяют правила
		if strings.HasPrefix(line, "#") {
			// Если есть текущее правило, сохраняем его перед началом нового
			if currentRule != nil && isValidRule(currentRule) {
				config.Rules = append(config.Rules, *currentRule)
				currentRule = nil
			}
			continue
		}

		// Парсим строки конфигурации
		if strings.HasPrefix(line, "type:") {
			// Создаём новое правило только если его нет
			if currentRule == nil {
				currentRule = &VAccessRule{}
			}
			currentRule.Type = strings.TrimSpace(strings.TrimPrefix(line, "type:"))

		} else if strings.HasPrefix(line, "type_file:") && currentRule != nil {
			fileTypes := strings.TrimSpace(strings.TrimPrefix(line, "type_file:"))
			// Разбиваем по запятым и очищаем пробелы
			for _, fileType := range strings.Split(fileTypes, ",") {
				fileType = strings.TrimSpace(fileType)
				if fileType != "" {
					currentRule.TypeFile = append(currentRule.TypeFile, fileType)
				}
			}

		} else if strings.HasPrefix(line, "path_access:") && currentRule != nil {
			pathAccess := strings.TrimSpace(strings.TrimPrefix(line, "path_access:"))
			// Разбиваем по запятым и очищаем пробелы
			for _, path := range strings.Split(pathAccess, ",") {
				path = strings.TrimSpace(path)
				if path != "" {
					currentRule.PathAccess = append(currentRule.PathAccess, path)
				}
			}

		} else if strings.HasPrefix(line, "ip_list:") && currentRule != nil {
			ipList := strings.TrimSpace(strings.TrimPrefix(line, "ip_list:"))
			// Разбиваем по запятым и очищаем пробелы
			for _, ip := range strings.Split(ipList, ",") {
				ip = strings.TrimSpace(ip)
				if ip != "" {
					currentRule.IPList = append(currentRule.IPList, ip)
				}
			}

		} else if strings.HasPrefix(line, "exceptions_dir:") && currentRule != nil {
			exceptionsDir := strings.TrimSpace(strings.TrimPrefix(line, "exceptions_dir:"))
			// Разбиваем по запятым и очищаем пробелы
			for _, exception := range strings.Split(exceptionsDir, ",") {
				exception = strings.TrimSpace(exception)
				if exception != "" {
					currentRule.ExceptionsDir = append(currentRule.ExceptionsDir, exception)
				}
			}

		} else if strings.HasPrefix(line, "url_error:") && currentRule != nil {
			currentRule.UrlError = strings.TrimSpace(strings.TrimPrefix(line, "url_error:"))
		}
	}

	// Добавляем последнее правило если оно валидно
	if currentRule != nil && isValidRule(currentRule) {
		config.Rules = append(config.Rules, *currentRule)
	}

	return config, scanner.Err()
}

// Поиск всех vAccess.conf файлов от корня сайта до запрашиваемого пути
func findVAccessFiles(requestPath string, host string) []string {
	var configFiles []string

	// Базовый путь к сайту (НЕ public_www, а уровень выше)
	basePath := "WebServer/www/" + host

	// Проверяем корневой vAccess.conf
	rootConfigPath := filepath.Join(basePath, "vAccess.conf")
	if _, err := os.Stat(rootConfigPath); err == nil {
		configFiles = append(configFiles, rootConfigPath)
	}

	// Разбиваем путь на части для поиска вложенных конфигов
	pathParts := strings.Split(strings.Trim(requestPath, "/"), "/")
	currentPath := basePath

	for _, part := range pathParts {
		if part == "" {
			continue
		}
		currentPath = filepath.Join(currentPath, part)
		configPath := filepath.Join(currentPath, "vAccess.conf")

		if _, err := os.Stat(configPath); err == nil {
			configFiles = append(configFiles, configPath)
		}
	}

	return configFiles
}

// Проверка соответствия пути правилу
func matchPath(rulePath, requestPath string) bool {
	// Если правило заканчивается на /*, проверяем префикс
	if strings.HasSuffix(rulePath, "/*") {
		prefix := strings.TrimSuffix(rulePath, "/*")

		// Специальный случай: /* должен совпадать со всеми путями
		if prefix == "" {
			return true
		}

		return strings.HasPrefix(requestPath, prefix)
	}

	// Точное совпадение
	return rulePath == requestPath
}

// Извлечение всех расширений из пути
func getAllExtensionsFromPath(filePath string) []string {
	var extensions []string

	// Разбиваем путь на части по слэшам
	parts := strings.Split(filePath, "/")

	for _, part := range parts {
		// Ищем все точки в каждой части пути
		if strings.Contains(part, ".") {
			// Находим все расширения в части (может быть несколько: file.tar.gz)
			dotIndex := strings.Index(part, ".")
			for dotIndex != -1 && dotIndex < len(part)-1 {
				// Извлекаем расширение от точки до следующей точки или конца
				nextDotIndex := strings.Index(part[dotIndex+1:], ".")
				if nextDotIndex == -1 {
					// Последнее расширение
					ext := strings.ToLower(part[dotIndex:])
					if ext != "." && len(ext) > 1 {
						extensions = append(extensions, ext)
					}
					break
				} else {
					// Промежуточное расширение
					ext := strings.ToLower(part[dotIndex : dotIndex+1+nextDotIndex+1])
					if ext != "." && len(ext) > 1 {
						extensions = append(extensions, ext)
					}
					dotIndex = dotIndex + 1 + nextDotIndex
				}
			}
		}
	}

	return extensions
}

// Проверка соответствия расширений файла
// Возвращает true если ВСЕ найденные расширения разрешены
func matchFileExtension(ruleExtensions []string, filePath string) bool {
	// Получаем все расширения из пути
	pathExtensions := getAllExtensionsFromPath(filePath)

	// Если расширений нет, проверяем есть ли no_extension в правилах
	if len(pathExtensions) == 0 {
		for _, ruleExt := range ruleExtensions {
			ruleExt = strings.ToLower(strings.TrimSpace(ruleExt))
			if ruleExt == "no_extension" {
				return true
			}
		}
		return false
	}

	// Проверяем каждое найденное расширение
	for _, pathExt := range pathExtensions {
		found := false
		for _, ruleExt := range ruleExtensions {
			ruleExt = strings.ToLower(strings.TrimSpace(ruleExt))

			// Поддержка паттернов типа *.php
			if strings.HasPrefix(ruleExt, "*.") {
				if pathExt == strings.TrimPrefix(ruleExt, "*") {
					found = true
					break
				}
			} else if ruleExt == pathExt {
				found = true
				break
			}
		}

		// Если хотя бы одно расширение не найдено в правилах - блокируем
		if !found {
			return false
		}
	}

	// Все расширения найдены в правилах
	return true
}

// Получение реального IP адреса клиента из соединения (без заголовков прокси)
func getClientIP(r *http.Request) string {
	// Извлекаем IP из RemoteAddr (формат: "IP:port")
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	// Убираем квадратные скобки для IPv6
	ip = strings.Trim(ip, "[]")

	return ip
}

// Проверка соответствия IP адреса правилу
func matchIPAddress(ruleIPs []string, clientIP string) bool {
	if len(ruleIPs) == 0 {
		return true // Если IP не указаны, то проверка пройдена
	}

	for _, ruleIP := range ruleIPs {
		ruleIP = strings.TrimSpace(ruleIP)
		if ruleIP == clientIP {
			return true
		}
	}

	return false
}

// Проверка исключений - возвращает true если путь находится в исключениях
func matchExceptions(exceptions []string, requestPath string) bool {
	if len(exceptions) == 0 {
		return false // Нет исключений
	}

	for _, exception := range exceptions {
		exception = strings.TrimSpace(exception)
		if matchPath(exception, requestPath) {
			return true // Путь найден в исключениях
		}
	}

	return false
}

// Основная функция проверки доступа
// Возвращает (разрешён_доступ, страница_ошибки)
func CheckVAccess(requestPath string, host string, r *http.Request) (bool, string) {
	// Находим все vAccess.conf файлы
	configFiles := findVAccessFiles(requestPath, host)

	if len(configFiles) == 0 {
		// Нет конфигурационных файлов - разрешаем доступ
		return true, ""
	}

	// Применяем правила по порядку (от корня к файлу)
	for _, configFile := range configFiles {
		config, err := parseVAccessFile(configFile)
		if err != nil {
			tools.Logs_file(1, "vAccess", "❌ Ошибка парсинга "+configFile+": "+err.Error(), "logs_vaccess.log", false)
			continue
		}

		// Проверяем каждое правило в конфиге
		for _, rule := range config.Rules {
			// Проверяем соответствие путей (если указаны)
			pathMatched := true // По умолчанию true, если путей нет
			if len(rule.PathAccess) > 0 {
				pathMatched = false
				for _, rulePath := range rule.PathAccess {
					if matchPath(rulePath, requestPath) {
						pathMatched = true
						break
					}
				}
			}

			// Если путь не совпадает - переходим к следующему правилу
			if !pathMatched {
				continue
			}

			// Проверяем исключения - если путь в исключениях, пропускаем правило
			if matchExceptions(rule.ExceptionsDir, requestPath) {
				continue
			}

			// Проверяем соответствие расширения файла (если указаны)
			fileMatches := true // По умолчанию true, если типов файлов нет
			if len(rule.TypeFile) > 0 {
				fileMatches = matchFileExtension(rule.TypeFile, requestPath)
			}

			// Проверяем соответствие IP адреса (если указаны)
			ipMatches := true // По умолчанию true, если IP не указаны
			if len(rule.IPList) > 0 {
				clientIP := getClientIP(r)
				ipMatches = matchIPAddress(rule.IPList, clientIP)
			}

			// Применяем правило в зависимости от типа
			switch rule.Type {
			case "Allow":
				// Allow правило: разрешаем только если ВСЕ условия выполнены
				if (len(rule.TypeFile) > 0 && !fileMatches) || (len(rule.IPList) > 0 && !ipMatches) {
					// Условия НЕ выполнены - блокируем
					errorPage := rule.UrlError
					if errorPage == "" {
						errorPage = "404" // По умолчанию 404
					}
					return false, errorPage
				}
				// Все условия Allow выполнены - разрешаем доступ
				return true, ""
			case "Disable":
				// Disable правило: запрещаем если ЛЮБОЕ условие выполнено
				if (len(rule.TypeFile) == 0 || fileMatches) && (len(rule.IPList) == 0 || ipMatches) {
					errorPage := rule.UrlError
					if errorPage == "" {
						errorPage = "404" // По умолчанию 404
					}
					return false, errorPage
				}
			default:
				// Неизвестный тип правила - игнорируем
				continue
			}
		}
	}

	// Все проверки пройдены - разрешаем доступ
	return true, ""
}

// Обработка страницы ошибки vAccess
func HandleVAccessError(w http.ResponseWriter, r *http.Request, errorPage string, host string) {
	switch {
	case errorPage == "404":
		// Стандартная 404 страница
		http.ServeFile(w, r, "WebServer/tools/error_page/index.html")

	case strings.HasPrefix(errorPage, "http://") || strings.HasPrefix(errorPage, "https://"):
		// Внешний сайт - редирект
		http.Redirect(w, r, errorPage, http.StatusFound)

	default:
		// Локальный путь от public_www
		localPath := "WebServer/www/" + host + "/public_www" + errorPage
		if _, err := os.Stat(localPath); err == nil {
			http.ServeFile(w, r, localPath)
		} else {
			// Файл не найден - показываем стандартную 404
			http.ServeFile(w, r, "WebServer/tools/error_page/index.html")
			tools.Logs_file(1, "vAccess", "❌ Страница ошибки не найдена: "+localPath, "logs_vaccess.log", false)
		}
	}
}
