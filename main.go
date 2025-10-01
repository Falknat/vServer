package main

import (
	"fmt"
	"time"
	webserver "vServer/Backend/WebServer"
	admin "vServer/Backend/admin/go"
	json_admin "vServer/Backend/admin/go/json"
	config "vServer/Backend/config"
	tools "vServer/Backend/tools"
)

func main() {

	if !tools.CheckSingleInstance() {
		println("")
		println(tools.Color("❌ ОШИБКА:", tools.Красный) + " vServer уже запущен!")
		println(tools.Color("💡 Подсказка:", tools.Жёлтый) + " Завершите уже запущенный процесс перед запуском нового.")
		println("")
		println("Нажмите Enter для завершения...")
		fmt.Scanln()
		return
	}

	// Освобождаем мьютекс при выходе (опционально, так как Windows сама освободит)
	defer tools.ReleaseMutex()

	println("")
	println(tools.Color("vServer", tools.Жёлтый) + tools.Color(" 1.0.0", tools.Голубой))
	println(tools.Color("Автор: ", tools.Зелёный) + tools.Color("Суманеев Роман (c) 2025", tools.Голубой))
	println(tools.Color("Официальный сайт: ", tools.Зелёный) + tools.Color("https://voxsel.ru", tools.Голубой))

	println("")
	println("🚀 Запуск vServer...")
	println("📁 Файлы сайта будут обслуживаться из папки 'www'")
	println("")
	println("⏳ Запуск сервисов...")
	println("")

	// Инициализируем время запуска сервера
	tools.ServerUptime("start")

	config.LoadConfig()
	time.Sleep(50 * time.Millisecond)

	webserver.StartHandler()
	time.Sleep(50 * time.Millisecond)

	// Запускаем серверы в горутинах
	go admin.StartAdmin()
	time.Sleep(50 * time.Millisecond)

	webserver.Cert_start()
	time.Sleep(50 * time.Millisecond)

	go webserver.StartHTTPS()
	json_admin.UpdateServerStatus("HTTPS Server", "running")
	time.Sleep(50 * time.Millisecond)

	go webserver.StartHTTP()
	json_admin.UpdateServerStatus("HTTP Server", "running")
	time.Sleep(50 * time.Millisecond)

	webserver.PHP_Start()
	json_admin.UpdateServerStatus("PHP Server", "running")
	time.Sleep(50 * time.Millisecond)

	webserver.StartMySQLServer(false)
	json_admin.UpdateServerStatus("MySQL Server", "running")
	time.Sleep(50 * time.Millisecond)

	println("")
	webserver.CommandListener()

}
