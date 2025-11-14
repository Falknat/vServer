package admin

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	webserver "vServer/Backend/WebServer"
	"vServer/Backend/admin/go/proxy"
	"vServer/Backend/admin/go/services"
	"vServer/Backend/admin/go/sites"
	"vServer/Backend/admin/go/vaccess"
	config "vServer/Backend/config"
	tools "vServer/Backend/tools"
)

type App struct {
	ctx context.Context
}

var appContext context.Context

func NewApp() *App {
	return &App{}
}

var isSingleInstance bool = false

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	appContext = ctx

	// Проверяем, не запущен ли уже vServer
	if !tools.CheckSingleInstance() {
		runtime.EventsEmit(ctx, "server:already_running", true)
		// Только мониторинг, не запускаем сервисы
		config.LoadConfig()
		go a.monitorServices()
		isSingleInstance = false
		return
	}

	isSingleInstance = true

	// Инициализируем время запуска
	tools.ServerUptime("start")

	// Загружаем конфигурацию
	config.LoadConfig()
	time.Sleep(50 * time.Millisecond)

	// Запускаем handler
	webserver.StartHandler()
	time.Sleep(50 * time.Millisecond)

	// Загружаем сертификаты
	webserver.Cert_start()
	time.Sleep(50 * time.Millisecond)

	// Запускаем серверы
	go webserver.StartHTTPS()
	time.Sleep(50 * time.Millisecond)

	go webserver.StartHTTP()
	time.Sleep(50 * time.Millisecond)

	// Запускаем PHP
	webserver.PHP_Start()
	time.Sleep(50 * time.Millisecond)

	// Запускаем MySQL асинхронно
	go webserver.StartMySQLServer(false)

	// Запускаем мониторинг статусов
	go a.monitorServices()
}

func (a *App) GetAllServicesStatus() services.AllServicesStatus {
	return services.GetAllServicesStatus()
}

func (a *App) CheckServicesReady() bool {
	status := services.GetAllServicesStatus()
	return status.HTTP.Status && status.HTTPS.Status && status.MySQL.Status && status.PHP.Status
}

func (a *App) Shutdown(ctx context.Context) {
	// Останавливаем все сервисы при закрытии приложения
	if isSingleInstance {
		webserver.StopHTTPServer()
		webserver.StopHTTPSServer()
		webserver.PHP_Stop()
		webserver.StopMySQLServer()

		// Освобождаем мьютекс
		tools.ReleaseMutex()
	}
}

func (a *App) monitorServices() {
	time.Sleep(1 * time.Second) // Ждём секунду перед первой проверкой

	for {
		time.Sleep(500 * time.Millisecond)
		status := services.GetAllServicesStatus()
		runtime.EventsEmit(appContext, "service:changed", status)
	}
}

func (a *App) GetSitesList() []sites.SiteInfo {
	return sites.GetSitesList()
}

func (a *App) GetProxyList() []proxy.ProxyInfo {
	return proxy.GetProxyList()
}

func (a *App) StartServer() string {
	webserver.Cert_start()

	go webserver.StartHTTPS()
	go webserver.StartHTTP()

	webserver.PHP_Start()
	go webserver.StartMySQLServer(false)

	return "Server started"
}

func (a *App) StopServer() string {
	webserver.StopHTTPServer()
	webserver.StopHTTPSServer()
	webserver.PHP_Stop()
	webserver.StopMySQLServer()

	return "Server stopped"
}

func (a *App) ReloadConfig() string {
	config.LoadConfig()
	return "Config reloaded"
}

func (a *App) GetConfig() interface{} {
	return config.ConfigData
}

func (a *App) SaveConfig(configJSON string) string {
	// Форматируем JSON перед сохранением
	var tempConfig interface{}
	err := json.Unmarshal([]byte(configJSON), &tempConfig)
	if err != nil {
		return "Error: Invalid JSON"
	}

	formattedJSON, err := json.MarshalIndent(tempConfig, "", "    ")
	if err != nil {
		return "Error: " + err.Error()
	}

	// Сохранение конфига в файл
	err = os.WriteFile(config.ConfigPath, formattedJSON, 0644)
	if err != nil {
		return "Error: " + err.Error()
	}

	// Перезагружаем конфигурацию
	config.LoadConfig()
	return "Config saved"
}

func (a *App) RestartAllServices() string {
	// Останавливаем все сервисы
	webserver.StopHTTPServer()
	webserver.StopHTTPSServer()
	webserver.PHP_Stop()
	webserver.StopMySQLServer()
	time.Sleep(500 * time.Millisecond)

	// Перезагружаем конфиг
	config.LoadConfig()
	time.Sleep(200 * time.Millisecond)

	// Обновляем кэш статусов сайтов
	webserver.UpdateSiteStatusCache()

	// Перезагружаем сертификаты
	webserver.Cert_start()
	time.Sleep(50 * time.Millisecond)

	// Запускаем серверы заново
	go webserver.StartHTTPS()
	time.Sleep(50 * time.Millisecond)

	go webserver.StartHTTP()
	time.Sleep(50 * time.Millisecond)

	webserver.PHP_Start()
	time.Sleep(200 * time.Millisecond)

	go webserver.StartMySQLServer(false)

	return "All services restarted"
}

// Управление отдельными сервисами
func (a *App) StartHTTPService() string {
	// Обновляем кэш перед запуском
	webserver.UpdateSiteStatusCache()
	go webserver.StartHTTP()
	return "HTTP started"
}

func (a *App) StopHTTPService() string {
	webserver.StopHTTPServer()
	return "HTTP stopped"
}

func (a *App) StartHTTPSService() string {
	// Обновляем кэш перед запуском
	webserver.UpdateSiteStatusCache()
	go webserver.StartHTTPS()
	return "HTTPS started"
}

func (a *App) StopHTTPSService() string {
	webserver.StopHTTPSServer()
	return "HTTPS stopped"
}

func (a *App) StartMySQLService() string {
	go webserver.StartMySQLServer(false)
	return "MySQL started"
}

func (a *App) StopMySQLService() string {
	webserver.StopMySQLServer()
	return "MySQL stopped"
}

func (a *App) StartPHPService() string {
	webserver.PHP_Start()
	return "PHP started"
}

func (a *App) StopPHPService() string {
	webserver.PHP_Stop()
	return "PHP stopped"
}

func (a *App) EnableProxyService() string {
	config.ConfigData.Soft_Settings.Proxy_enabled = true

	// Сохраняем в файл
	configJSON, _ := json.Marshal(config.ConfigData)
	os.WriteFile(config.ConfigPath, configJSON, 0644)

	return "Proxy enabled"
}

func (a *App) DisableProxyService() string {
	config.ConfigData.Soft_Settings.Proxy_enabled = false

	// Сохраняем в файл
	configJSON, _ := json.Marshal(config.ConfigData)
	os.WriteFile(config.ConfigPath, configJSON, 0644)

	return "Proxy disabled"
}

func (a *App) OpenSiteFolder(host string) string {
	folderPath := "WebServer/www/" + host

	// Получаем абсолютный путь
	absPath, err := tools.AbsPath(folderPath)
	if err != nil {
		return "Error: " + err.Error()
	}

	// Открываем папку в проводнике
	cmd := exec.Command("explorer", absPath)
	err = cmd.Start()
	if err != nil {
		return "Error: " + err.Error()
	}

	return "Folder opened"
}

func (a *App) GetVAccessRules(host string, isProxy bool) *vaccess.VAccessConfig {
	config, err := vaccess.GetVAccessConfig(host, isProxy)
	if err != nil {
		return &vaccess.VAccessConfig{Rules: []vaccess.VAccessRule{}}
	}
	return config
}

func (a *App) SaveVAccessRules(host string, isProxy bool, configJSON string) string {
	var config vaccess.VAccessConfig
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return "Error: Invalid JSON"
	}

	err = vaccess.SaveVAccessConfig(host, isProxy, &config)
	if err != nil {
		return "Error: " + err.Error()
	}

	return "vAccess saved"
}

func (a *App) UpdateSiteCache() string {
	webserver.UpdateSiteStatusCache()
	return "Cache updated"
}
