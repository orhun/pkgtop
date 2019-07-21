package main

import (
	"fmt"
	"log"
	"strings"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var diskUsageText string

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	diskUsage:= map[string]int{
		"dev": 0,
		"run": 1,
		"/dev/sda1": 75,
		"tmpfs": 4,
	}
	termWidth, termHeight := ui.TerminalDimensions()

	for k, v := range diskUsage {
		diskUsageText += fmt.Sprintf("  %s%s[%s %d%%] \n", k, 
			strings.Repeat(" ", 15-len(k)), strings.Repeat("|", ((termWidth/3)*v)/100), v)
	}

	dfText := widgets.NewParagraph()
	dfText.Text = diskUsageText
	//dfText.Border = false

	pkgText := widgets.NewParagraph()
	pkgText.Text = "~"
	//pkgText.Border = false

	termGrid := ui.NewGrid()
	termGrid.SetRect(0, 0, termWidth, termHeight)
	termGrid.Set(
		ui.NewRow(1.0/4,
			ui.NewCol(0.8, dfText),
			ui.NewCol(0.2, pkgText),
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
				diskUsageText = ""
				for k, v := range diskUsage {
					diskUsageText += fmt.Sprintf("  %s%s[%s %d%%] \n", k, 
						strings.Repeat(" ", 15-len(k)), strings.Repeat("|", ((payload.Width/3)*v)/100), v)
				}
				dfText.Text = diskUsageText
				ui.Clear()
				ui.Render(termGrid)
			}
		}
	}

}