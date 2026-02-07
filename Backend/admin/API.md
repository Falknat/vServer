# 📡 vServer Admin API

API панели управления vServer. Все методы вызываются через Wails IPC биндинги.

> **Доступ из фронтенда:** `window.go.admin.App.MethodName()`

---

## 📋 Содержание

- [Сервисы](#-сервисы)
- [Сайты](#-сайты)
- [Прокси](#-прокси)
- [Сертификаты](#-сертификаты)
- [Конфигурация](#-конфигурация)
- [vAccess](#-vaccess)
- [Типы данных](#-типы-данных)
- [События](#-события)

---

## 🔧 Сервисы

### `GetAllServicesStatus()`
Возвращает статусы всех сервисов.

- **Параметры:** нет
- **Возвращает:** `AllServicesStatus`

```json
{
  "http":  { "name": "HTTP",  "status": true,  "port": "80",  "info": "" },
  "https": { "name": "HTTPS", "status": true,  "port": "443", "info": "" },
  "mysql": { "name": "MySQL", "status": true,  "port": "3306","info": "" },
  "php":   { "name": "PHP",   "status": true,  "port": "8000","info": "" },
  "proxy": { "name": "Proxy", "status": false, "port": "",    "info": "" }
}
```

### `CheckServicesReady()`
Проверяет, готовы ли все основные сервисы (HTTP, HTTPS, MySQL, PHP).

- **Параметры:** нет
- **Возвращает:** `bool`

---

### `StartServer()`
Запускает все сервисы (HTTP, HTTPS, PHP, MySQL, SSL).

- **Параметры:** нет
- **Возвращает:** `string` — `"Server started"`

### `StopServer()`
Останавливает все сервисы.

- **Параметры:** нет
- **Возвращает:** `string` — `"Server stopped"`

### `RestartAllServices()`
Перезапускает все сервисы с перезагрузкой конфига и сертификатов.

- **Параметры:** нет
- **Возвращает:** `string` — `"All services restarted"`

---

### `StartHTTPService()` / `StopHTTPService()`
Управление HTTP сервером (порт 80).

- **Параметры:** нет
- **Возвращает:** `string` — `"HTTP started"` / `"HTTP stopped"`

### `StartHTTPSService()` / `StopHTTPSService()`
Управление HTTPS сервером (порт 443).

- **Параметры:** нет
- **Возвращает:** `string` — `"HTTPS started"` / `"HTTPS stopped"`

### `StartMySQLService()` / `StopMySQLService()`
Управление MySQL сервером.

- **Параметры:** нет
- **Возвращает:** `string` — `"MySQL started"` / `"MySQL stopped"`

### `StartPHPService()` / `StopPHPService()`
Управление PHP FastCGI пулом.

- **Параметры:** нет
- **Возвращает:** `string` — `"PHP started"` / `"PHP stopped"`

### `EnableProxyService()` / `DisableProxyService()`
Включение/отключение прокси-сервиса. Сохраняет изменение в конфиг.

- **Параметры:** нет
- **Возвращает:** `string` — `"Proxy enabled"` / `"Proxy disabled"`

### `EnableACMEService()` / `DisableACMEService()`
Включение/отключение ACME (Let's Encrypt). Сохраняет изменение в конфиг.

- **Параметры:** нет
- **Возвращает:** `string` — `"ACME enabled"` / `"ACME disabled"`

---

## 🌐 Сайты

### `GetSitesList()`
Получает список всех сайтов.

- **Параметры:** нет
- **Возвращает:** `[]SiteInfo`

```json
[
  {
    "name": "My Site",
    "host": "example.com",
    "alias": ["www.example.com"],
    "status": "active",
    "root_file": "index.html",
    "root_file_routing": false,
    "auto_create_ssl": true
  }
]
```

### `CreateNewSite(siteJSON)`
Создаёт новый сайт. Создаёт папку, конфиг и шаблон index.html.

- **Параметры:** `siteJSON: string` — JSON строка с данными `SiteInfo`
- **Возвращает:** `string` — `"Site created successfully"` или `"Error: ..."`

### `DeleteSite(host)`
Удаляет сайт (папку и запись из конфига).

- **Параметры:** `host: string` — домен сайта
- **Возвращает:** `string` — `"Site deleted successfully"` или `"Error: ..."`

### `UpdateSiteCache()`
Обновляет кэш статусов сайтов.

- **Параметры:** нет
- **Возвращает:** `string` — `"Cache updated"`

### `OpenSiteFolder(host)`
Открывает папку сайта в проводнике Windows.

- **Параметры:** `host: string` — домен сайта
- **Возвращает:** `string` — `"Folder opened"` или `"Error: ..."`

---

## 🔀 Прокси

### `GetProxyList()`
Получает список всех прокси-сервисов.

- **Параметры:** нет
- **Возвращает:** `[]ProxyInfo`

```json
[
  {
    "enable": true,
    "external_domain": "app.example.com",
    "local_address": "127.0.0.1",
    "local_port": "3000",
    "service_https_use": false,
    "auto_https": true,
    "auto_create_ssl": true,
    "status": "active"
  }
]
```

---

## 🔒 Сертификаты

### `GetCertInfo(domain)`
Получает информацию о сертификате для домена.

- **Параметры:** `domain: string`
- **Возвращает:** `CertInfo`

```json
{
  "domain": "example.com",
  "issuer": "Let's Encrypt",
  "not_before": "2025-01-01T00:00:00Z",
  "not_after": "2025-03-31T00:00:00Z",
  "days_left": 60,
  "is_expired": false,
  "has_cert": true,
  "dns_names": ["example.com", "www.example.com"]
}
```

### `GetAllCertsInfo()`
Получает информацию о всех сертификатах.

- **Параметры:** нет
- **Возвращает:** `[]CertInfo`

### `ObtainSSLCertificate(domain)`
Получает SSL сертификат через Let's Encrypt для указанного домена.

- **Параметры:** `domain: string`
- **Возвращает:** `string` — `"SSL certificate obtained successfully for ..."` или `"Error: ..."`

### `ObtainAllSSLCertificates()`
Получает сертификаты для всех доменов с флагом `auto_create_ssl: true`.

- **Параметры:** нет
- **Возвращает:** `string` — `"Completed: X success, Y errors"`

### `UploadCertificate(host, certType, certDataBase64)`
Загружает сертификат вручную.

- **Параметры:**
  - `host: string` — домен
  - `certType: string` — тип файла (`"cert"` или `"key"`)
  - `certDataBase64: string` — содержимое файла в Base64
- **Возвращает:** `string` — `"Certificate uploaded successfully"` или `"Error: ..."`

### `DeleteCertificate(domain)`
Удаляет сертификат для домена.

- **Параметры:** `domain: string`
- **Возвращает:** `string` — `"Certificate deleted successfully"` или `"Error: ..."`

### `ReloadSSLCertificates()`
Перезагружает все SSL сертификаты.

- **Параметры:** нет
- **Возвращает:** `string` — `"SSL certificates reloaded"`

---

## ⚙️ Конфигурация

### `GetConfig()`
Возвращает текущую конфигурацию сервера.

- **Параметры:** нет
- **Возвращает:** `object` — полный объект конфигурации

### `SaveConfig(configJSON)`
Сохраняет конфигурацию в файл и перезагружает.

- **Параметры:** `configJSON: string` — JSON строка с конфигурацией
- **Возвращает:** `string` — `"Config saved"` или `"Error: ..."`

### `ReloadConfig()`
Перезагружает конфигурацию из файла.

- **Параметры:** нет
- **Возвращает:** `string` — `"Config reloaded"`

---

## 🛡️ vAccess

### `GetVAccessRules(host, isProxy)`
Получает правила доступа для сайта или прокси.

- **Параметры:**
  - `host: string` — домен
  - `isProxy: bool` — `true` для прокси, `false` для сайта
- **Возвращает:** `VAccessConfig`

```json
{
  "rules": [
    {
      "type": "allow",
      "type_file": [".php", ".html"],
      "path_access": ["/admin"],
      "ip_list": ["192.168.1.0/24"],
      "exceptions_dir": ["/public"],
      "url_error": "/403.html"
    }
  ]
}
```

### `SaveVAccessRules(host, isProxy, configJSON)`
Сохраняет правила доступа.

- **Параметры:**
  - `host: string` — домен
  - `isProxy: bool` — `true` для прокси, `false` для сайта
  - `configJSON: string` — JSON строка с `VAccessConfig`
- **Возвращает:** `string` — `"vAccess saved"` или `"Error: ..."`

---

## 📦 Типы данных

### `ServiceStatus`
```json
{
  "name":   "string",
  "status": "bool",
  "port":   "string",
  "info":   "string"
}
```

### `AllServicesStatus`
```json
{
  "http":  "ServiceStatus",
  "https": "ServiceStatus",
  "mysql": "ServiceStatus",
  "php":   "ServiceStatus",
  "proxy": "ServiceStatus"
}
```

### `SiteInfo`
```json
{
  "name":              "string",
  "host":              "string",
  "alias":             ["string"],
  "status":            "string",
  "root_file":         "string",
  "root_file_routing": "bool",
  "auto_create_ssl":   "bool"
}
```

### `ProxyInfo`
```json
{
  "enable":            "bool",
  "external_domain":   "string",
  "local_address":     "string",
  "local_port":        "string",
  "service_https_use": "bool",
  "auto_https":        "bool",
  "auto_create_ssl":   "bool",
  "status":            "string"
}
```

### `CertInfo`
```json
{
  "domain":     "string",
  "issuer":     "string",
  "not_before": "string",
  "not_after":  "string",
  "days_left":  "int",
  "is_expired": "bool",
  "has_cert":   "bool",
  "dns_names":  ["string"]
}
```

### `VAccessRule`
```json
{
  "type":           "string",
  "type_file":      ["string"],
  "path_access":    ["string"],
  "ip_list":        ["string"],
  "exceptions_dir": ["string"],
  "url_error":      "string"
}
```

### `VAccessConfig`
```json
{
  "rules": ["VAccessRule"]
}
```

---

## 📡 События (Wails Events)

### `service:changed`
Эмитится каждые 500мс с текущими статусами сервисов.

- **Данные:** `AllServicesStatus`

### `server:already_running`
Эмитится при запуске, если vServer уже запущен в другом экземпляре.

- **Данные:** `bool` — `true`
