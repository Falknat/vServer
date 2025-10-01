package command

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"
	"vServer/Backend/config"
	"vServer/Backend/tools"
)

var entries []os.DirEntry

func path_www() {

	wwwPath, err := tools.AbsPath("WebServer/www")
	tools.Error_check(err, "Ошибка получения пути")

	entries, err = os.ReadDir(wwwPath)
	tools.Error_check(err, "Ошибка чтения директории")

}

func site_type(entry os.DirEntry) string {

	certPath, err := tools.AbsPath("WebServer/cert/" + entry.Name())
	if err == nil {
		if _, err := os.Stat(certPath); err == nil {
			return "https"
		}
	}
	return "http"
}

func site_alliace(siteName string) []string {
	// Получаем абсолютный путь к config.json
	configPath, err := tools.AbsPath("WebServer/config.json")
	tools.Error_check(err, "Ошибка получения пути к config.json")

	// Читаем содержимое config.json
	configData, err := os.ReadFile(configPath)
	tools.Error_check(err, "Ошибка чтения config.json")

	// Структура для парсинга Site_www
	type SiteConfig struct {
		Name   string   `json:"name"`
		Host   string   `json:"host"`
		Alias  []string `json:"alias"`
		Status string   `json:"status"`
	}
	type Config struct {
		SiteWWW []SiteConfig `json:"Site_www"`
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	tools.Error_check(err, "Ошибка парсинга config.json")

	// Ищем алиасы для конкретного сайта
	for _, site := range config.SiteWWW {
		if site.Host == siteName {
			return site.Alias
		}
	}

	// Возвращаем пустой массив если сайт не найден
	return []string{}
}

func site_status(siteName string) string {
	configPath := "WebServer/config.json"

	// Читаем конфигурационный файл
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "error"
	}

	// Парсим JSON
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return "error"
	}

	// Получаем список сайтов
	siteWww, ok := config["Site_www"].([]interface{})
	if !ok {
		return "error"
	}

	// Ищем сайт по host
	for _, siteInterface := range siteWww {
		site, ok := siteInterface.(map[string]interface{})
		if !ok {
			continue
		}

		// Проверяем только по host (имя папки)
		if host, ok := site["host"].(string); ok && host == siteName {
			if status, ok := site["status"].(string); ok {
				return status
			}
		}
	}

	return "inactive"
}

func SiteList(w http.ResponseWriter, r *http.Request, path string) bool {

	switch path {

	case "/service/Site_List":
		sites := []map[string]interface{}{}

		path_www()

		for _, entry := range entries {
			if entry.IsDir() {
				site := map[string]interface{}{
					"host":    entry.Name(),
					"type":    site_type(entry),
					"aliases": site_alliace(entry.Name()),
					"status":  site_status(entry.Name()),
				}
				sites = append(sites, site)
			}
		}

		metrics := map[string]interface{}{
			"sites": sites,
		}

		data, _ := json.MarshalIndent(metrics, "", "  ")
		w.Write(data)

		return true

	}

	return false

}

func addSiteToConfig(siteName string) error {
	// Получаем абсолютный путь к config.json
	configPath, err := tools.AbsPath("WebServer/config.json")
	if err != nil {
		return err
	}

	// Читаем содержимое config.json
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Парсим как общий JSON объект
	var config map[string]interface{}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return err
	}

	// Получаем массив сайтов
	siteWWW, ok := config["Site_www"].([]interface{})
	if !ok {
		siteWWW = []interface{}{}
	}

	// Создаем новый сайт в том же формате что уже есть
	newSite := map[string]interface{}{
		"name":   siteName,
		"host":   siteName,
		"alias":  []string{""},
		"status": "active",
	}

	// Добавляем новый сайт в массив
	config["Site_www"] = append(siteWWW, newSite)

	// Сохраняем обновленный конфиг
	updatedData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	// Делаем массивы алиасов в одну строку
	re := regexp.MustCompile(`"alias": \[\s+"([^"]*?)"\s+\]`)
	compactData := re.ReplaceAll(updatedData, []byte(`"alias": ["$1"]`))

	// Исправляем отступы после Site_www
	dataStr := string(compactData)
	dataStr = strings.ReplaceAll(dataStr, `    ],
    "Soft_Settings"`, `    ],

    "Soft_Settings"`)
	compactData = []byte(dataStr)

	err = os.WriteFile(configPath, compactData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func site_add(w http.ResponseWriter, path string) bool {

	// URL параметры: /service/Site_Add/sitename
	if strings.HasPrefix(path, "/service/Site_Add/") {
		siteName := strings.TrimPrefix(path, "/service/Site_Add/")

		if siteName == "" {
			w.WriteHeader(400)
			w.Write([]byte(`{"status":"error","message":"Не указано имя сайта в URL"}`))
			return true
		}

		wwwPath := "WebServer/www/" + siteName

		// Проверяем существует ли уже такой сайт
		if _, err := os.Stat(wwwPath); err == nil {
			w.WriteHeader(409) // Conflict
			w.Write([]byte(`{"status":"error","message":"Сайт ` + siteName + ` уже существует"}`))
			return true
		}

		err := os.MkdirAll(wwwPath, 0755)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`{"status":"error","message":"Ошибка создания папки сайта"}`))
			return true
		}

		// Добавляем сайт в config.json
		err = addSiteToConfig(siteName)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`{"status":"error","message":"Ошибка добавления в конфигурацию: ` + err.Error() + `"}`))
			return true
		}

		// Создаем папку public_www
		publicWwwPath := wwwPath + "/public_www"
		err = os.MkdirAll(publicWwwPath, 0755)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`{"status":"error","message":"Ошибка создания папки public_www"}`))
			return true
		}

		indexFilePath := wwwPath + "/public_www/index.html"
		indexContent := "Привет друг! Твой сайт создан!"
		err = os.WriteFile(indexFilePath, []byte(indexContent), 0644)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`{"status":"error","message":"Ошибка создания index.html: ` + err.Error() + `"}`))
			return true
		}

		w.Write([]byte(`{"status":"ok","message":"Сайт ` + siteName + ` успешно создан и добавлен в конфигурацию"}`))
		config.LoadConfig()
		return true
	}

	return false
}
