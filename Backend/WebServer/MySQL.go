package webserver

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
	config "vServer/Backend/config"
	tools "vServer/Backend/tools"
)

var mysqlProcess *exec.Cmd
var mysql_status bool = false
var mysql_secure bool = false

var mysqldPath string
var configPath string
var dataDirAbs string
var binDirAbs string
var binPathAbs string

var mysql_port int
var mysql_ip string

var console_mysql bool = false

func AbsPathMySQL() {

	var err error

	mysqldPath, err = tools.AbsPath(filepath.Join("WebServer/soft/MySQL/bin", "mysqld.exe"))
	tools.CheckError(err)

	configPath, err = tools.AbsPath("WebServer/soft/MySQL/my.ini")
	tools.CheckError(err)

	dataDirAbs, err = tools.AbsPath("WebServer/soft/MySQL/bin/data")
	tools.CheckError(err)

	binDirAbs, err = tools.AbsPath("WebServer/soft/MySQL/bin")
	tools.CheckError(err)

	binPathAbs, err = tools.AbsPath("WebServer/soft/MySQL/bin")
	tools.CheckError(err)

}

// config_patch возвращает путь к mysqld, аргументы и бинарную директорию
func config_patch(secures bool) (string, []string, string) {

	// Получаем абсолютные пути
	AbsPathMySQL()

	// Объявляем args на уровне функции
	var args []string

	if secures {

		args = []string{
			"--defaults-file=" + configPath,
			"--datadir=" + dataDirAbs,
			"--shared-memory",
			"--skip-grant-tables",
			"--console",
		}

	} else {

		args = []string{
			"--defaults-file=" + configPath,
			"--port=" + fmt.Sprintf("%d", mysql_port),
			"--bind-address=" + mysql_ip,
			"--datadir=" + dataDirAbs,
			"--console",
		}

	}

	return mysqldPath, args, binDirAbs
}

// StartMySQLServer запускает MySQL сервер
func StartMySQLServer(secure bool) {

	mysql_port = config.ConfigData.Soft_Settings.Mysql_port
	mysql_ip = config.ConfigData.Soft_Settings.Mysql_host

	if tools.Port_check("MySQL", mysql_ip, strconv.Itoa(mysql_port)) {
		return
	}

	if mysql_status {
		tools.Logs_file(1, "MySQL", "Сервер MySQL уже запущен", "logs_mysql.log", true)
		return
	}

	if false {
		return
	}

	// Настройка режима
	mysql_secure = secure
	mysqldPath, args, binDirAbs := config_patch(secure)

	// Выбор сообщения
	if secure {
		tools.Logs_file(0, "MySQL", "Запуск сервера MySQL в режиме безопасности", "logs_mysql.log", true)
	} else {
		tools.Logs_file(0, "MySQL", "Запуск сервера MySQL в обычном режиме", "logs_mysql.log", true)
	}

	// Общая логика запуска
	mysqlProcess = exec.Command(mysqldPath, args...)
	mysqlProcess.Dir = binDirAbs
	tools.Logs_console(mysqlProcess, console_mysql)

	tools.Logs_file(0, "MySQL", fmt.Sprintf("Сервер MySQL запущен на %s:%d", mysql_ip, mysql_port), "logs_mysql.log", true)

	mysql_status = true

}

// StopMySQLServer останавливает MySQL сервер
func StopMySQLServer() {

	if mysql_status {

		cmd := exec.Command("taskkill", "/F", "/IM", "mysqld.exe")

		err := cmd.Run()
		tools.CheckError(err)

		tools.Logs_file(0, "MySQL", "Сервер MySQL остановлен", "logs_mysql.log", true)
		mysql_status = false

	} else {

		tools.Logs_file(1, "MySQL", "Сервер MySQL уже остановлен", "logs_mysql.log", true)

	}

}

func ResetPasswordMySQL() {

	NewPasswordMySQL := "root"

	StopMySQLServer()
	time.Sleep(2 * time.Second)
	mysql_secure = true
	StartMySQLServer(true)
	time.Sleep(2 * time.Second)
	query := "FLUSH PRIVILEGES; ALTER USER 'root'@'%' IDENTIFIED BY '" + NewPasswordMySQL + "';"
	СheckMySQLPassword(query)
	tools.Logs_file(0, "MySQL", "Новый пароль: "+NewPasswordMySQL, "logs_mysql.log", true)
	println()
	StopMySQLServer()
	StartMySQLServer(false)

}

// СheckMySQLPassword проверяет пароль для MySQL
func СheckMySQLPassword(query string) {

	AbsPathMySQL()

	if mysql_secure {

		// В безопасном режиме подключаемся без пароля
		cmd := exec.Command(filepath.Join(binPathAbs, "mysql.exe"), "-u", "root", "-pRoot", "-e", query)
		cmd.Dir = binPathAbs

		// Захватываем вывод для логирования
		err := tools.Logs_console(cmd, false)

		if err != nil {
			tools.Logs_file(1, "MySQL", "Вывод MySQL (stdout/stderr):", "logs_mysql.log", true)
		} else {
			tools.Logs_file(0, "MySQL", "Команда выполнена успешно", "logs_mysql.log", true)
		}

	}

}
