package webserver

import (
	"net/http"
	tools "vServer/Backend/tools"
)

var httpServer *http.Server
var port_http string = "80"

// Запуск HTTP сервера
func StartHTTP() {

	if tools.Port_check("HTTP", "localhost", port_http) {
		return
	}

	// Создаем HTTP сервер
	httpServer = &http.Server{
		Addr:    ":" + port_http,
		Handler: nil,
	}

	tools.Logs_file(0, "HTTP ", "💻 HTTP сервер запущен на порту 80", "logs_http.log", true)

	if err := httpServer.ListenAndServe(); err != nil {
		// Игнорируем нормальную ошибку при остановке сервера
		if err.Error() != "http: Server closed" {
			tools.Logs_file(1, "HTTP", "❌ Ошибка запуска сервера: "+err.Error(), "logs_http.log", true)
		}
	}
}

// StopHTTPServer останавливает HTTP сервер
func StopHTTPServer() {
	if httpServer != nil {
		httpServer.Close()
		httpServer = nil
		tools.Logs_file(0, "HTTP", "HTTP сервер остановлен", "logs_http.log", true)
	}
}
