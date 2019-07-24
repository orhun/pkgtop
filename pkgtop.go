package main

import (
	"log"
	"strings"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var i int
var termGrid, dfGrid, pkgGrid *ui.Grid
var pkgText *widgets.Paragraph
var dfgau *widgets.Gauge
var pkgl *widgets.List
var pkglx []*widgets.List

func initWidgets() {
	termGrid, dfGrid, pkgGrid = 
		ui.NewGrid(),
		ui.NewGrid(),
		ui.NewGrid()
	pkgText = widgets.NewParagraph()
}

func setDiskUsage(diskUsage map[string]int) bool {
	i := 0
	entries := make([]interface{}, len(diskUsage))
	for name, perc := range diskUsage {
		dfgau = widgets.NewGauge()
		dfgau.Title = name
		dfgau.Percent = perc
		entries[i] = ui.NewRow(
			1.0/float64(len(diskUsage)),
			ui.NewCol(1.0, dfgau),
		)
		i++
	}
	dfGrid.Set(entries...)
	return true
}

func setPkgList(pkgs []string, titles []string) bool {
	entries := make([]interface{}, len(titles))
	for i = 0; i < len(titles); i++ {
		var rows []string
		for _, pkg := range pkgs {
			rows = append(rows, strings.Split(pkg, "~")[i])
		}
		pkgl = widgets.NewList()
		pkgl.Title = titles[i]
		pkgl.Rows = rows
		pkgl.WrapText = false
		pkgl.Border = false
		pkgl.TextStyle = ui.NewStyle(ui.ColorYellow)
		entries[i] = ui.NewCol(1.0/float64(len(titles)), pkgl)
		pkglx = append(pkglx, pkgl)
	}
	pkgGrid.Set(ui.NewRow(1.0, entries...),)
	return true
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	initWidgets()
	defer ui.Close()

	diskUsage := map[string]int {
		"dev": 0,
		"run": 1,
		"/dev/sda1": 75,
		"tmpfs": 4,
	}

	setDiskUsage(diskUsage)

	pkgs := []string {
		"apache~2.4.39-1~6.25MiB~'Fri 11 Jan 2019 03:34:39'",
		"autoconf~2.69-5~2.06MiB~'Fri 11 Jan 2019 03:34:39'",
		"automake~1.16.1-1~1598.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"bind-tools~9.14.2-1~5.85MiB~'Fri 11 Jan 2019 03:34:39'",
		"bison~3.3.2-1~2013.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"brook~1-1~13.98MiB~'Fri 11 Jan 2019 03:34:39'",
		"chafa~1.0.1-1~327.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"cmatrix~2.0-1~95.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"compton~6.2-2~306.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"docker~1-1~170.98MiB~'Fri 11 Jan 2019 03:34:39'",
	}
	titles := []string{"1", "2", "3", "4",}
	
	setPkgList(pkgs, titles)
	
	pkgText.Text = "~"

	termWidth, termHeight := ui.TerminalDimensions()
	termGrid.SetRect(0, 0, termWidth, termHeight)
	termGrid.Set(
		ui.NewRow(1.0/4,
			ui.NewCol(1.0/2, dfGrid),
			ui.NewCol(1.0/4, pkgText),
			ui.NewCol(1.0/4, pkgText),
		),
		ui.NewRow(1.0/1.6,
			ui.NewCol(1.0/1, pkgGrid),
		),
		ui.NewRow(1.0/8,
			ui.NewCol(1.0/1, pkgText),
		),
	)
	ui.Render(termGrid)
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>", "<C-d>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				termGrid.SetRect(0, 0, payload.Width, payload.Height)		
				ui.Clear()
				ui.Render(termGrid)
			case "j", "<Down>":
				for _, l := range pkglx {
					l.ScrollDown()
				}
			case "k", "<Up>":
				for _, l := range pkglx {
					l.ScrollUp()
				}
			}
		}
		for _, l := range pkglx {
			ui.Render(l)
		}
	}

}