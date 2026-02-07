package webserver

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	tools "vServer/Backend/tools"
)

var certDir = "WebServer/cert/"
var certMap map[string]*tls.Certificate
var fallbackCert *tls.Certificate
var httpsServer *http.Server
var port_https string = "443"

// GetHTTPSStatus возвращает статус HTTPS сервера
func GetHTTPSStatus() bool {
	return httpsServer != nil
}

// Запуск https сервера
func StartHTTPS() {

	if tools.Port_check("HTTPS", "localhost", port_https) {
		return
	}

	// Отключаем вывод ошибок TLS в консоль
	log.SetOutput(io.Discard)

	// Конфигурация TLS
	tlsConfig := &tls.Config{
		GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
			serverName := chi.ServerName

			if serverName == "" {

			} else if cert, ok := certMap[serverName]; ok {
				// Найден точный сертификат для домена
				return cert, nil

			} else {
				// Пробуем найти сертификат для родительского домена
				parentDomain := getParentDomain(serverName)
				if parentDomain != "" {
					if cert, ok := certMap[parentDomain]; ok {
						tools.Logs_file(1, "HTTPS", "✅ Используем сертификат родительского домена "+parentDomain+" для "+serverName, "logs_https.log", false)
						return cert, nil
					}
				}

				tools.Logs_file(1, "HTTPS", "⚠️ Нет сертификата для: "+serverName, "logs_https.log", false)
			}

			if fallbackCert != nil {
				return fallbackCert, nil
			}

			tools.Logs_file(1, "HTTPS", "❌ Нет fallback-сертификата — соединение будет отклонено", "logs_https.log", true)
			return nil, nil
		},
	}

	// Запуск сервера
	httpsServer = &http.Server{
		Addr:      ":" + port_https,
		TLSConfig: tlsConfig,
		Handler:   nil,
	}

	tools.Logs_file(0, "HTTPS", "✅ HTTPS сервер запущен на порту "+port_https, "logs_https.log", true)

	if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
		// Игнорируем нормальную ошибку при остановке сервера
		if err.Error() != "http: Server closed" {
			tools.Logs_file(1, "HTTPS", "❌ Ошибка запуска сервера: "+err.Error(), "logs_https.log", true)
		}
	}
}

// Извлекает родительский домен из поддомена
func getParentDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) <= 2 {
		return "" // Уже основной домен или некорректный формат
	}
	// Возвращаем домен без первого поддомена
	return strings.Join(parts[1:], ".")
}

// Проверяет, существует ли сертификат для домена

func Cert_start() {
	fallbackCert = loadFallbackCertificate(filepath.Join(certDir, "no_cert"))
	certMap = loadCertificates(certDir)
}

func checkHostCert(r *http.Request) bool {
	host := r.Host

	// Убираем порт если есть
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}

	// Проверяем точное совпадение
	if _, err := os.Stat(certDir + host); err == nil {
		return true
	}

	// Проверяем родительский домен
	parentDomain := getParentDomain(host)
	if parentDomain != "" {
		if _, err := os.Stat(certDir + parentDomain); err == nil {
			return true
		}
	}

	return false
}

func loadCertificates(certDir string) map[string]*tls.Certificate {
	certMap := make(map[string]*tls.Certificate)

	entries, err := os.ReadDir(certDir)
	if err != nil {
		tools.Logs_file(1, "HTTPS", "📁 Ошибка чтения каталога сертификатов: "+err.Error(), "logs_https.log", true)
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "no_cert" {
			continue
		}

		domain := entry.Name()
		certPath := filepath.Join(certDir, domain, "certificate.crt")
		keyPath := filepath.Join(certDir, domain, "private.key")

		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			tools.Logs_file(1, "HTTPS", "⚠️ Ошибка загрузки сертификата для "+domain+": "+err.Error(), "logs_https.log", true)
			continue
		}

		certMap[domain] = &cert
		tools.Logs_file(0, "HTTPS", "✅ Загрузили сертификат для: "+tools.Color(domain, tools.Голубой), "logs_https.log", true)
	}

	return certMap
}

func loadFallbackCertificate(fallbackDir string) *tls.Certificate {
	certPath := filepath.Join(fallbackDir, "certificate.crt")
	keyPath := filepath.Join(fallbackDir, "private.key")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		tools.Logs_file(1, "HTTPS", "⚠️ Не удалось загрузить fallback-сертификат: "+err.Error(), "logs_https.log", true)
		return nil
	}

	tools.Logs_file(0, "HTTPS", "✅ Fallback-сертификат загружен", "logs_https.log", true)
	return &cert
}

func ReloadCertificates() {
	fmt.Println("")
	fmt.Println("🔒 Перезагружаем SSL сертификаты...")
	fmt.Println("")

	// Выгружаем старые сертификаты
	certMap = make(map[string]*tls.Certificate)
	fallbackCert = nil

	fmt.Println("⏹️ Старые сертификаты выгружены")

	// Загружаем сертификаты заново
	Cert_start()

	fmt.Println("✅ SSL сертификаты успешно перезагружены!")
	fmt.Println("")
}

// StopHTTPSServer останавливает HTTPS сервер
func StopHTTPSServer() {
	// Останавливаем HTTPS сервер
	if httpsServer != nil {
		httpsServer.Close()
		httpsServer = nil
		tools.Logs_file(0, "HTTPS", "HTTPS сервер остановлен", "logs_https.log", true)
	}
}
