package main

import (
	"context"
	"embed"
	"errors"
	"log"
	"os"

	"github.com/issueye/macro_operation/configs"
	"github.com/issueye/macro_operation/internal/model"
	"github.com/issueye/macro_operation/internal/repository"
	"github.com/issueye/macro_operation/internal/service"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// App struct
type App struct {
	log         *log.Logger
	macroRepo   repository.MacroRepository
	recordService *service.RecordService
	playService  *service.PlayService
	macroService *service.MacroService
}

// NewApp creates a new App instance
func NewApp() *App {
	// Set up logging
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(logFile, "", log.LstdFlags)

	return &App{
		log: logger,
	}
}

func main() {
	// Create an instance of the app struct
	app := NewApp()

	// Run Wails application
	err := wails.Run(&options.App{
		Title:            "操作宏录制和回放工具",
		Width:            800,
		Height:           600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			if err := app.onStartup(ctx); err != nil {
				log.Fatal(err)
			}
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

// onStartup is called when the app starts.
func (a *App) onStartup(ctx context.Context) error {
	// Load configuration
	a.log.Println("Starting Macro Recorder App")

	configPath := configs.GetConfigPath()
	config, err := configs.Load(configPath)
	if err != nil {
		a.log.Printf("Warning: Failed to load config from %s: %v", configPath, err)
		a.log.Println("Using default configuration")
		// 使用默认配置
		config = &configs.Config{}
		a.setDefaults(config)
	}

	a.log.Printf("Config loaded: %+v", config.App)

	// Ensure macros directory exists
	if err := os.MkdirAll(config.Storage.MacrosDir, 0755); err != nil {
		a.log.Printf("Failed to create macros directory: %v", err)
		return err
	}

	// Initialize services
	a.macroRepo, err = repository.NewFileMacroRepository(config.Storage.MacrosDir)
	if err != nil {
		a.log.Fatalf("Failed to initialize macro repository: %v", err)
		return err
	}

	a.recordService = service.NewRecordService()
	a.playService = service.NewPlayService()
	a.macroService = service.NewMacroService(a.macroRepo)

	return nil
}

// setDefaults 设置默认配置
func (a *App) setDefaults(config *configs.Config) {
	if config.App.Name == "" {
		config.App.Name = "Macro Recorder"
	}
	if config.App.Version == "" {
		config.App.Version = "1.0.0"
	}
	if config.Storage.MacrosDir == "" {
		config.Storage.MacrosDir = "./macros"
	}
}

// StartRecording starts recording keyboard and mouse events
func (a *App) StartRecording() error {
	a.log.Println("Start recording")
	if a.recordService.IsStarted() {
		return errors.New("already recording")
	}
	return a.recordService.Start()
}

// StopRecording stops recording
func (a *App) StopRecording() error {
	a.log.Println("Stop recording")
	if !a.recordService.IsStarted() {
		return errors.New("not recording")
	}
	return a.recordService.Stop()
}

// GetEventCount returns the number of events recorded
func (a *App) GetEventCount() int {
	return a.recordService.GetEventCount()
}

// IsRecording returns whether recording is in progress
func (a *App) IsRecording() bool {
	return a.recordService.IsStarted()
}

// SaveMacro saves the recorded macro
func (a *App) SaveMacro(name string) (string, error) {
	a.log.Printf("Saving macro: %s", name)

	if a.recordService.IsStarted() {
		return "", errors.New("please stop recording first")
	}

	eventCount := a.recordService.GetEventCount()
	if eventCount == 0 {
		return "", errors.New("no events recorded")
	}

	script, err := a.recordService.GenerateScript(name)
	if err != nil {
		a.log.Printf("Failed to generate script: %v", err)
		return "", err
	}

	if err := a.macroService.SaveFromScript(name, script); err != nil {
		a.log.Printf("Failed to save macro: %v", err)
		return "", err
	}

	return script, nil
}

// ListMacros returns all saved macros
func (a *App) ListMacros() ([]*model.Macro, error) {
	a.log.Println("List macros")
	return a.macroService.List()
}

// PlayMacro plays a saved macro
func (a *App) PlayMacro(name string) error {
	a.log.Printf("Play macro: %s", name)

	macro, err := a.macroService.Load(name)
	if err != nil {
		return err
	}

	return a.playService.PlayScript(macro.Script.Code)
}

// DeleteMacro deletes a saved macro
func (a *App) DeleteMacro(name string) error {
	a.log.Printf("Delete macro: %s", name)
	return a.macroService.Delete(name)
}

// GenerateCurrentScript generates the current script from recorded events
func (a *App) GenerateCurrentScript() (string, error) {
	return a.recordService.GenerateCurrentScript()
}

// RecordAndSave records and saves a macro in one operation
func (a *App) RecordAndSave(name string) (string, error) {
	a.log.Printf("Record and save macro: %s", name)

	if err := a.StartRecording(); err != nil {
		return "", err
	}

	// 这里需要注意，Wails的回调是异步的，所以我们不能在这里直接等待用户操作
	// 实际实现中，应该在前端处理录制和保存的流程，或者使用频道等机制
	return "", errors.New("record and save not implemented in this way")
}

//go:embed all:frontend/dist
var assets embed.FS
