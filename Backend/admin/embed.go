package admin

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed html
var AdminHTML embed.FS

// Флаг для переключения между embed и файловой системой
var UseEmbedded bool = true

// Получает файловую систему в зависимости от флага
func GetFileSystem() fs.FS {
	if UseEmbedded {
		return AdminHTML
	}
	return os.DirFS("Backend/admin/")
}
