package vaccess

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetVAccessPath(host string, isProxy bool) string {
	if isProxy {
		return fmt.Sprintf("WebServer/tools/Proxy_vAccess/%s_vAccess.conf", host)
	}
	return fmt.Sprintf("WebServer/www/%s/vAccess.conf", host)
}

func GetVAccessConfig(host string, isProxy bool) (*VAccessConfig, error) {
	filePath := GetVAccessPath(host, isProxy)

	// Получаем абсолютный путь БЕЗ проверки существования
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return &VAccessConfig{Rules: []VAccessRule{}}, nil
	}

	// Проверяем существование файла
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		// Файл не существует - возвращаем пустую конфигурацию
		return &VAccessConfig{Rules: []VAccessRule{}}, nil
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &VAccessConfig{}
	scanner := bufio.NewScanner(file)

	var currentRule *VAccessRule

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Парсим параметры
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "type":
					// Начало нового правила - сохраняем предыдущее
					if currentRule != nil && currentRule.Type != "" {
						config.Rules = append(config.Rules, *currentRule)
					}
					// Создаём новое правило
					currentRule = &VAccessRule{}
					currentRule.Type = value
				case "type_file":
					if currentRule != nil {
						currentRule.TypeFile = splitAndTrim(value)
					}
				case "path_access":
					if currentRule != nil {
						currentRule.PathAccess = splitAndTrim(value)
					}
				case "ip_list":
					if currentRule != nil {
						currentRule.IPList = splitAndTrim(value)
					}
				case "exceptions_dir":
					if currentRule != nil {
						currentRule.ExceptionsDir = splitAndTrim(value)
					}
				case "url_error":
					if currentRule != nil {
						currentRule.UrlError = value
					}
				}
			}
		}
	}

	// Добавляем последнее правило
	if currentRule != nil && currentRule.Type != "" {
		config.Rules = append(config.Rules, *currentRule)
	}

	return config, nil
}

func SaveVAccessConfig(host string, isProxy bool, config *VAccessConfig) error {
	filePath := GetVAccessPath(host, isProxy)

	// Создаём директорию если не существует
	dir := ""
	if isProxy {
		dir = "WebServer/tools/Proxy_vAccess"
	} else {
		dir = fmt.Sprintf("WebServer/www/%s", host)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	os.MkdirAll(absDir, 0755)

	// Получаем абсолютный путь к файлу
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	// Формируем содержимое файла
	var content strings.Builder

	content.WriteString("# vAccess Configuration\n")
	content.WriteString("# Правила применяются сверху вниз\n\n")

	for i, rule := range config.Rules {
		content.WriteString(fmt.Sprintf("# Правило %d\n", i+1))
		content.WriteString(fmt.Sprintf("type: %s\n", rule.Type))

		if len(rule.TypeFile) > 0 {
			content.WriteString(fmt.Sprintf("type_file: %s\n", strings.Join(rule.TypeFile, ", ")))
		}
		if len(rule.PathAccess) > 0 {
			content.WriteString(fmt.Sprintf("path_access: %s\n", strings.Join(rule.PathAccess, ", ")))
		}
		if len(rule.IPList) > 0 {
			content.WriteString(fmt.Sprintf("ip_list: %s\n", strings.Join(rule.IPList, ", ")))
		}
		if len(rule.ExceptionsDir) > 0 {
			content.WriteString(fmt.Sprintf("exceptions_dir: %s\n", strings.Join(rule.ExceptionsDir, ", ")))
		}
		if rule.UrlError != "" {
			content.WriteString(fmt.Sprintf("url_error: %s\n", rule.UrlError))
		}

		content.WriteString("\n")
	}

	return os.WriteFile(absPath, []byte(content.String()), 0644)
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := []string{}
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
