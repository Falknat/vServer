package acme

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"vServer/Backend/config"
	tools "vServer/Backend/tools"

	"golang.org/x/crypto/acme"
)

var (
	// DefaultManager глобальный менеджер ACME
	DefaultManager *Manager
	
	// Let's Encrypt URLs
	LetsEncryptProduction = "https://acme-v02.api.letsencrypt.org/directory"
	LetsEncryptStaging    = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

// Init инициализирует ACME менеджер
func Init(production bool) error {
	certDir := "WebServer/cert"
	acmeDir := filepath.Join(certDir, ".acme")
	
	// Создаём директорию для ACME данных
	if err := os.MkdirAll(acmeDir, 0700); err != nil {
		return fmt.Errorf("не удалось создать директорию ACME: %w", err)
	}
	
	DefaultManager = &Manager{
		challenges:   make(map[string]*ChallengeData),
		certDir:      certDir,
		acmeDir:      acmeDir,
		isProduction: production,
	}
	
	// Загружаем или создаём account key
	if err := DefaultManager.loadOrCreateAccountKey(); err != nil {
		return fmt.Errorf("ошибка account key: %w", err)
	}
	
	mode := "STAGING"
	if production {
		mode = "PRODUCTION"
	}
	tools.Logs_file(0, "ACME", "✅ ACME менеджер инициализирован ("+mode+")", "logs_acme.log", true)
	
	return nil
}

// CollectDomainsForSSL собирает все домены для получения SSL
func CollectDomainsForSSL() []string {
	domains := make(map[string]bool)
	
	// Из Site_www
	for _, site := range config.ConfigData.Site_www {
		if site.AutoCreateSSL && site.Status == "active" {
			if isValidDomain(site.Host) {
				domains[site.Host] = true
			}
			for _, alias := range site.Alias {
				if isValidDomain(alias) && !strings.Contains(alias, "*") {
					domains[alias] = true
				}
			}
		}
	}
	
	// Из Proxy_Service
	for _, proxy := range config.ConfigData.Proxy_Service {
		if proxy.AutoCreateSSL && proxy.Enable {
			if isValidDomain(proxy.ExternalDomain) {
				domains[proxy.ExternalDomain] = true
			}
		}
	}
	
	// Конвертируем map в slice
	result := make([]string, 0, len(domains))
	for domain := range domains {
		result = append(result, domain)
	}
	
	return result
}

// ObtainCertificate получает сертификат для домена (метод для кнопки в админке)
func ObtainCertificate(domain string) ObtainResult {
	if DefaultManager == nil {
		return ObtainResult{
			Success: false,
			Domain:  domain,
			Error:   "ACME менеджер не инициализирован",
		}
	}
	
	return DefaultManager.obtainCertificate(domain)
}

// ObtainAllCertificates получает сертификаты для всех доменов с AutoCreateSSL
func ObtainAllCertificates() []ObtainResult {
	domains := CollectDomainsForSSL()
	results := make([]ObtainResult, 0, len(domains))
	
	for _, domain := range domains {
		// Проверяем нужно ли получать сертификат
		if !needsCertificate(domain) {
			continue
		}
		
		result := ObtainCertificate(domain)
		results = append(results, result)
		
		// Пауза между запросами чтобы не превысить лимиты
		time.Sleep(time.Second * 2)
	}
	
	return results
}

// CheckAndRenewCertificates проверяет и обновляет истекающие сертификаты
func CheckAndRenewCertificates() []ObtainResult {
	domains := CollectDomainsForSSL()
	results := make([]ObtainResult, 0)
	
	for _, domain := range domains {
		daysLeft := getCertDaysLeft(domain)
		
		// Обновляем если до истечения менее 30 дней
		if daysLeft >= 0 && daysLeft < 30 {
			tools.Logs_file(0, "ACME", "🔄 Обновление сертификата для "+domain+" (осталось "+fmt.Sprintf("%d", daysLeft)+" дней)", "logs_acme.log", true)
			result := ObtainCertificate(domain)
			results = append(results, result)
			time.Sleep(time.Second * 2)
		}
	}
	
	return results
}

// StartBackgroundRenewal запускает фоновую проверку сертификатов
func StartBackgroundRenewal(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for range ticker.C {
			tools.Logs_file(0, "ACME", "🔍 Проверка сертификатов...", "logs_acme.log", false)
			results := CheckAndRenewCertificates()
			
			for _, r := range results {
				if r.Success {
					tools.Logs_file(0, "ACME", "✅ Обновлён: "+r.Domain, "logs_acme.log", true)
				} else {
					tools.Logs_file(1, "ACME", "❌ Ошибка обновления "+r.Domain+": "+r.Error, "logs_acme.log", true)
				}
			}
			
			// Очистка старых challenges
			if DefaultManager != nil {
				DefaultManager.cleanupOldChallenges()
			}
		}
	}()
}

// obtainCertificate внутренний метод получения сертификата
func (m *Manager) obtainCertificate(domain string) ObtainResult {
	tools.Logs_file(0, "ACME", "🔐 Получение сертификата для: "+domain, "logs_acme.log", true)
	
	// Определяем ACME сервер
	acmeURL := LetsEncryptStaging
	if m.isProduction {
		acmeURL = LetsEncryptProduction
	}
	
	// Создаём ACME клиент
	client := &acme.Client{
		Key:          m.accountKey,
		DirectoryURL: acmeURL,
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	
	// Регистрируем аккаунт (если ещё не зарегистрирован)
	_, err := client.Register(ctx, &acme.Account{}, acme.AcceptTOS)
	if err != nil && err != acme.ErrAccountAlreadyExists {
		return ObtainResult{
			Success: false,
			Domain:  domain,
			Error:   "Ошибка регистрации: " + err.Error(),
		}
	}
	
	// Создаём заказ на сертификат
	order, err := client.AuthorizeOrder(ctx, acme.DomainIDs(domain))
	if err != nil {
		return ObtainResult{
			Success: false,
			Domain:  domain,
			Error:   "Ошибка создания заказа: " + err.Error(),
		}
	}
	
	// Обрабатываем авторизации
	for _, authURL := range order.AuthzURLs {
		auth, err := client.GetAuthorization(ctx, authURL)
		if err != nil {
			return ObtainResult{
				Success: false,
				Domain:  domain,
				Error:   "Ошибка получения авторизации: " + err.Error(),
			}
		}
		
		if auth.Status == acme.StatusValid {
			continue
		}
		
		// Ищем HTTP-01 challenge
		var challenge *acme.Challenge
		for _, c := range auth.Challenges {
			if c.Type == "http-01" {
				challenge = c
				break
			}
		}
		
		if challenge == nil {
			return ObtainResult{
				Success: false,
				Domain:  domain,
				Error:   "HTTP-01 challenge не найден",
			}
		}
		
		// Получаем KeyAuth для challenge
		keyAuth, err := client.HTTP01ChallengeResponse(challenge.Token)
		if err != nil {
			return ObtainResult{
				Success: false,
				Domain:  domain,
				Error:   "Ошибка генерации ответа: " + err.Error(),
			}
		}
		
		// Добавляем challenge для обработки HTTP запросов
		m.addChallenge(challenge.Token, keyAuth, domain)
		defer m.removeChallenge(challenge.Token)
		
		// Уведомляем ACME сервер что готовы к проверке
		if _, err := client.Accept(ctx, challenge); err != nil {
			return ObtainResult{
				Success: false,
				Domain:  domain,
				Error:   "Ошибка принятия challenge: " + err.Error(),
			}
		}
		
		// Ждём завершения авторизации
		if _, err := client.WaitAuthorization(ctx, authURL); err != nil {
			return ObtainResult{
				Success: false,
				Domain:  domain,
				Error:   "Ошибка авторизации: " + err.Error(),
			}
		}
	}
	
	// Генерируем приватный ключ для сертификата
	certKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return ObtainResult{
			Success: false,
			Domain:  domain,
			Error:   "Ошибка генерации ключа: " + err.Error(),
		}
	}
	
	// Создаём CSR
	csr, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		DNSNames: []string{domain},
	}, certKey)
	if err != nil {
		return ObtainResult{
			Success: false,
			Domain:  domain,
			Error:   "Ошибка создания CSR: " + err.Error(),
		}
	}
	
	// Ждём готовности заказа и финализируем
	order, err = client.WaitOrder(ctx, order.URI)
	if err != nil {
		return ObtainResult{
			Success: false,
			Domain:  domain,
			Error:   "Ошибка ожидания заказа: " + err.Error(),
		}
	}
	
	// Получаем сертификат
	der, _, err := client.CreateOrderCert(ctx, order.FinalizeURL, csr, true)
	if err != nil {
		return ObtainResult{
			Success: false,
			Domain:  domain,
			Error:   "Ошибка получения сертификата: " + err.Error(),
		}
	}
	
	// Сохраняем сертификат и ключ
	if err := m.saveCertificate(domain, der, certKey); err != nil {
		return ObtainResult{
			Success: false,
			Domain:  domain,
			Error:   "Ошибка сохранения: " + err.Error(),
		}
	}
	
	tools.Logs_file(0, "ACME", "✅ Сертификат получен для: "+domain, "logs_acme.log", true)
	
	return ObtainResult{
		Success: true,
		Domain:  domain,
		Message: "Сертификат успешно получен",
	}
}

// saveCertificate сохраняет сертификат и ключ в файлы
func (m *Manager) saveCertificate(domain string, certDER [][]byte, key *ecdsa.PrivateKey) error {
	certDir := filepath.Join(m.certDir, domain)
	
	// Создаём директорию
	if err := os.MkdirAll(certDir, 0700); err != nil {
		return err
	}
	
	// Сохраняем сертификат (certificate.crt)
	certPath := filepath.Join(certDir, "certificate.crt")
	certFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()
	
	for _, der := range certDER {
		pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: der})
	}
	
	// Сохраняем приватный ключ (private.key)
	keyPath := filepath.Join(certDir, "private.key")
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	
	pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	
	return nil
}

// loadOrCreateAccountKey загружает или создаёт ключ аккаунта
func (m *Manager) loadOrCreateAccountKey() error {
	keyPath := filepath.Join(m.acmeDir, "account.key")
	
	// Пробуем загрузить существующий ключ
	if data, err := os.ReadFile(keyPath); err == nil {
		block, _ := pem.Decode(data)
		if block != nil {
			key, err := x509.ParseECPrivateKey(block.Bytes)
			if err == nil {
				m.accountKey = key
				tools.Logs_file(0, "ACME", "🔑 Account key загружен", "logs_acme.log", false)
				return nil
			}
		}
	}
	
	// Создаём новый ключ
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	
	// Сохраняем ключ
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	
	pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	
	m.accountKey = key
	tools.Logs_file(0, "ACME", "🔑 Новый account key создан", "logs_acme.log", true)
	
	return nil
}

// Вспомогательные функции

func isValidDomain(domain string) bool {
	// Исключаем localhost, IP адреса и локальные домены
	if domain == "" || domain == "localhost" {
		return false
	}
	
	// Исключаем IP адреса
	if net.ParseIP(domain) != nil {
		return false
	}
	
	// Исключаем локальные домены
	if strings.HasSuffix(domain, ".local") || strings.HasSuffix(domain, ".localhost") {
		return false
	}
	
	// Исключаем wildcard
	if strings.Contains(domain, "*") {
		return false
	}
	
	// Базовая проверка формата домена
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]*(\.[a-zA-Z0-9][a-zA-Z0-9-]*)+$`)
	return domainRegex.MatchString(domain)
}

func needsCertificate(domain string) bool {
	certPath := filepath.Join("WebServer/cert", domain, "certificate.crt")
	
	// Если сертификата нет - нужен
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return true
	}
	
	// Если сертификат истекает скоро - нужно обновить
	daysLeft := getCertDaysLeft(domain)
	return daysLeft >= 0 && daysLeft < 30
}

func getCertDaysLeft(domain string) int {
	certPath := filepath.Join("WebServer/cert", domain, "certificate.crt")
	
	data, err := os.ReadFile(certPath)
	if err != nil {
		return -1
	}
	
	block, _ := pem.Decode(data)
	if block == nil {
		return -1
	}
	
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return -1
	}
	
	daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)
	return daysLeft
}

func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetCertInfo получает информацию о сертификате для домена
func GetCertInfo(domain string) CertInfo {
	certPath := filepath.Join("WebServer/cert", domain, "certificate.crt")
	
	info := CertInfo{
		Domain:  domain,
		HasCert: false,
	}
	
	// Проверяем существует ли сертификат
	data, err := os.ReadFile(certPath)
	if err != nil {
		return info
	}
	
	block, _ := pem.Decode(data)
	if block == nil {
		return info
	}
	
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return info
	}
	
	info.HasCert = true
	info.Issuer = cert.Issuer.CommonName
	info.NotBefore = cert.NotBefore.Format("2006-01-02 15:04:05")
	info.NotAfter = cert.NotAfter.Format("2006-01-02 15:04:05")
	info.DaysLeft = int(time.Until(cert.NotAfter).Hours() / 24)
	info.IsExpired = time.Now().After(cert.NotAfter)
	info.DNSNames = cert.DNSNames
	
	return info
}

// DeleteCertificate удаляет сертификат для домена
func DeleteCertificate(domain string) error {
	certDir := filepath.Join("WebServer/cert", domain)
	
	// Проверяем существует ли директория
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		return fmt.Errorf("сертификат для %s не найден", domain)
	}
	
	// Удаляем директорию с сертификатами
	err := os.RemoveAll(certDir)
	if err != nil {
		return fmt.Errorf("ошибка удаления сертификата: %w", err)
	}
	
	tools.Logs_file(0, "ACME", "🗑️ Сертификат удалён для: "+domain, "logs_acme.log", true)
	return nil
}

// GetAllCertsInfo получает информацию о всех сертификатах
func GetAllCertsInfo() []CertInfo {
	certs := make([]CertInfo, 0)
	certBaseDir := "WebServer/cert"
	
	entries, err := os.ReadDir(certBaseDir)
	if err != nil {
		return certs
	}
	
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "no_cert" && entry.Name() != ".acme" {
			info := GetCertInfo(entry.Name())
			if info.HasCert {
				certs = append(certs, info)
			}
		}
	}
	
	return certs
}
