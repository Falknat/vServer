package services

import (
	"fmt"
	webserver "vServer/Backend/WebServer"
	config "vServer/Backend/config"
)

func GetAllServicesStatus() AllServicesStatus {
	return AllServicesStatus{
		HTTP:  getHTTPStatus(),
		HTTPS: getHTTPSStatus(),
		MySQL: getMySQLStatus(),
		PHP:   getPHPStatus(),
		Proxy: getProxyStatus(),
	}
}

func getHTTPStatus() ServiceStatus {
	// Используем внутренний статус вместо TCP проверки
	return ServiceStatus{
		Name:   "HTTP",
		Status: webserver.GetHTTPStatus(),
		Port:   "80",
		Info:   "",
	}
}

func getHTTPSStatus() ServiceStatus {
	// Используем внутренний статус вместо TCP проверки
	return ServiceStatus{
		Name:   "HTTPS",
		Status: webserver.GetHTTPSStatus(),
		Port:   "443",
		Info:   "",
	}
}

func getMySQLStatus() ServiceStatus {
	port := fmt.Sprintf("%d", config.ConfigData.Soft_Settings.Mysql_port)

	// Используем внутренний статус вместо TCP проверки
	// чтобы не вызывать connect_errors в MySQL
	return ServiceStatus{
		Name:   "MySQL",
		Status: webserver.GetMySQLStatus(),
		Port:   port,
		Info:   "",
	}
}

func getPHPStatus() ServiceStatus {
	basePort := config.ConfigData.Soft_Settings.Php_port

	// Диапазон портов для 4 воркеров
	portRange := fmt.Sprintf("%d-%d", basePort, basePort+3)

	// Используем внутренний статус вместо TCP проверки
	return ServiceStatus{
		Name:   "PHP",
		Status: webserver.GetPHPStatus(),
		Port:   portRange,
		Info:   "",
	}
}

func getProxyStatus() ServiceStatus {
	activeCount := 0
	totalCount := len(config.ConfigData.Proxy_Service)

	for _, proxy := range config.ConfigData.Proxy_Service {
		if proxy.Enable {
			activeCount++
		}
	}

	info := fmt.Sprintf("%d из %d", activeCount, totalCount)

	// Проверяем глобальный флаг и статус HTTP/HTTPS
	proxyEnabled := config.ConfigData.Soft_Settings.Proxy_enabled
	httpRunning := webserver.GetHTTPStatus()
	httpsRunning := webserver.GetHTTPSStatus()

	status := proxyEnabled && (httpRunning || httpsRunning)

	return ServiceStatus{
		Name:   "Proxy",
		Status: status,
		Port:   "-",
		Info:   info,
	}
}
