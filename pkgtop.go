package main

import (
	"log"
	"strconv"
	"os/exec"
	str "strings"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var i int
var termGrid, dfGrid, pkgGrid *ui.Grid
var pkgText, sysInfoText *widgets.Paragraph
var dfgau *widgets.Gauge
var pkgl *widgets.List
var lists []*widgets.List
var sysInfoCmd = "printf \"Hostname: $(uname -n)\n" + 
		"Kernel: $(uname -s)\n" + 
		"Kernel Release: $(uname -r)\n" + 
		"Kernel Version: $(uname -v)\n" + 
		"Processor Type: $(uname -p)\n" + 
		"Hardware: $(uname --m)\n" + 
		"Hardware Platform: $(uname -i)\n" + 
		"OS: $(uname -o)\n\""

func initWidgets() {
	termGrid, dfGrid, pkgGrid = 
		ui.NewGrid(), 
		ui.NewGrid(), 
		ui.NewGrid()
	pkgText, sysInfoText = 
		widgets.NewParagraph(), 
		widgets.NewParagraph()
}

func getDfEntries(diskUsage []string) []interface {} {
	entries := make([]interface{}, len(diskUsage))
	for i, val := range diskUsage {
		dfval := str.Split(val, "~")
		dfgau = widgets.NewGauge()
		dfgau.Title = dfval[0]
		percent, err := strconv.Atoi(dfval[1])
		if err != nil {
			return nil
		}
		dfgau.Percent = percent
		entries[i] = ui.NewRow(
			1.0/float64(len(diskUsage)),
			ui.NewCol(1.0, dfgau),
		)
	}
	return entries
}

func getPkgListEntries(pkgs []string, titles []string) ([]*widgets.List, []interface {}) {
	var pkgls []*widgets.List
	entries := make([]interface{}, len(titles))
	for i = 0; i < len(titles); i++ {
		var rows []string
		for _, pkg := range pkgs {
			rows = append(rows, str.Split(pkg, "~")[i])
		}
		pkgl = widgets.NewList()
		pkgl.Title = titles[i]
		pkgl.Rows = rows
		pkgl.WrapText = false
		pkgl.Border = false
		pkgl.TextStyle = ui.NewStyle(ui.ColorBlue)
		entries[i] = ui.NewCol(1.0/float64(len(titles)), pkgl)
		pkgls = append(pkgls, pkgl)
	}
	return pkgls, entries
}

func execCmd(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Execution of '%s' failed with %s\n", name, err)
	}
	return string(out)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	initWidgets()
	defer ui.Close()

	diskUsage := []string {
		"dev~0",
		"run~1",
		"/dev/sda1~75",
		"tmpfs~10",
	}
	dfGrid.Set(getDfEntries(diskUsage)...)

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

	lists, entries := getPkgListEntries(pkgs, titles)
	pkgGrid.Set(ui.NewRow(1.0, entries...),)

	sysInfoText.Text = execCmd("sh", "-c", sysInfoCmd)
	
	termWidth, termHeight := ui.TerminalDimensions()
	termGrid.SetRect(0, 0, termWidth, termHeight)
	termGrid.Set(
		ui.NewRow(1.0/4,
			ui.NewCol(1.0/2, dfGrid),
			ui.NewCol(1.0/4, sysInfoText),
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
				for _, l := range lists {
					l.ScrollDown()
				}
			case "k", "<Up>":
				for _, l := range lists {
					l.ScrollUp()
				}
			}
		}
		for _, l := range lists {
			ui.Render(l)
		}
	}
}