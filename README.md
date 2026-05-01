# PSX Screener

![Wails](https://img.shields.io/badge/Wails-v2-red?style=flat-square&logo=go)
![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-3-4FC08D?style=flat-square&logo=vuedotjs&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey?style=flat-square)

Anyone who has tried to do serious financial analysis on Pakistani equities knows the problem: there is "Screener" for PSX, no clean API, and the official exchange website is not built for research. You end up copying numbers out of PDFs, switching between scattered websites, and spending more time gathering data than actually analysing it.

I built PSX Screener to fix that for my own workflow. It pulls income statements, balance sheets, and financial ratios for any PSX-listed company into a single clean view — the kind of tool I wished existed before I had to build it myself.

## Features

- Search multiple tickers at once (comma-separated)
- Income statement, balance sheet, and financial ratios side by side
- TTM + last 4 fiscal years for all metrics
- Fast parallel fetching for multi-ticker queries
- Dark-themed UI with keyboard-friendly navigation

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Wails v2](https://wails.io/docs/gettingstarted/installation) — `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

## Getting Started

```bash
# Install frontend dependencies (first run only)
cd frontend && npm install && cd ..

# Run in development mode with hot reload
wails dev
```

## Building

```bash
# Produce a native binary in build/bin/
wails build
```

## Usage

Type one or more PSX ticker symbols in the search box, separated by commas, and press **Enter** or **Fetch**.

```
FFC, ENGRO, LUCK
```

When multiple tickers are entered, a dropdown lets you switch between companies without re-fetching.

## Next Steps

- **Multiple scraper sources** — the scraper is currently tied to stockanalysis.com, which itself has gaps. The plan is to define a common interface and add additional sources (PSX official site, Investing.com) so data can be cross-referenced and coverage improves over time.

- **Local database cache** — right now every lookup hits the network. Adding a local SQLite cache would store results with a timestamp and only re-fetch when stale, which matters a lot when you are flipping through a watchlist of 20 companies.

- **Export to Excel** — the natural end point for most financial workflows is a spreadsheet. An export button that dumps the current company's data into a `.xlsx` file (one sheet per statement) would close that loop without any manual copy-pasting.

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go — scraping via [Colly](https://github.com/gocolly/colly) |
| Frontend | Vue 3 (Composition API), plain CSS |
| Desktop bridge | [Wails v2](https://wails.io) |
| Data source | [stockanalysis.com](https://stockanalysis.com) |
