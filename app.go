package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/xuri/excelize/v2"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	scraper "psx-screener-gui/scrapers"
)

type App struct {
	ctx context.Context
}

type CompanyExport struct {
	Ticker string                `json:"ticker"`
	Data   scraper.FinancialData `json:"data"`
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
	if err != nil || savePath == "" {
		return err
	}
	return os.WriteFile(savePath, []byte(content), 0644)
}

func (a *App) ExportExcel(companies []CompanyExport, statementHeaders []string, ratioHeaders []string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not supported on this platform")
	}

	f := excelize.NewFile()
	defer f.Close()

	// Styles
	boldStyle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Size: 11}})
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 10},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1C2128"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	sectionStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "F0A500", Size: 10},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"161B22"}, Pattern: 1},
	})
	priceStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "F0A500", Size: 13},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	numFmt := "#,##0.00"
	numStyle, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &numFmt,
		Alignment:    &excelize.Alignment{Horizontal: "right"},
	})
	pctStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt:    10, // 0.00%
		Alignment: &excelize.Alignment{Horizontal: "right"},
	})

	firstSheet := true
	for _, company := range companies {
		sheet := company.Ticker
		if firstSheet {
			f.SetSheetName("Sheet1", sheet)
			firstSheet = false
		} else {
			f.NewSheet(sheet)
		}

		row := 1
		f.SetCellValue(sheet, cell(1, row), company.Data.CompanyName)
		f.SetCellStyle(sheet, cell(1, row), cell(1, row), boldStyle)
		f.SetRowHeight(sheet, row, 20)
		row++

		f.SetCellValue(sheet, cell(1, row), fmt.Sprintf("PKR %.2f  |  Figures in PKR Millions", company.Data.CurrentPrice))
		f.SetCellStyle(sheet, cell(1, row), cell(1, row), priceStyle)
		row += 2

		sections := []struct {
			title   string
			items   []scraper.LineItem
			headers []string
		}{
			{"Income Statement", company.Data.IncomeStatement, statementHeaders},
			{"Balance Sheet", company.Data.BalanceSheet, statementHeaders},
			{"Ratios", company.Data.Ratios, ratioHeaders},
		}

		for _, sec := range sections {
			// Section title
			f.SetCellValue(sheet, cell(1, row), sec.title)
			f.SetCellStyle(sheet, cell(1, row), cell(len(sec.headers)+1, row), sectionStyle)
			row++

			// Period headers
			f.SetCellValue(sheet, cell(1, row), "Metric")
			for j, h := range sec.headers {
				f.SetCellValue(sheet, cell(j+2, row), h)
			}
			f.SetCellStyle(sheet, cell(1, row), cell(len(sec.headers)+1, row), headerStyle)
			row++

			// Data rows
			for _, item := range sec.items {
				f.SetCellValue(sheet, cell(1, row), item.Label)
				for j, h := range sec.headers {
					v, ok := item.Values[h]
					if !ok {
						row++; continue
					}
					col := cell(j+2, row)
					if item.IsPercent {
						f.SetCellValue(sheet, col, float64(v)/100)
						f.SetCellStyle(sheet, col, col, pctStyle)
					} else {
						f.SetCellValue(sheet, col, v)
						f.SetCellStyle(sheet, col, col, numStyle)
					}
				}
				row++
			}
			row++ // blank between sections
		}

		// Column widths
		f.SetColWidth(sheet, "A", "A", 30)
		for i := range statementHeaders {
			col, _ := excelize.ColumnNumberToName(i + 2)
			f.SetColWidth(sheet, col, col, 14)
		}
	}

	tmp, err := os.CreateTemp("", "psx-screener-*.xlsx")
	if err != nil {
		return err
	}
	tmp.Close()
	if err := f.SaveAs(tmp.Name()); err != nil {
		return err
	}
	return openFile(tmp.Name())
}

func cell(col, row int) string {
	name, _ := excelize.ColumnNumberToName(col)
	return fmt.Sprintf("%s%d", name, row)
}
