package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"speaking_hearts/cmd/posttool/rebuild"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend
var assets embed.FS

// App struct defines the application state and backend-to-frontend bindings.
type App struct {
	ctx       context.Context
	rebuilder *rebuild.AudioRebuilder
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{
		rebuilder: rebuild.NewAudioRebuilder(),
	}
}

// startup is called when the app starts. The context is saved so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SelectFolderAndProcess opens a native directory picker and processes the ceremony data.
func (a *App) SelectFolderAndProcess() string {
	// Open a native directory dialog
	selectedPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Speaking Hearts Storage Directory",
	})

	if err != nil {
		return fmt.Sprintf("Error selecting directory: %v", err)
	}

	if selectedPath == "" {
		return "Operation cancelled by user."
	}

	log.Printf("PostTool: Processing directory: %s", selectedPath)

	// Scan the directory for metadata.json files
	files, err := os.ReadDir(selectedPath)
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err)
	}

	count := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" && file.Name() != "metadata.json" { // Match fragment_*.json
			err := a.rebuilder.ProcessMetadata(filepath.Join(selectedPath, file.Name()))
			if err != nil {
				log.Printf("Error processing %s: %v", file.Name(), err)
				continue
			}
			count++
		}
	}

	if count == 0 {
		return "No ceremony fragments found in the selected folder."
	}

	// Export the reconstructed SRT file to the same directory
	outputPath := filepath.Join(selectedPath, "reconstructed_subtitles.srt")
	err = a.rebuilder.ExportSRT(outputPath)
	if err != nil {
		return fmt.Errorf("failed to export SRT: %v", err).Error()
	}

	return fmt.Sprintf("Successfully processed %d fragments! Subtitles exported to: %s", count, outputPath)
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Speaking Hearts - Post-Event Toolset",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Bind: []interface{}{
			app,
		},
		OnStartup: app.startup,
	})

	if err != nil {
		log.Fatal(err)
	}
}
