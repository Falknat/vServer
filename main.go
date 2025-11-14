package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	admin "vServer/Backend/admin/go"
)

//go:embed all:Backend/admin/frontend
var assets embed.FS

func main() {
	// Создаём экземпляр приложения
	app := admin.NewApp()

	// Настройки окна
	err := wails.Run(&options.App{
		Title:     "vServer - Панель управления",
		Width:     1600,
		Height:    900,
		MinWidth:  1400,
		MinHeight: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 10, G: 14, B: 26, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
		Frameless: true,
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			Theme:                windows.Dark,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
