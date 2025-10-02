package webserver

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"vServer/Backend/config"
	tools "vServer/Backend/tools"
)

var (
	configMutex sync.RWMutex
)

func StartHandlerProxy(w http.ResponseWriter, r *http.Request) (valid bool) {
	valid = false

	configMutex.RLock()
	defer configMutex.RUnlock()

	// Проходим по всем прокси конфигурациям
	for _, proxyConfig := range config.ConfigData.Proxy_Service {
		// Пропускаем отключенные прокси
		if !proxyConfig.Enable {
			continue
		}

		// Проверяем совпадение домена
		if r.Host != proxyConfig.ExternalDomain {
			continue
		}

		valid = true

		// Проверяем vAccess для прокси
		accessAllowed, errorPage := CheckProxyVAccess(r.URL.Path, proxyConfig.ExternalDomain, r)
		if !accessAllowed {
			// Доступ запрещён - обрабатываем страницу ошибки
			HandleProxyVAccessError(w, r, errorPage)
			return valid
		}

		// Проверяем AutoHTTPS - редирект с HTTP на HTTPS
		https_check := !(r.TLS == nil)
		if !https_check && proxyConfig.AutoHTTPS {
			// Перенаправляем на HTTPS
			httpsURL := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
			tools.Logs_file(0, "P-HTTP", "🔀 IP клиента: "+r.RemoteAddr+" Редирект HTTP → HTTPS: "+r.Host+r.URL.Path, "logs_http.log", false)
			return valid
		}

		// Логирование прокси-запроса
		if https_check {
			tools.Logs_file(0, "P-HTTPS", "🔍 IP клиента: "+r.RemoteAddr+" Обработка запроса: https://"+r.Host+r.URL.Path+" → "+proxyConfig.LocalAddress+":"+proxyConfig.LocalPort, "logs_https.log", false)
		} else {
			tools.Logs_file(0, "P-HTTP", "🔍 IP клиента: "+r.RemoteAddr+" Обработка запроса: http://"+r.Host+r.URL.Path+" → "+proxyConfig.LocalAddress+":"+proxyConfig.LocalPort, "logs_http.log", false)
		}

		// Определяем протокол для локального соединения
		protocol := "http"
		if proxyConfig.ServiceHTTPSuse {
			protocol = "https"
		}

		// Проксирование на локальный адрес
		proxyURL := protocol + "://" + proxyConfig.LocalAddress + ":" + proxyConfig.LocalPort + r.URL.RequestURI()
		proxyReq, err := http.NewRequest(r.Method, proxyURL, r.Body)
		if err != nil {
			http.Error(w, "Ошибка создания прокси-запроса", http.StatusInternalServerError)
			return valid
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
		if proxyConfig.ServiceHTTPSuse {
			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // Простая настройка для внутренних соединений
				},
			}
		}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, "Ошибка прокси-запроса", http.StatusBadGateway)
			tools.Logs_file(1, "PROXY", "Ошибка прокси-запроса: "+err.Error(), "logs_proxy.log", false)
			return valid
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
	}

	return valid
}
