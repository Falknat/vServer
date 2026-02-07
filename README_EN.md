# 🚀 vServer - Functional Web Server on Go
**🇷🇺 [Русская версия](README.md)**
> Full-featured web server with HTTP/HTTPS, MySQL, PHP, Let's Encrypt and GUI admin panel

**👨‍💻 Author:** Roman Sumaneev  
**🌐 Website:** [vserf.ru](https://vserf.ru)  
**📞 Contact:** [VK](https://vk.com/felias)

## 🎯 Features

<img src="https://vserf.ru/images/11.jpg" alt="Main page" width="600">
<img src="https://vserf.ru/images/12.jpg" alt="Main page" width="600">

### 🌐 Web Server
- ✅ **HTTP/HTTPS** server with SSL certificate support
- ✅ **Let's Encrypt** — automatic SSL issuance and renewal
- ✅ **Proxy server** — reverse proxy for local services
- ✅ **PHP 8** — built-in support (FastCGI pool)
- ✅ **MySQL** — built-in database server
- ✅ **vAccess** — access control system
- ✅ **Wildcard** — wildcard aliases and certificates support

### 🎛️ GUI Admin Panel
- ✅ **Service management** — start/stop HTTP, HTTPS, MySQL, PHP, Proxy
- ✅ **Site management** — create, edit, delete
- ✅ **Proxy management** — visual reverse proxy configuration
- ✅ **SSL manager** — issue, renew, upload certificates
- ✅ **vAccess editor** — access rules with drag-and-drop
- ✅ **Settings** — MySQL, PHP, proxy, ACME configuration
- ✅ **Dark/Light theme** and **RU/EN** localization

## 🏗️ Architecture

```
vServer/
├── 🎯 main.go              # Entry point (Wails)
│
├── 🔧 Backend/             # Core logic
│   ├── admin/
│   │   ├── go/             # Go admin backend
│   │   └── API.md          # API documentation
│   ├── config/             # Go configuration
│   ├── tools/              # Utilities and helpers
│   └── WebServer/          # Web server modules
│       └── acme/           # Let's Encrypt
│
├── 🖥️ front_vue/           # Vue 3 admin frontend
│   └── src/
│       ├── Core/           # API, stores, i18n, router
│       └── Design/         # Components, views, styles
│
├── 🌐 WebServer/           # Server working files
│   ├── config.json         # Configuration
│   ├── cert/               # SSL certificates
│   ├── soft/               # MySQL and PHP
│   ├── tools/              # Logs, error page, vAccess
│   └── www/                # Web content (sites)
│
├── 🔨 build_admin.ps1      # Build script
└── 🚀 vSerf.exe            # Built application
```

## 🚀 Installation and Launch

### For Users

1. Download the latest [release](https://github.com/AiVoxel/vServer/releases)
2. Extract `WebServer/soft/soft.rar` to `WebServer/soft/`
3. Run `vSerf.exe` — the GUI admin panel will open
4. Server starts automatically, manage everything through the interface

> 🔑 **Default MySQL password:** `root`

### For Developers

```powershell
# Build (checks/installs Go, Node.js, Wails)
./build_admin.ps1
```

The script automatically:
- Checks dependencies (Go, Node.js, npm) — offers to install via `winget`
- Installs Go modules and Wails CLI
- Builds Vue frontend
- Compiles → `vSerf.exe`

## 🔒 vAccess — Access Control

Flexible rules system for sites and proxies. Configured through GUI admin panel (vAccess section).

**Features:**
- IP filtering — allow/block by IP addresses
- Path control — restrict access to directories
- File filtering — block by extensions
- Exceptions — paths excluded from rules
- Custom errors — redirects or error pages

## 🔐 SSL Certificates

### Automatic (Let's Encrypt)
Enable ACME in settings → certificates are issued and renewed automatically.

### Manual
Upload through GUI admin panel (SSL Manager section) or place files:
```
WebServer/cert/{domain}/
├── certificate.crt
├── private.key
└── ca_bundle.crt
```

> 💡 **Wildcard:** one certificate in the main domain folder covers all subdomains.

## 📝 Logging

Logs in `WebServer/tools/logs/`:

| File | Contents |
|------|----------|
| `logs_http.log` | HTTP requests |
| `logs_https.log` | HTTPS requests |
| `logs_proxy.log` | Proxy errors |
| `logs_mysql.log` | MySQL operations |
| `logs_php.log` | PHP errors |
| `logs_config.log` | Configuration |
| `logs_vaccess.log` | Access control |
| `logs_acme.log` | Let's Encrypt |

## 📡 API

Full admin API documentation: [`Backend/admin/API.md`](Backend/admin/API.md)
