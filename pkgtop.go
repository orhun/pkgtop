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

	appParagraph := widgets.NewParagraph()
	appParagraph.Text = "pkgtop"
	appParagraph.SetRect(0, 0, 10, 5)
	appParagraph.BorderStyle.Fg = ui.ColorBlue

	ui.Render(appParagraph)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>", "<C-d>":
			return
		}
	}
}