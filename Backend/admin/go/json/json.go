package json

import "encoding/json"

// Данные серверов
var ServerStatus = []map[string]interface{}{
	{"NameService": "HTTP Server", "Port": 80, "Status": "stopped"},
	{"NameService": "HTTPS Server", "Port": 443, "Status": "stopped"},
	{"NameService": "PHP Server", "Port": 9000, "Status": "stopped"},
	{"NameService": "MySQL Server", "Port": 3306, "Status": "stopped"},
}

// Данные меню
var MenuData = []map[string]interface{}{
	{"name": "Dashboard", "icon": "🏠", "url": "#dashboard", "active": true},
	{"name": "Серверы", "icon": "🖥️", "url": "#servers", "active": false},
	{"name": "Сайты", "icon": "🌐", "url": "#sites", "active": false},
	{"name": "SSL Сертификаты", "icon": "🔒", "url": "#certificates", "active": false},
	{"name": "Файловый менеджер", "icon": "📁", "url": "#files", "active": false},
	{"name": "Базы данных", "icon": "🗄️", "url": "#databases", "active": false},
	{"name": "Логи", "icon": "📋", "url": "#logs", "active": false},
	{"name": "Настройки", "icon": "⚙️", "url": "#settings", "active": false},
}

// Функция обновления статуса сервера
func UpdateServerStatus(serviceName, status string) {
	for i := range ServerStatus {
		if ServerStatus[i]["NameService"] == serviceName {
			ServerStatus[i]["Status"] = status
			break
		}
	}
}

// Получить JSON серверов
func GetServerStatusJSON() []byte {
	data, _ := json.Marshal(ServerStatus)
	return data
}

// Получить JSON меню
func GetMenuJSON() []byte {
	data, _ := json.Marshal(MenuData)
	return data
}
