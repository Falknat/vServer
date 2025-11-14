package tools

import (
	"encoding/base64"
	"fmt"
	"net"
	"time"
)

// Время запуска сервера
var ServerStartTime time.Time

// isPortInUse проверяет, занят ли указанный порт
func Port_check(service string, host string, port string) bool {
	conn, err := net.DialTimeout("tcp", host+":"+port, time.Millisecond*300)
	if err != nil {
		return false // порт свободен
	}
	conn.Close()
	Logs_file(1, service, "⚠️ Порт "+port+" уже занят, сервис не запущен", "logs_error.log", true)
	return true // порт занят
}

// Управление временем работы сервера
func ServerUptime(action string, asSeconds ...bool) interface{} {
	switch action {
	case "start":
		// Инициализация времени запуска
		ServerStartTime = time.Now()
		return nil

	case "get":
		// Получить время работы
		if ServerStartTime.IsZero() {
			if len(asSeconds) > 0 && asSeconds[0] {
				return int64(0)
			}
			return "Сервер не запущен"
		}

		uptime := time.Since(ServerStartTime)

		// Возвращаем секунды
		if len(asSeconds) > 0 && asSeconds[0] {
			return int64(uptime.Seconds())
		}

		// Возвращаем читаемый формат
		days := int(uptime.Hours()) / 24
		hours := int(uptime.Hours()) % 24
		minutes := int(uptime.Minutes()) % 60
		seconds := int(uptime.Seconds()) % 60

		if days > 0 {
			return fmt.Sprintf("%dд %dч %dм", days, hours, minutes)
		} else if hours > 0 {
			return fmt.Sprintf("%dч %dм", hours, minutes)
		} else if minutes > 0 {
			return fmt.Sprintf("%dм", minutes)
		} else {
			return fmt.Sprintf("%dс", seconds)
		}

	default:
		return "Неизвестное действие"
	}
}

func Error_check(err error, message string) bool {
	if err != nil {
		fmt.Printf("Ошибка: %v\n", message)
		return false
	}

	return true
}

// DecodeBase64 декодирует строку из base64
func DecodeBase64(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования base64: %w", err)
	}
	return decoded, nil
}
