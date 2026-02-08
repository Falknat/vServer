package config

import (
	"encoding/json"
	"os"
	tools "vServer/Backend/tools"
)

var ConfigPath = "WebServer/config.json"

var ConfigData struct {
	Site_www      []Site_www      `json:"Site_www"`
	Soft_Settings Soft_Settings   `json:"Soft_Settings"`
	Proxy_Service []Proxy_Service `json:"Proxy_Service"`
}

type Site_www struct {
	Name              string   `json:"name"`
	Host              string   `json:"host"`
	Alias             []string `json:"alias"`
	Status            string   `json:"status"`
	Root_file         string   `json:"root_file"`
	Root_file_routing bool     `json:"root_file_routing"`
	AutoCreateSSL     bool     `json:"AutoCreateSSL"`
	Compression       *bool    `json:"Compression"`
}

func (s Site_www) IsCompressionEnabled() bool {
	if s.Compression == nil {
		return true
	}
	return *s.Compression
}

type Soft_Settings struct {
	Php_port      int    `json:"php_port"`
	Php_host      string `json:"php_host"`
	Mysql_port    int    `json:"mysql_port"`
	Mysql_host    string `json:"mysql_host"`
	Proxy_enabled bool   `json:"proxy_enabled"`
	ACME_enabled  bool   `json:"ACME_enabled"`
}

type Proxy_Service struct {
	Name            string `json:"Name"`
	Enable          bool   `json:"Enable"`
	ExternalDomain  string `json:"ExternalDomain"`
	LocalAddress    string `json:"LocalAddress"`
	LocalPort       string `json:"LocalPort"`
	ServiceHTTPSuse bool   `json:"ServiceHTTPSuse"`
	AutoHTTPS       bool   `json:"AutoHTTPS"`
	AutoCreateSSL   bool   `json:"AutoCreateSSL"`
	Compression     *bool  `json:"Compression"`
}

func (p Proxy_Service) IsCompressionEnabled() bool {
	if p.Compression == nil {
		return true
	}
	return *p.Compression
}

func LoadConfig() {

	data, err := os.ReadFile(ConfigPath)

	if err != nil {
		tools.Logs_file(1, "JSON", "Ошибка загрузки конфигурационного файла", "logs_error.log", true)
		return
	}

	err = json.Unmarshal(data, &ConfigData)
	if err != nil {
		tools.Logs_file(1, "JSON", "Ошибка парсинга конфигурационного файла", "logs_error.log", true)
	}

	// Миграция: добавляем новые поля если их нет
	migrateConfig(data)

	println()

}

// migrateConfig проверяет и добавляет новые поля в конфиг
func migrateConfig(originalData []byte) {
	needsSave := false

	// Парсим оригинальный JSON как map для проверки наличия полей
	var rawConfig map[string]json.RawMessage
	if err := json.Unmarshal(originalData, &rawConfig); err != nil {
		return
	}

	// Проверяем Site_www
	if rawSites, ok := rawConfig["Site_www"]; ok {
		var sites []map[string]interface{}
		if err := json.Unmarshal(rawSites, &sites); err == nil {
			for _, site := range sites {
				if _, exists := site["AutoCreateSSL"]; !exists {
					needsSave = true
					break
				}
			}
		}
	}

	// Проверяем Site_www на наличие Compression
	if rawSites, ok := rawConfig["Site_www"]; ok {
		var sites2 []map[string]interface{}
		if err := json.Unmarshal(rawSites, &sites2); err == nil {
			for i, site := range sites2 {
				if _, exists := site["Compression"]; !exists {
					needsSave = true
					compressionTrue := true
					ConfigData.Site_www[i].Compression = &compressionTrue
				}
			}
		}
	}

	// Проверяем Proxy_Service
	if rawProxies, ok := rawConfig["Proxy_Service"]; ok {
		var proxies []map[string]interface{}
		if err := json.Unmarshal(rawProxies, &proxies); err == nil {
			for _, proxy := range proxies {
				if _, exists := proxy["AutoCreateSSL"]; !exists {
					needsSave = true
					break
				}
			}
		}
	}

	// Проверяем Proxy_Service на наличие Compression
	if rawProxies, ok := rawConfig["Proxy_Service"]; ok {
		var proxies2 []map[string]interface{}
		if err := json.Unmarshal(rawProxies, &proxies2); err == nil {
			for i, proxy := range proxies2 {
				if _, exists := proxy["Compression"]; !exists {
					needsSave = true
					compressionTrue := true
					ConfigData.Proxy_Service[i].Compression = &compressionTrue
				}
			}
		}
	}

	// Проверяем Soft_Settings на наличие ACME_enabled
	if rawSettings, ok := rawConfig["Soft_Settings"]; ok {
		var settings map[string]interface{}
		if err := json.Unmarshal(rawSettings, &settings); err == nil {
			if _, exists := settings["ACME_enabled"]; !exists {
				needsSave = true
			}
		}
	}

	// Если нужно обновить - сохраняем конфиг с новыми полями
	if needsSave {
		tools.Logs_file(0, "JSON", "🔄 Миграция конфига: добавляем новые поля", "logs_error.log", true)
		saveConfig()
	}
}

// saveConfig сохраняет текущий конфиг в файл
func saveConfig() {
	formattedJSON, err := json.MarshalIndent(ConfigData, "", "    ")
	if err != nil {
		tools.Logs_file(1, "JSON", "Ошибка форматирования конфига: "+err.Error(), "logs_error.log", true)
		return
	}

	err = os.WriteFile(ConfigPath, formattedJSON, 0644)
	if err != nil {
		tools.Logs_file(1, "JSON", "Ошибка сохранения конфига: "+err.Error(), "logs_error.log", true)
		return
	}

	tools.Logs_file(0, "JSON", "✅ Конфиг обновлён с новыми полями", "logs_error.log", true)
}
