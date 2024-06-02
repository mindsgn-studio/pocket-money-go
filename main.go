package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := &App{}

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "Pocket Money",
		Width:             800,
		Height:            600,
		DisableResize:     false,
		Fullscreen:        false,
		WindowStartState:  options.Maximised,
		Frameless:         false,
		MinWidth:          400,
		MinHeight:         400,
		MaxWidth:          1280,
		MaxHeight:         1024,
		StartHidden:       false,
		HideWindowOnClose: true,
		BackgroundColour:  &options.RGBA{R: 237, G: 234, B: 228, A: 1},
		AlwaysOnTop:       false,
		AssetServer: &assetserver.Options{
			Assets: assets,
			// Handler:    assetsHandler,
			// Middleware: assetsMidldeware,
		},
		// Menu:                             app.applicationMenu(),
		Logger:             nil,
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.ERROR,
		OnStartup:          app.StartUp,
		OnDomReady:         app.OnDomReady,
		// OnShutdown:                       app.shutdown,
		// OnBeforeClose:                    app.beforeClose,
		CSSDragProperty:                  "--wails-draggable",
		CSSDragValue:                     "drag",
		EnableDefaultContextMenu:         false,
		EnableFraudulentWebsiteDetection: false,
		Bind: []interface{}{
			app,
		},
		// EnumBind: []interface{}{
		//	AllWeekdays,
		//},
		ErrorFormatter: func(err error) any { return err.Error() },
		// SingleInstanceLock: &options.SingleInstanceLock{
		//	UniqueId:               "c9c8fd93-6758-4144-87d1-34bdb0a8bd60",
		//	OnSecondInstanceLaunch: app.onSecondInstanceLaunch,
		// },
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			BackdropType:                      windows.Mica,
			DisablePinchZoom:                  true,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
			// ZoomFactor:           float64,
			// IsZoomControlEnabled: bool,
			// Messages:             *windows.Messages,
			// OnSuspend:            func(),
			// OnResume:             func(),
			// WebviewGpuDisabled:   false,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       false,
				// OnFileOpen:                 app.onFileOpen,
				// OnUrlOpen:                  app.onUrlOpen,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "Pocket Money",
				Message: "Made by Mindsgn Studio",
				// Icon:    icon,
			},
		},
		Linux: &linux.Options{
			// Icon:                icon,
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways,
			ProgramName:         "wails",
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
