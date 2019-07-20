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

	pkgText := widgets.NewParagraph()
	pkgText.Text = "pkgtop"
	pkgText.SetRect(0, 0, 10, 5)
	pkgText.BorderStyle.Fg = ui.ColorBlue

	termGrid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	termGrid.SetRect(0, 0, termWidth, termHeight)
	termGrid.Set(
		ui.NewRow(1.0/1,
			ui.NewCol(1.0/2, pkgText),
			ui.NewCol(1.0/2, pkgText),
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