package sites

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	config "vServer/Backend/config"
	tools "vServer/Backend/tools"
)

//go:embed templates/index.tmpl
var indexTemplate string

// CreateNewSite создаёт новый сайт со всей необходимой структурой
func CreateNewSite(siteData SiteInfo) error {
	// 1. Валидация данных
	if err := ValidateSiteData(siteData); err != nil {
		return err
	}

	// 2. Создание структуры папок
	if err := CreateSiteFolder(siteData.Host); err != nil {
		return fmt.Errorf("ошибка создания папок: %w", err)
	}

	// 3. Создание стартового файла
	if err := CreateStarterFile(siteData.Host, siteData.RootFile); err != nil {
		return fmt.Errorf("ошибка создания стартового файла: %w", err)
	}

	// 4. Создание пустого vAccess.conf
	if err := CreateVAccessFile(siteData.Host); err != nil {
		return fmt.Errorf("ошибка создания vAccess.conf: %w", err)
	}

	// 5. Добавление сайта в конфиг
	if err := AddSiteToConfig(siteData); err != nil {
		return fmt.Errorf("ошибка добавления в конфиг: %w", err)
	}

	tools.Logs_file(0, "SITES", fmt.Sprintf("✅ Новый сайт создан: %s (%s)", siteData.Name, siteData.Host), "logs_error.log", true)
	return nil
}

// ValidateSiteData проверяет данные нового сайта
func ValidateSiteData(siteData SiteInfo) error {
	// Проверка обязательных полей
	if strings.TrimSpace(siteData.Host) == "" {
		return errors.New("поле Host обязательно для заполнения")
	}

	if strings.TrimSpace(siteData.Name) == "" {
		return errors.New("поле Name обязательно для заполнения")
	}

	if strings.TrimSpace(siteData.RootFile) == "" {
		return errors.New("поле RootFile обязательно для заполнения")
	}

	// Проверка уникальности host
	for _, site := range config.ConfigData.Site_www {
		if strings.EqualFold(site.Host, siteData.Host) {
			return fmt.Errorf("сайт с host '%s' уже существует", siteData.Host)
		}
	}

	// Проверка валидности status
	if siteData.Status != "active" && siteData.Status != "inactive" {
		return errors.New("status должен быть 'active' или 'inactive'")
	}

	return nil
}

// CreateSiteFolder создаёт структуру папок для нового сайта
func CreateSiteFolder(host string) error {
	// Создаём путь: WebServer/www/{host}/public_www/
	folderPath := filepath.Join("WebServer", "www", host, "public_www")

	absPath, err := tools.AbsPath(folderPath)
	if err != nil {
		return err
	}

	// Создаём все необходимые папки
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку: %w", err)
	}

	tools.Logs_file(0, "SITES", fmt.Sprintf("📁 Создана папка: %s", folderPath), "logs_error.log", false)
	return nil
}

// CreateStarterFile создаёт стартовый файл (index.html или index.php)
func CreateStarterFile(host, rootFile string) error {
	filePath := filepath.Join("WebServer", "www", host, "public_www", rootFile)

	// Получаем абсолютный путь БЕЗ проверки существования
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("ошибка получения абсолютного пути: %w", err)
	}

	// Генерируем контент из шаблона
	content := generateTemplate(host, rootFile)

	// Записываем файл
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}

	tools.Logs_file(0, "SITES", fmt.Sprintf("📄 Создан стартовый файл: %s", rootFile), "logs_error.log", false)
	return nil
}

// CreateVAccessFile создаёт пустой конфиг vAccess
func CreateVAccessFile(host string) error {
	filePath := filepath.Join("WebServer", "www", host, "vAccess.conf")

	// Получаем абсолютный путь БЕЗ проверки существования
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("ошибка получения абсолютного пути: %w", err)
	}

	content := `# vAccess Configuration
# Правила применяются сверху вниз

# Пример правила (закомментировано):
# type: Disable
# type_file: *.php
# path_access: /uploads/*
# url_error: 404

`

	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("не удалось создать vAccess.conf: %w", err)
	}

	tools.Logs_file(0, "SITES", "🔒 Создан vAccess.conf", "logs_error.log", false)
	return nil
}

// AddSiteToConfig добавляет новый сайт в config.json
func AddSiteToConfig(siteData SiteInfo) error {
	// Создаём новую запись
	newSite := config.Site_www{
		Name:              siteData.Name,
		Host:              siteData.Host,
		Alias:             siteData.Alias,
		Status:            siteData.Status,
		Root_file:         siteData.RootFile,
		Root_file_routing: siteData.RootFileRouting,
	}

	// Добавляем в массив
	config.ConfigData.Site_www = append(config.ConfigData.Site_www, newSite)

	// Сохраняем конфиг в файл
	if err := SaveConfig(); err != nil {
		return err
	}

	tools.Logs_file(0, "SITES", "💾 Конфигурация обновлена", "logs_error.log", false)
	return nil
}

// SaveConfig сохраняет текущую конфигурацию в файл
func SaveConfig() error {
	// Форматируем JSON с отступами
	jsonData, err := json.MarshalIndent(config.ConfigData, "", "    ")
	if err != nil {
		return fmt.Errorf("ошибка форматирования JSON: %w", err)
	}

	// Получаем абсолютный путь к файлу конфига
	absPath, err := tools.AbsPath(config.ConfigPath)
	if err != nil {
		return err
	}

	// Записываем в файл
	if err := os.WriteFile(absPath, jsonData, 0644); err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}

// UploadSiteCertificate загружает SSL сертификат для сайта
func UploadSiteCertificate(host, certType string, certData []byte) error {
	// Создаём папку для сертификатов
	certDir := filepath.Join("WebServer", "cert", host)

	absCertDir, err := tools.AbsPath(certDir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(absCertDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку для сертификатов: %w", err)
	}

	// Определяем имя файла
	var fileName string
	switch certType {
	case "certificate":
		fileName = "certificate.crt"
	case "privatekey":
		fileName = "private.key"
	case "cabundle":
		fileName = "ca_bundle.crt"
	default:
		return fmt.Errorf("неизвестный тип сертификата: %s", certType)
	}

	// Путь к файлу сертификата
	certFilePath := filepath.Join(absCertDir, fileName)

	// Записываем файл
	if err := os.WriteFile(certFilePath, certData, 0644); err != nil {
		return fmt.Errorf("не удалось сохранить сертификат: %w", err)
	}

	tools.Logs_file(0, "SITES", fmt.Sprintf("🔒 Загружен сертификат: %s для %s", fileName, host), "logs_error.log", true)
	return nil
}

// DeleteSiteCertificates удаляет сертификаты сайта
func DeleteSiteCertificates(host string) error {
	certDir := filepath.Join("WebServer", "cert", host)

	// Получаем абсолютный путь к папке сертификатов
	absCertDir, err := filepath.Abs(certDir)
	if err != nil {
		return fmt.Errorf("ошибка получения пути: %w", err)
	}

	// Проверяем, существует ли папка
	if _, err := os.Stat(absCertDir); os.IsNotExist(err) {
		return nil // Папки нет - ничего удалять не нужно
	}

	// Удаляем папку со всем содержимым
	if err := os.RemoveAll(absCertDir); err != nil {
		return fmt.Errorf("не удалось удалить папку сертификатов: %w", err)
	}

	tools.Logs_file(0, "SITES", fmt.Sprintf("🗑️ Удалены сертификаты для: %s", host), "logs_error.log", true)
	return nil
}

// DeleteSite полностью удаляет сайт
func DeleteSite(host string) error {
	// 1. Проверяем, существует ли сайт в конфиге
	siteIndex := -1
	for i, site := range config.ConfigData.Site_www {
		if site.Host == host {
			siteIndex = i
			break
		}
	}

	if siteIndex == -1 {
		return fmt.Errorf("сайт с host '%s' не найден в конфигурации", host)
	}

	// 2. Удаляем папку сайта
	siteDir := filepath.Join("WebServer", "www", host)
	absSiteDir, err := filepath.Abs(siteDir)
	if err != nil {
		return fmt.Errorf("ошибка получения пути: %w", err)
	}

	if _, err := os.Stat(absSiteDir); err == nil {
		if err := os.RemoveAll(absSiteDir); err != nil {
			return fmt.Errorf("не удалось удалить папку сайта: %w", err)
		}
		tools.Logs_file(0, "SITES", fmt.Sprintf("🗑️ Удалена папка сайта: %s", siteDir), "logs_error.log", false)
	}

	// 3. Удаляем сертификаты
	if err := DeleteSiteCertificates(host); err != nil {
		// Логируем ошибку, но продолжаем удаление
		tools.Logs_file(1, "SITES", fmt.Sprintf("Ошибка удаления сертификатов: %v", err), "logs_error.log", false)
	}

	// 4. Удаляем из конфига
	config.ConfigData.Site_www = append(
		config.ConfigData.Site_www[:siteIndex],
		config.ConfigData.Site_www[siteIndex+1:]...,
	)

	// 5. Сохраняем конфиг
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("ошибка сохранения конфигурации: %w", err)
	}

	tools.Logs_file(0, "SITES", fmt.Sprintf("✅ Сайт '%s' полностью удалён", host), "logs_error.log", true)
	return nil
}

// generateTemplate генерирует шаблон для нового сайта
func generateTemplate(host, rootFile string) string {
	// Для всех типов файлов используем один HTML шаблон
	return strings.ReplaceAll(indexTemplate, "{{.Host}}", host)
}
