package webserver

import (
	"fmt"
	"os"
	"time"
	admin "vServer/Backend/admin"
	config "vServer/Backend/config"
	tools "vServer/Backend/tools"
)

var Secure_post bool = false

func CommandListener() {

	fmt.Println("Введите help для получения списка команд")
	fmt.Println("")

	for {
		var cmd string

		fmt.Print(tools.Color(" > ", tools.Оранжевый))

		fmt.Scanln(&cmd)

		switch cmd {
		case "help":

			fmt.Println(" ------------------------------------------")
			fmt.Println(" 1: mysql_stop - Остановить MySQL")
			fmt.Println(" 2: mysql_start - Запустить MySQL")
			fmt.Println(" 3: mysql_pass - Сбросить пароль MySQL")
			fmt.Println(" 4: clear - Очистить консоль")
			fmt.Println(" 5: cert_reload - Перезагрузить SSL сертификаты")
			fmt.Println(" 6: admin_toggle - Переключить режим админки (embed/файловая система)")
			fmt.Println(" 7: config_reload - Перезагрузить конфигурацию")
			fmt.Println(" 8: restart - Перезапустить сервер")
			fmt.Println(" 9: php_console - Открыть PHP консоль")
			fmt.Println(" 10: exit - выйти из программы")
			fmt.Println(" ------------------------------------------")
			fmt.Println("")

		case "mysql_stop":
			StopMySQLServer()

		case "mysql_start":
			StartMySQLServer(false)

		case "mysql_pass":
			ResetPasswordMySQL()

		case "clear":
			ClearConsole()

		case "cert_reload":
			ReloadCertificates()

		case "admin_toggle":
			AdminToggle()

		case "config_reload":
			ConfigReload()

		case "restart":
			RestartServer()

		case "time_run":
			fmt.Println(tools.ServerUptime("get"))

		case "secure_post":

			if Secure_post {
				Secure_post = false
				fmt.Println("Secure post is disabled")
			} else {
				Secure_post = true
				fmt.Println("Secure post is enabled")
			}

		case "exit":
			fmt.Println("Завершение...")
			os.Exit(0)

		default:

			fmt.Println(" Неизвестная команда. Введите 'help' для получения списка команд")
			fmt.Println("")

		}
	}
}

func RestartServer() {
	fmt.Println("")
	fmt.Println("⏹️ Перезагрузка сервера...")

	// Останавливаем все сервисы
	fmt.Println("⏹️ Останавливаем сервисы...")
	fmt.Println("")

	// Останавливаем HTTP/HTTPS серверы
	StopHTTPServer()
	StopHTTPSServer()

	// Останавливаем MySQL
	if mysql_status {
		StopMySQLServer()
		time.Sleep(1 * time.Second)
	}

	// Останавливаем PHP
	PHP_Stop()
	time.Sleep(1 * time.Second)

	fmt.Println("")
	fmt.Println("✅ Все сервисы остановлены")

	// Перезагружаем конфигурацию
	fmt.Println("📋 Перезагружаем конфигурацию...")
	fmt.Println("")
	config.LoadConfig()

	// Запускаем сервисы заново
	fmt.Println("🚀 Запускаем сервисы...")
	fmt.Println("")

	// Запускаем HTTP/HTTPS серверы
	go StartHTTP()
	time.Sleep(100 * time.Millisecond)
	go StartHTTPS()
	time.Sleep(100 * time.Millisecond)

	// Запускаем PHP
	PHP_Start()
	time.Sleep(100 * time.Millisecond)

	// Запускаем MySQL
	StartMySQLServer(false)
	time.Sleep(100 * time.Millisecond)

	fmt.Println("✅ Сервер успешно перезагружен!")
	fmt.Println("")
}

func ClearConsole() {
	// Очищаем консоль, но сохраняем первые три строки
	fmt.Print("\033[H\033[2J") // ANSI escape code для очистки экрана

	println("")
	println(tools.Color("vServer", tools.Жёлтый) + tools.Color(" 1.0.0", tools.Голубой))
	println(tools.Color("Автор: ", tools.Зелёный) + tools.Color("Суманеев Роман (c) 2025", tools.Голубой))
	println(tools.Color("Официальный сайт: ", tools.Зелёный) + tools.Color("https://voxsel.ru", tools.Голубой))

	println("")

	// Восстанавливаем первые три строки
	fmt.Println("Введите help для получения списка команд")
	fmt.Println("")
}

// Переключает режим админки между embed и файловой системой
func AdminToggle() {
	fmt.Println("")

	if admin.UseEmbedded {
		// Переключаем на файловую систему
		admin.UseEmbedded = false
		fmt.Println("🔄 Режим изменен: Embedded → Файловая система")
		fmt.Println("✅ Админка переключена на файловую систему")
		fmt.Println("📁 Файлы будут загружаться с диска из Backend/admin/html/")
		fmt.Println("💡 Теперь можно редактировать файлы и изменения будут видны сразу")
	} else {
		// Переключаем обратно на embedded
		admin.UseEmbedded = true
		fmt.Println("🔄 Режим изменен: Файловая система → Embedded")
		fmt.Println("✅ Админка переключена на embedded режим")
		fmt.Println("📦 Файлы загружаются из встроенных ресурсов")
		fmt.Println("🚀 Быстрая загрузка, но изменения требуют перекомпиляции")
	}

	fmt.Println("")
}

// Перезагружает конфигурацию без перезапуска сервисов
func ConfigReload() {
	fmt.Println("")
	fmt.Println("📋 Перезагружаем конфигурацию...")

	// Загружаем новую конфигурацию
	config.LoadConfig()

	fmt.Println("✅ Конфигурация успешно перезагружена!")
	fmt.Println("💡 Изменения применятся к новым запросам")
	fmt.Println("")
}
