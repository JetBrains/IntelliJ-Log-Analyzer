package main

import (
	"embed"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var Version = "1.0.5" //Version variable is being overwritten on built step and passed to the program via ldflag

//go:embed frontend/src
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "JetBrains Log analyzer",
		Width:             1280,
		Height:            800,
		MinWidth:          720,
		MinHeight:         570,
		MaxWidth:          5000,
		MaxHeight:         2000,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		Logger:            NewAppLogger(),
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		Mac: &mac.Options{
			TitleBar:             TitleBarHiddenInset(),
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "JetBrains Logs Analyzer",
				Message: "Version: " + Version + "\n",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

}
func TitleBarHiddenInset() *mac.TitleBar {

	return &mac.TitleBar{
		TitlebarAppearsTransparent: false,
		HideTitle:                  false,
		HideTitleBar:               false,
		FullSizeContent:            false,
		UseToolbar:                 false,
		HideToolbarSeparator:       true,
	}

}
