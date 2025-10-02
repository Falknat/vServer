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
}

type Soft_Settings struct {
	Php_port   int    `json:"php_port"`
	Php_host   string `json:"php_host"`
	Mysql_port int    `json:"mysql_port"`
	Mysql_host string `json:"mysql_host"`
	Admin_port string `json:"admin_port"`
	Admin_host string `json:"admin_host"`
}

type Proxy_Service struct {
	Enable          bool   `json:"Enable"`
	ExternalDomain  string `json:"ExternalDomain"`
	LocalAddress    string `json:"LocalAddress"`
	LocalPort       string `json:"LocalPort"`
	ServiceHTTPSuse bool   `json:"ServiceHTTPSuse"`
	AutoHTTPS       bool   `json:"AutoHTTPS"`
}

func LoadConfig() {

	data, err := os.ReadFile(ConfigPath)

	if err != nil {
		tools.Logs_file(0, "JSON", "Ошибка загрузки конфигурационного файла", "logs_config.log", true)
	} else {
		tools.Logs_file(0, "JSON", "config.json успешно загружен", "logs_config.log", true)
	}

	err = json.Unmarshal(data, &ConfigData)
	if err != nil {
		tools.Logs_file(0, "JSON", "Ошибка парсинга конфигурационного файла", "logs_config.log", true)
	} else {
		tools.Logs_file(0, "JSON", "config.json успешно прочитан", "logs_config.log", true)
	}

	println()

}
