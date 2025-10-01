package webserver

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	tools "vServer/Backend/tools"
)

// ProxyConfig хранит конфигурацию для прокси
type ProxyConfig struct {
	ExternalDomain string
	LocalAddress   string
	LocalPort      string
	UseHTTPS       bool
}

var (
	proxyConfigs  = make(map[int]*ProxyConfig)
	configMutex   sync.RWMutex
	configsLoaded = false
)

// InitProxyConfigs инициализирует конфигурации прокси один раз при старте
func InitProxyConfigs() {
	configMutex.Lock()
	defer configMutex.Unlock()

	if configsLoaded {
		return
	}

	// Конфигурация 1
	config1 := &ProxyConfig{
		ExternalDomain: "git.voxsel.ru",
		LocalAddress:   "127.0.0.1",
		LocalPort:      "3333",
		UseHTTPS:       false, // Локальный сервис работает по HTTP
	}
	proxyConfigs[1] = config1

	// Конфигурация 2
	config2 := &ProxyConfig{
		ExternalDomain: "localhost",
		LocalAddress:   "127.0.0.1",
		LocalPort:      "8000",
		UseHTTPS:       false, // Локальный сервис работает по HTTP
	}
	proxyConfigs[2] = config2

	configsLoaded = true
}

func StartHandlerProxy(w http.ResponseWriter, r *http.Request) (valid bool) {
	valid = false

	// Инициализируем конфигурации если еще не сделано
	if !configsLoaded {
		InitProxyConfigs()
	}

	configMutex.RLock()
	defer configMutex.RUnlock()

	// Выбираем конфигурацию (пока используем 1)
	config := proxyConfigs[1]
	if config == nil {
		return false
	}

	if r.Host == config.ExternalDomain {
		valid = true

		// Определяем протокол для локального соединения
		protocol := "http"
		if config.UseHTTPS {
			protocol = "https"
		}

		// Проксирование на локальный адрес
		proxyURL := protocol + "://" + config.LocalAddress + ":" + config.LocalPort + r.URL.RequestURI()
		proxyReq, err := http.NewRequest(r.Method, proxyURL, r.Body)
		if err != nil {
			http.Error(w, "Ошибка создания прокси-запроса", http.StatusInternalServerError)
			return
		}

		// Копируем ВСЕ заголовки без изменений (кроме технических)
		for name, values := range r.Header {
			// Пропускаем только технические заголовки HTTP/1.1
			lowerName := strings.ToLower(name)
			if lowerName == "connection" || lowerName == "upgrade" ||
				lowerName == "proxy-connection" || lowerName == "te" ||
				lowerName == "trailers" || lowerName == "transfer-encoding" {
				continue
			}

			// Копируем заголовок как есть
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}

		// Прозрачная передача - никаких дополнительных заголовков
		// Все заголовки уже скопированы выше "как есть"

		// Выполняем прокси-запрос
		client := &http.Client{
			// Отключаем автоматическое следование редиректам для корректной работы с авторизацией
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		// Для HTTPS соединений настраиваем TLS (если понадобится)
		if config.UseHTTPS {
			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // Простая настройка для внутренних соединений
				},
			}
		}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, "Ошибка прокси-запроса", http.StatusBadGateway)
			tools.Logs_file(1, "PROXY", "Ошибка прокси-запроса: "+err.Error(), "logs_proxy.log", true)
			return
		}
		defer resp.Body.Close()

		// Прозрачно копируем ВСЕ заголовки ответа без изменений
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// Устанавливаем статус код
		w.WriteHeader(resp.StatusCode)

		// Копируем тело ответа
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("Ошибка копирования тела ответа: %v", err)
		}

		return valid

	} else {
		return valid
	}
}
