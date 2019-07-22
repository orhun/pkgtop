package main

import (
	"log"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	diskUsage := map[string]int{
		"dev": 0,
		"run": 1,
		"/dev/sda1": 75,
		"tmpfs": 4,
	}

	pkgs := []string{
		"apache~2.4.39-1~6.25MiB~'Fri 11 Jan 2019 03:34:39'",
		"autoconf~2.69-5~2.06MiB~'Fri 11 Jan 2019 03:34:39'",
		"automake~1.16.1-1~1598.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"bind-tools~9.14.2-1~5.85MiB~'Fri 11 Jan 2019 03:34:39'",
		"bison~3.3.2-1~2013.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"brook~20190401-1~13.98MiB~'Fri 11 Jan 2019 03:34:39'",
		"chafa~1.0.1-1~327.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"cmatrix~2.0-1~95.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"compton~6.2-2~306.00KiB~'Fri 11 Jan 2019 03:34:39'",
		"docker~1:18.09.6-1~170.98MiB~'Fri 11 Jan 2019 03:34:39'",
	}

	_ = diskUsage
	_ = pkgs

	dfGrid := ui.NewGrid()

	g0 := widgets.NewGauge()
	g0.Title = "Slim Gauge"
	g0.SetRect(20, 20, 30, 30)
	g0.Percent = 75
	g0.BarColor = ui.ColorRed
	g0.BorderStyle.Fg = ui.ColorWhite
	g0.TitleStyle.Fg = ui.ColorCyan

	
	dfGrid.Set(
		ui.NewRow(0.25,
			ui.NewCol(1.0, g0),
		),
		ui.NewRow(0.25,
			ui.NewCol(1.0, g0),
		),
		ui.NewRow(0.25,
			ui.NewCol(1.0, g0),
		),
		ui.NewRow(0.25,
			ui.NewCol(1.0, g0),
		),
	)

	pkgText := widgets.NewParagraph()
	pkgText.Text = "~"
	//pkgText.Border = false

	termGrid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	termGrid.SetRect(0, 0, termWidth, termHeight)
	termGrid.Set(
		ui.NewRow(1.0/8,
			ui.NewCol(1.0/2, pkgText),
			ui.NewCol(1.0/2, pkgText),
		),
		ui.NewRow(1.0/1.5,
			ui.NewCol(1.0/1, pkgText),
		),
		ui.NewRow(1.0/5,
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
			}
		}
	}

}