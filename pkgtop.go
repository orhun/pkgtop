package main

import (
	"fmt"
	"log"
	"strings"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var diskUsage map[string]int

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func maxValMap(m map[string]int) string {
	var max int = 0
	var key string = ""
	for k, v := range m {
        if max < v {
			max = v
			key = k
        }
    }
    return key
}

func getDfText(diskUsage map[string]int, width int) string {
	var diskUsageText string
	width /= 3
	for k, v := range diskUsage {
		diskUsageText += fmt.Sprintf(" %s%s[%s %s%d%%] \n", k, 
			strings.Repeat(" ", len(maxValMap(diskUsage)) + 1 - len(k)), 
			strings.Repeat("|", (width*v)/100), 
			strings.Repeat(" ", width-(width*v)/100 + btoi(v < 10)), v)
	}
	return diskUsageText
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	diskUsage = map[string]int{
		"dev": 0,
		"run": 1,
		"/dev/sda1": 75,
		"tmpfs": 4,
	}
	termWidth, termHeight := ui.TerminalDimensions()


	dfText := widgets.NewParagraph()
	dfText.Text = getDfText(diskUsage, termWidth)
	//dfText.Border = false

	pkgText := widgets.NewParagraph()
	pkgText.Text = "~"
	//pkgText.Border = false

	termGrid := ui.NewGrid()
	termGrid.SetRect(0, 0, termWidth, termHeight)
	termGrid.Set(
		ui.NewRow(1.0/4,
			ui.NewCol(0.6, dfText),
			ui.NewCol(0.4, pkgText),
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
				dfText.Text = getDfText(diskUsage, payload.Width)
				ui.Clear()
				ui.Render(termGrid)
			}
		}
	}

}