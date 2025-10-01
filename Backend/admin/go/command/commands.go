package command

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	webserver "vServer/Backend/WebServer"
	admin "vServer/Backend/admin"
	json "vServer/Backend/admin/go/json"
)

func SecurePost(w http.ResponseWriter, r *http.Request) bool {
	// Проверяем, что запрос POST (не GET из браузера)

	if webserver.Secure_post {

		if r.Method != "POST" {
			http.Error(w, "Метод не разрешен. Используйте POST", http.StatusMethodNotAllowed)
			return false
		}

	}

	return true
}

// API обработчик для /api/*
func ApiHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch path {
	case "/api/metrics":
		json.GetAllMetrics(w)
	default:
		http.NotFound(w, r)
	}
}

// JSON обработчик для /json/*
func JsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path

	switch path {
	case "/json/server_status.json":
		w.Write(json.GetServerStatusJSON())
	case "/json/menu.json":
		w.Write(json.GetMenuJSON())
	default:
		http.NotFound(w, r)
	}
}

// Обработчик сервисных команд для /service/*
func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if !SecurePost(w, r) {
		return
	}

	if SiteList(w, r, path) {
		return
	}

	if site_add(w, path) {
		return
	}

	if Service_Run(w, r, path) {
		return
	}
	http.NotFound(w, r)
}

// Обработчик статических файлов из embed
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Убираем ведущий слэш
	path = strings.TrimPrefix(path, "/")

	// Если пустой путь, показываем index.html
	if path == "" {
		path = "html/index.html"
	} else {
		path = "html/" + path
	}

	// Читаем файл из файловой системы (embed или диск)
	fileSystem := admin.GetFileSystem()
	content, err := fs.ReadFile(fileSystem, path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Устанавливаем правильный Content-Type
	ext := filepath.Ext(path)
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	}

	// Предотвращаем кеширование в режиме разработки (когда UseEmbedded = false)
	if !admin.UseEmbedded {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}

	// Отдаем файл
	w.Write(content)

}
