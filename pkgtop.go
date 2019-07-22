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
	width = int(float64(width)/2.5)
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

	pkgList := widgets.NewList()
	pkgList.Rows = []string{
		"apache 2.4.39-1 6.25MiB 'Fri 11 Jan 2019 03:34:39'",
		"autoconf 2.69-5 2.06MiB 'Fri 11 Jan 2019 03:34:39'",
		"automake 1.16.1-1 1598.00KiB 'Fri 11 Jan 2019 03:34:39'",
		"bind-tools 9.14.2-1 5.85MiB 'Fri 11 Jan 2019 03:34:39'",
		"bison 3.3.2-1 2013.00KiB 'Fri 11 Jan 2019 03:34:39'",
		"brook 20190401-1 13.98MiB 'Fri 11 Jan 2019 03:34:39'",
		"chafa 1.0.1-1 327.00KiB 'Fri 11 Jan 2019 03:34:39'",
		"cmatrix 2.0-1 95.00KiB 'Fri 11 Jan 2019 03:34:39'",
		"compton 6.2-2 306.00KiB 'Fri 11 Jan 2019 03:34:39'",
		"docker 1:18.09.6-1 170.98MiB 'Fri 11 Jan 2019 03:34:39'",
	}
	
	termGrid := ui.NewGrid()
	termGrid.SetRect(0, 0, termWidth, termHeight)
	termGrid.Set(
		ui.NewRow(1.0/4,
			ui.NewCol(0.5, dfText),
			ui.NewCol(0.5, pkgText),
		),
		ui.NewRow(3.0/4,
			ui.NewCol(1.0, pkgList),
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