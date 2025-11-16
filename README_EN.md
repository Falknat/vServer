# 🚀 vServer - Functional Web Server on Go
**🇷🇺 [Русская версия](README.md)**
> Full-featured web server with HTTP/HTTPS, MySQL, PHP support and GUI admin panel

**👨‍💻 Author:** Roman Sumaneev  
**🌐 Website:** [voxsel.ru](https://voxsel.ru)  
**📞 Contact:** [VK](https://vk.com/felias)

## 🎯 Features

<img src="https://vserf.ru/images/11.jpg" alt="Main page" width="600">
<img src="https://vserf.ru/images/12.jpg" alt="Main page" width="600">

### 🌐 Web Server
- ✅ **HTTP/HTTPS** server with SSL certificate support
- ✅ **Proxy server** for request proxying
- ✅ **PHP server** with built-in PHP 8 support
- ✅ **Static content** for hosting websites
- ✅ **vAccess** - access control system for sites and proxies

### 🗄️ Database
- ✅ **MySQL server** with full support

### 🔧 Administration
- ✅ **GUI Admin Panel** - Wails desktop application with modern interface
- ✅ **Service Management** - start/stop HTTP, HTTPS, MySQL, PHP, Proxy
- ✅ **Site and Proxy Editor** - visual configuration management
- ✅ **vAccess Editor** - access rules configuration through interface

## 🏗️ Architecture

```
vServer/
├── 🎯 main.go              # Main server entry point
│
├── 🔧 Backend/             # Core logic
│   │
│   ├── admin/              # | 🎛️ GUI Admin Panel (Wails) |
│   │   ├── go/             # | Go backend for admin panel |
│   │   └── frontend/       # | Modern UI |
│   │
│   ├── config/             # | 🔧 Go configuration files |
│   ├── tools/              # | 🛠️ Utilities and helpers |
│   └── WebServer/          # | 🌐 Web server modules |
│
├── 🌐 WebServer/           # Web content and configuration
│   │
│   ├── cert/               # | 🔐 SSL certificates |
│   ├── soft/               # | 📦 MySQL and PHP |
│   ├── tools/              # | 📊 Logs and tools |
│   └── www/                # | 🌍 Web content |
│
├── 📄 go.mod               # Go modules
├── 🔨 build_admin.ps1      # Build GUI admin panel
└── 🚀 vSerf.exe            # GUI admin panel (after build)
```

## 🚀 Installation and Launch

### 🔨 Building the Main Server
```powershell
./build_admin.ps1
```

The script will automatically:
- Check/create `go.mod`
- Install dependencies (`go mod tidy`)
- Check/install Wails CLI
- Build the application → `vSerf.exe`

### 📦 Component Preparation
1. Extract `WebServer/soft/soft.rar` archive to `WebServer/soft/` folder
2. Run `vServer.exe` - main server
3. Run `vSerf.exe` - GUI admin panel for management

> 🔑 **Important:** Default MySQL password is `root`

### 📦 Ready Project for Users
Required for operation:
- 📄 `vSerf.exe` - GUI admin panel (optional)
- 📁 `WebServer/` - configuration and resources

> 💡 The `Backend/` folder and `go.mod`, `main.go` files are only needed for development

## ⚙️ Configuration

Configuration via `WebServer/config.json`:

```json
{
  "Site_www": [
    { 
      "name": "Local Site", 
      "host": "127.0.0.1", 
      "alias": ["localhost"],
      "status": "active",
      "root_file": "index.html",
      "root_file_routing": true
    }
  ],
  "Proxy_Service": [
    {
      "Enable": true,
      "ExternalDomain": "git.example.com",
      "LocalAddress": "127.0.0.1",
      "LocalPort": "3333",
      "ServiceHTTPSuse": false,
      "AutoHTTPS": true
    }
  ],
  "Soft_Settings": {
    "mysql_port": 3306, "mysql_host": "127.0.0.1",
    "php_port": 8000, "php_host": "localhost",
    "proxy_enabled": true
  }
}
```

**Main Parameters:**
- `Site_www` - website settings
- `Proxy_Service` - proxy service configuration
- `Soft_Settings` - service ports and hosts (MySQL, PHP, proxy_enabled)

### 🌐 Alias with Wildcard Support

Wildcard (`*`) support in aliases for sites:

```json
{
  "alias": [
    "*.test.com",     // All subdomains of test.com
    "*.test.ru",      // All subdomains of test.ru
    "test.com",       // Exact match
    "api.*"           // api with any zone
  ],
  "host": "test.com"
}
```

**Wildcard Examples:**
- `*.example.com` → `api.example.com`, `admin.example.com`, `test.example.com` ✅
- `example.*` → `example.com`, `example.ru`, `example.org` ✅
- `*example.com` → `test-example.com`, `my-example.com` ✅
- `*` → any domain ✅ (use carefully!)
- `example.com` → only `example.com` ✅ (without wildcard)

### 🔄 Proxy Server

The proxy server allows redirecting external requests to local services.

**Proxy_Service Parameters:**
- `Enable` - enable/disable proxy (true/false)
- `ExternalDomain` - external domain for request interception
- `LocalAddress` - local service address
- `LocalPort` - local service port
- `ServiceHTTPSuse` - use HTTPS for connecting to local service (true/false)
- `AutoHTTPS` - automatically redirect HTTP → HTTPS (true/false)

**Multiple Proxy Example:**
```json
"Proxy_Service": [
  {
    "Enable": true,
    "ExternalDomain": "git.example.com",
    "LocalAddress": "127.0.0.1",
    "LocalPort": "3000",
    "ServiceHTTPSuse": false,
    "AutoHTTPS": true
  },
  {
    "Enable": false,
    "ExternalDomain": "api.example.com",
    "LocalAddress": "127.0.0.1",
    "LocalPort": "8080",
    "ServiceHTTPSuse": false,
    "AutoHTTPS": false
  }
]
```

#### 📖 Detailed Parameter Description:

**`ServiceHTTPSuse`** - protocol for connecting to local service:
- `false` - vServer connects to local service via HTTP (default)
- `true` - vServer connects to local service via HTTPS

**`AutoHTTPS`** - automatic HTTPS redirect:
- `true` - all HTTP requests are automatically redirected to HTTPS (recommended)
- `false` - both HTTP and HTTPS requests are allowed

**How it Works:**
```
Client (HTTP/HTTPS) → vServer (AutoHTTPS check) → Local Service (ServiceHTTPSuse)
```

**Applying Changes:**
- Enter `config_reload` command in console to reload configuration
- Changes will apply to new requests without server restart

## 🔒 vAccess - Access Control System

vServer includes a flexible access control system **vAccess** for sites and proxy services.

### 📁 Configuration Locations

**For Sites:**
```
WebServer/www/{host}/vAccess.conf
```

**For Proxy:**
```
WebServer/tools/Proxy_vAccess/{domain}_vAccess.conf
```

### ⚙️ Main Features

- ✅ **IP Filtering** - allow/block by IP addresses
- ✅ **Path Control** - restrict access to specific directories
- ✅ **File Filtering** - block by extensions (*.php, *.exe)
- ✅ **Exceptions** - flexible rules with exceptions_dir
- ✅ **Custom Errors** - redirects or error pages

### 📝 Configuration Example

```conf
# Allow admin panel only from local IPs
type: Allow
path_access: /admin/*, /api/admin/*
ip_list: 127.0.0.1, 192.168.1.100
url_error: 404

# Block dangerous files in uploads
type: Disable
type_file: *.php, *.exe, *.sh
path_access: /uploads/*
url_error: 404
```

### 📚 Documentation

Detailed vAccess documentation:
- **For Sites:** see `WebServer/www/{host}/vAccess.conf` (examples in file)
- **For Proxy:** see `WebServer/tools/Proxy_vAccess/README.md`

## 📝 Logging

All logs are saved in `WebServer/tools/logs/`:

- 🌐 `logs_http.log` - HTTP requests (including proxy P-HTTP)
- 🔒 `logs_https.log` - HTTPS requests (including proxy P-HTTPS)
- 🔄 `logs_proxy.log` - Proxy server errors
- 🗄️ `logs_mysql.log` - MySQL operations
- 🐘 `logs_php.log` - PHP errors
- ⚙️ `logs_config.log` - Configuration
- 🔐 `logs_vaccess.log` - Access control for sites
- 🔐 `logs_vaccess_proxy.log` - Access control for proxy

## 🔐 SSL Certificates

### Certificate Installation

1. Open `WebServer/` directory
2. Create `cert/` folder (if it doesn't exist)
3. Create a folder with your domain name or IP address
4. Place certificate files with **exact** names:
   ```
   certificate.crt
   private.key
   ca_bundle.crt
   ```
5. Certificate will be automatically loaded at server startup

### 📁 Certificate Structure

```
WebServer/
└── cert/
    ├── example.com/          # Main domain
    │   ├── certificate.crt
    │   ├── private.key
    │   └── ca_bundle.crt
    │
    └── sub.example.com/      # Subdomain (optional)
        ├── certificate.crt
        ├── private.key
        └── ca_bundle.crt
```

### 🎯 Working with Subdomains

**Important:** If no separate folder is created in `cert/` for a subdomain, the parent domain's certificate will be used automatically.

**Examples:**
- ✅ Request to `example.com` → uses certificate from `cert/example.com/`
- ✅ Request to `sub.example.com` (folder exists) → uses `cert/sub.example.com/`
- ✅ Request to `sub.example.com` (folder does NOT exist) → uses `cert/example.com/`

**This is convenient for wildcard certificates:** one certificate in the main domain folder is enough for all subdomains! 🌟

