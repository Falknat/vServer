package admin

import (
	"net/http"
	command "vServer/Backend/admin/go/command"
	config "vServer/Backend/config"
	tools "vServer/Backend/tools"
)

var adminServer *http.Server

// Запуск Admin сервера
func StartAdmin() {

	// Получаем значения из конфига во время выполнения
	port_admin := config.ConfigData.Soft_Settings.Admin_port
	host_admin := config.ConfigData.Soft_Settings.Admin_host

	if tools.Port_check("ADMIN", host_admin, port_admin) {
		return
	}

	// Создаем оптимизированный мультиплексор для админ сервера
	mux := http.NewServeMux()

	// Регистрируем специализированные обработчики (быстрая маршрутизация)
	mux.HandleFunc("/api/", command.ApiHandler)         // API эндпоинты
	mux.HandleFunc("/json/", command.JsonHandler)       // JSON данные
	mux.HandleFunc("/service/", command.ServiceHandler) // Сервисные команды POST
	mux.HandleFunc("/", command.StaticHandler)          // Статические файлы

	// Создаем Admin сервер (только localhost для безопасности)
	adminServer = &http.Server{
		Addr:    host_admin + ":" + port_admin,
		Handler: mux,
	}

	tools.Logs_file(0, "ADMIN", "🛠️ Admin панель запущена на порту "+port_admin, "logs_http.log", true)

	if err := adminServer.ListenAndServe(); err != nil {
		// Игнорируем нормальную ошибку при остановке сервера
		if err.Error() != "http: Server closed" {
			tools.Logs_file(1, "ADMIN", "❌ Ошибка запуска админ сервера: "+err.Error(), "logs_http.log", true)
		}
	}
}
