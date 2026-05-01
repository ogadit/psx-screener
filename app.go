package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	scraper "psx-screener-gui/scrapers"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetData(ticker string) (scraper.FinancialData, error) {
	return scraper.GetFinancials(ticker)
}

func (a *App) GetPlatform() string {
	return runtime.GOOS
}

func (a *App) SaveCSV(filename string, content string) error {
	savePath, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: filename,
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "CSV Files (*.csv)", Pattern: "*.csv"},
		},
	})
	if err != nil {
		return err
	}
	if savePath == "" {
		return nil // user cancelled
	}
	return os.WriteFile(savePath, []byte(content), 0644)
}

func (a *App) OpenInExcel(content string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not supported on this platform")
	}
	f, err := os.CreateTemp("", "psx-screener-*.csv")
	if err != nil {
		return err
	}
	f.WriteString(content)
	f.Close()
	return exec.Command("cmd", "/c", "start", "", f.Name()).Run()
}
