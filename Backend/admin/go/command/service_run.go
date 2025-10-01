package command

import (
	"net/http"
	webserver "vServer/Backend/WebServer"
	json "vServer/Backend/admin/go/json"
)

// Обработчик команд управления серверами
func Service_Run(w http.ResponseWriter, r *http.Request, path string) bool {

	switch path {
	case "/service/MySql_Stop":
		webserver.StopMySQLServer()
		json.UpdateServerStatus("MySQL Server", "stopped")
		return true

	case "/service/MySql_Start":
		webserver.StartMySQLServer(false)
		json.UpdateServerStatus("MySQL Server", "running")
		return true

	case "/service/Http_Stop":
		webserver.StopHTTPServer()
		json.UpdateServerStatus("HTTP Server", "stopped")
		return true

	case "/service/Http_Start":
		go webserver.StartHTTP()
		json.UpdateServerStatus("HTTP Server", "running")
		return true

	case "/service/Https_Stop":
		webserver.StopHTTPSServer()
		json.UpdateServerStatus("HTTPS Server", "stopped")
		return true

	case "/service/Https_Start":
		go webserver.StartHTTPS()
		json.UpdateServerStatus("HTTPS Server", "running")
		return true

	case "/service/Php_Start":
		webserver.PHP_Start()
		json.UpdateServerStatus("PHP Server", "running")
		return true

	case "/service/Php_Stop":
		webserver.PHP_Stop()
		json.UpdateServerStatus("PHP Server", "stopped")
		return true

	default:
		http.NotFound(w, r)
		return false // Команда не найдена
	}
}
