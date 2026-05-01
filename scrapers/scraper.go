package scraper

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/gocolly/colly"
)

var PreferredItems = map[string][]string{
	"incomeStatement": {"Revenue", "Gross Profit", "Gross Margin", "Net Income", "Profit Margin", "EBITDA", "EBITDA Margin"},
	"balanceSheet":    {"Cash & Short-Term Investments", "Total Current Assets", "Total Assets", "Total Current Liabilities", "Total Liabilities", "Total Common Equity", "Net Cash (Debt)"},
	"ratios":          {"Market Capitalization", "Enterprise Value", "PE Ratio", "EV/Sales Ratio", "EV/EBITDA Ratio", "Return on Equity (ROE)"},
}

var PreferredHeaders = map[string][]string{
	"statements": {"TTM", "FY 2025", "FY 2024", "FY 2023", "FY 2022", "FY 2021"},
	"ratios":     {"Current", "FY 2025", "FY 2024", "FY 2023", "FY 2022", "FY 2021"},
}

// LineItem represents a single financial metric across time periods
// e.g. {"label": "Revenue", "isPercent": false, "values": {"TTM": 1234.5, "FY 2024": 1100.0}}
type LineItem struct {
	Label     string             `json:"label"`
	IsPercent bool               `json:"isPercent"`
	Values    map[string]float32 `json:"values"`
}

type FinancialData struct {
	CompanyName     string     `json:"companyName"`
	CurrentPrice    float32    `json:"currentPrice"`
	IncomeStatement []LineItem `json:"incomeStatement"`
	BalanceSheet    []LineItem `json:"balanceSheet"`
	Ratios          []LineItem `json:"ratios"`
}

func GetFinancials(ticker string) (FinancialData, error) {
	companyName, currentPrice, err := getCompany(ticker)
	if err != nil {
		return FinancialData{}, err
	}
	return FinancialData{
		CompanyName:     companyName,
		CurrentPrice:    currentPrice,
		IncomeStatement: getData(ticker, "https://stockanalysis.com/quote/psx/%s/financials/", PreferredItems["incomeStatement"], PreferredHeaders["statements"]),
		BalanceSheet:    getData(ticker, "https://stockanalysis.com/quote/psx/%s/financials/balance-sheet/", PreferredItems["balanceSheet"], PreferredHeaders["statements"]),
		Ratios:          getData(ticker, "https://stockanalysis.com/quote/psx/%s/financials/ratios/", PreferredItems["ratios"], PreferredHeaders["ratios"]),
	}, nil
}

// parseValue handles "12.34%", "1,234,567", "N/A", "-" etc.
// Returns the float32 value and whether it was a percentage
func parseValue(raw string) (float32, bool) {
	s := strings.TrimSpace(raw)
	isPercent := strings.HasSuffix(s, "%")
	s = strings.TrimSuffix(s, "%")
	s = strings.ReplaceAll(s, ",", "")

	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, isPercent // covers "N/A", "-", ""
	}
	return float32(val), isPercent
}

func PrintFinancials(ticker string, fd FinancialData) {
	fmt.Printf("\n=== %s Financial Summary ===\n\n", ticker)
	fmt.Println("--- Income Statement ---")
	PrintLineItems(fd.IncomeStatement, PreferredHeaders["statements"])
	fmt.Println("\n--- Balance Sheet ---")
	PrintLineItems(fd.BalanceSheet, PreferredHeaders["statements"])
	fmt.Println("\n--- Ratios ---")
	PrintLineItems(fd.Ratios, PreferredHeaders["ratios"])
}

func PrintLineItems(items []LineItem, headers []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Header row
	fmt.Fprintf(w, "\t%s\n", strings.Join(headers, "\t"))

	// Separator
	sep := make([]string, len(headers)+1)
	for i := range sep {
		sep[i] = strings.Repeat("-", 12)
	}
	fmt.Fprintln(w, strings.Join(sep, "\t"))

	// Data rows
	for _, item := range items {
		fmt.Fprintf(w, "%s", item.Label)
		for _, h := range headers {
			val, ok := item.Values[h]
			if !ok {
				fmt.Fprintf(w, "\tN/A")
			} else if item.IsPercent {
				fmt.Fprintf(w, "\t%.2f%%", val)
			} else {
				fmt.Fprintf(w, "\t%.2f", val)
			}
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func getCompany(ticker string) (string, float32, error) {
	var companyName string
	var currentPrice float32
	var httpErr error

	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 404 {
			httpErr = fmt.Errorf("ticker %q not found on PSX", ticker)
		} else {
			httpErr = fmt.Errorf("request failed (HTTP %d): %w", r.StatusCode, err)
		}
	})

	c.OnHTML("#main h1", func(h *colly.HTMLElement) {
		companyName = strings.TrimSpace(h.Text)
	})

	c.OnHTML("div.text-4xl.font-bold", func(h *colly.HTMLElement) {
		price, err := strconv.ParseFloat(strings.TrimSpace(h.Text), 64)
		if err == nil {
			currentPrice = float32(price)
		}
	})

	if err := c.Visit(fmt.Sprintf("https://stockanalysis.com/quote/psx/%s/", ticker)); err != nil {
		return "", 0, fmt.Errorf("could not reach stockanalysis.com: %w", err)
	}
	if httpErr != nil {
		return "", 0, httpErr
	}
	if companyName == "" {
		return "", 0, fmt.Errorf("no data found for %q — verify the ticker", ticker)
	}
	return companyName, currentPrice, nil
}

func getData(ticker, url string, preferredLineItems, preferredHeaders []string) []LineItem {
	c := colly.NewCollector()
	headerIndexMap := map[string]int{}
	result := []LineItem{}

	c.OnHTML("table", func(table *colly.HTMLElement) {
		// Map preferred headers to their column indices
		table.ForEach("th", func(i int, th *colly.HTMLElement) {
			text := strings.TrimSpace(th.Text)
			if slices.Contains(preferredHeaders, text) {
				headerIndexMap[text] = i
			}
		})

		table.ForEach("tr", func(_ int, row *colly.HTMLElement) {
			var allCols []string
			row.ForEach("td", func(_ int, td *colly.HTMLElement) {
				allCols = append(allCols, strings.TrimSpace(td.Text))
			})

			if len(allCols) == 0 || !slices.Contains(preferredLineItems, allCols[0]) {
				return
			}

			item := LineItem{
				Label:  allCols[0],
				Values: map[string]float32{},
			}

			for _, h := range preferredHeaders {
				idx, exists := headerIndexMap[h]
				if !exists || idx >= len(allCols) {
					// Leave missing keys absent — PrintLineItems handles the N/A display
					continue
				}
				val, isPercent := parseValue(allCols[idx])
				item.Values[h] = val
				item.IsPercent = isPercent // all values in a row share the same type
			}

			result = append(result, item)
		})
	})

	c.Visit(fmt.Sprintf(url, ticker))
	return result
}
