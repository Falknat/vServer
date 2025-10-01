package tools

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

const (
	Красный     = "\033[31m"
	Зелёный     = "\033[32m"
	Жёлтый      = "\033[33m"
	Синий       = "\033[34m"
	Голубой     = "\033[36m"
	Фиолетовый  = "\033[35m"
	Белый       = "\033[37m"
	Серый       = "\033[90m"
	Оранжевый   = "\033[38;5;208m"
	Сброс_Цвета = "\033[0m"
)

// Функция окрашивания текста
func Color(text, ansi string) string {
	return ansi + text + Сброс_Цвета
}

// Функция для удаления ANSI-кодов из строки
func RemoveAnsiCodes(text string) string {
	// Регулярное выражение для удаления ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(text, "")
}

// Логирование в файл
/*
	type_log:
		0 - INFO
		1 - ERROR
		2 - WARNING
*/
func Logs_file(type_log int, service string, message string, log_file string, console bool) {

	color_data := ""

	service_str := Color(" ["+service+"] ", Жёлтый)
	type_log_str := Color(" [INFO] ", Голубой)
	log_files := log_file

	switch type_log {
	case 0:
		type_log_str = Color(" [-INFOS-]", Голубой)
	case 1:
		type_log_str = Color(" [-ERROR-]", Красный)
	case 2:
		type_log_str = Color(" [WARNING]", Жёлтый)
	}

	if type_log == 1 {
		color_data = Красный
	} else {
		color_data = Зелёный
	}

	if console {
		// Очищаем текущую строку (стираем промпт >) и выводим лог с новой строки
		fmt.Print("\r\033[K")
		fmt.Println(Color(time.Now().Format("2006-01-02 15:04:05")+type_log_str+service_str+message, color_data))
	}

	// Создаем текст с цветами, затем удаляем ANSI-коды для файла
	colored_text := time.Now().Format("2006-01-02 15:04:05") + type_log_str + service_str + message
	text := RemoveAnsiCodes(colored_text) + "\n"

	// Открываем файл для дозаписи, создаём если нет, права на запись.
	file, err := os.OpenFile("WebServer/tools/logs/"+log_files, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Пишем строку в файл
	if _, err := file.WriteString(text); err != nil {
		log.Fatal(err)
	}

}
