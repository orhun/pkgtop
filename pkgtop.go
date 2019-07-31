package main

import (
	"fmt"
	"log"
	"strconv"
	"os/exec"
	str "strings"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var termGrid, dfGrid, pkgGrid *ui.Grid
var pkgText, sysInfoText *widgets.Paragraph
var lists []*widgets.List
var sysInfoCmd = "printf \"Hostname: $(uname -n)\n" + 
		"Kernel: $(uname -s)\n" + 
		"Kernel Release: $(uname -r)\n" + 
		"Kernel Version: $(uname -v)\n" + 
		"Processor Type: $(uname -p)\n" + 
		"Hardware: $(uname --m)\n" + 
		"Hardware Platform: $(uname -i)\n" + 
		"OS: $(uname -o)\n\""
var dfCmd = "df -h | awk '{$1=$1};1 {if(NR>1)print}'"
var dfCount, dfIndex = 4, 0

/*!
 * Parse the 'df' command output as Gauge and GridItem.
 *
 * \param diskUsage (array of 'df' lines)
 * \param s (starting index)
 * \param n (n * entry)
 * \return gauges, entries
 */
func getDfEntries(diskUsage []string, s int, n int) ([]*widgets.Gauge, 
		[]interface {}) {
	/* Use the length of 'df' array if "n"
	 * (entry count to show) is greater. 
	 */
	if len(diskUsage) < n {
		n = len(diskUsage)
	}
	entries := make([]interface{}, n)
	var gauges []*widgets.Gauge
	for i := s; i < s + n; i++ {
		/* Pass the insufficient lines */
		if len(diskUsage[i]) < 5 {
			continue
		}
		/* Create gauge widget from the splitted  
		 * line and add it to the entries slice.
		 */
		dfVal := str.Split(diskUsage[i], " ")
		dfGau := widgets.NewGauge()
		dfGau.Title = fmt.Sprintf("%s ~ (%s/%s) [%s]", 
			dfVal[0], dfVal[2], dfVal[1], dfVal[len(dfVal)-1])
		percent, err := strconv.Atoi(
			str.Replace(dfVal[4], "%", "", 1))
		if err != nil {
			return gauges, nil
		}
		dfGau.Percent = percent
		gauges = append(gauges, dfGau)
		entries[i - s] = ui.NewRow(
			1.0/float64(n),
			ui.NewCol(1.0, dfGau),
		)
	}
	return gauges, entries
}

/*!
 * Execute the 'df' command and show parsed output values with widgets.
 *
 * \param dfIndex (starting index of entries to render)
 * \return dfIndex
 */
func showDfInfo(dfIndex int) int {
	/* Prevent underflow and return the first index. */
	if dfIndex < 0 {
		return 0
	}
	/* Execute the 'df' command and split the output by newline. */
	dfOutput := str.Split(execCmd("sh", "-c", dfCmd), "\n")
	/* Remove the last line on invalid length like '\n' */
	if len(dfOutput) > 0 && len(dfOutput[len(dfOutput)-1]) < 5 {
		dfOutput = dfOutput[:len(dfOutput)-1]
	}
	/* Return the maximum index on overflow */
	if len(dfOutput) - dfIndex < dfCount && len(dfOutput) > dfCount {
		return len(dfOutput) - dfCount
	/* Use the first index on invalid entry count. */
	}else if len(dfOutput) <= dfCount {
		dfIndex = 0
	}
	/* Create and render the widgets. */
	gauges, dfEntries := getDfEntries(
		dfOutput, 
		dfIndex, 
		dfCount)
	dfGrid.Set(dfEntries...)
	ui.Render(dfGrid)
	for _, g := range gauges {
		ui.Render(g)
	}
	return dfIndex
}

func getPkgListEntries(pkgs []string, titles []string) ([]*widgets.List, 
		[]interface {}) {
	var pkgls []*widgets.List
	entries := make([]interface{}, len(titles))
	for i := 0; i < len(titles); i++ {
		var rows []string
		for _, pkg := range pkgs {
			rows = append(rows, str.Split(pkg, "~")[i])
		}
		pkgl := widgets.NewList()
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

/*!
 * Execute a operating system command and capture the output.
 *
 * \param name
 * \param arg
 * \return output
 */
func execCmd(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Execution of '%s' failed with %s\n", name, err)
	}
	return string(out)
}

/*!
 * Initialize the termui and render widgets.
 *
 * \return 0 on success
 */
func initUi() int {
	/* Initialize the termui library */
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	/* Close the UI on function exit */
	defer ui.Close()
	/* Initialize the widgets */
	termGrid, dfGrid, pkgGrid = 
		ui.NewGrid(), 
		ui.NewGrid(), 
		ui.NewGrid()
	pkgText, sysInfoText = 
		widgets.NewParagraph(), 
		widgets.NewParagraph()
	
		
	pkgs := []string {
		"apache~2.4.39-1~6.25MiB~'Fri 11 Jan 2019 03:34:39'",
		"autoconf~2.69-5~2.06MiB~'Fri 11 Jan 2019 03:34:39'",
	}
	titles := []string{"1", "2", "3", "4",}
	lists, pkgEntries := getPkgListEntries(pkgs, titles)
	pkgGrid.Set(ui.NewRow(1.0, pkgEntries...),)

	/* Show the disk usage information */
	dfIndex = showDfInfo(dfIndex)
	/* Show the OS information */
	sysInfoText.Text = execCmd("sh", "-c", sysInfoCmd)
	/* Configure the main grid layout */
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
			ui.NewCol(1.0, pkgText),
		),
	)
	/* Render Grid widgets */
	ui.Render(pkgGrid)
	ui.Render(termGrid)
	
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>", "<C-d>":
				return 0
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				termGrid.SetRect(0, 0, 
					payload.Width, payload.Height)		
				ui.Clear()
				ui.Render(termGrid)
				dfIndex = showDfInfo(dfIndex)
			case "j", "<Down>":
				for _, l := range lists {
					l.ScrollDown()
				}
			case "k", "<Up>":
				for _, l := range lists {
					l.ScrollUp()
				}
			case "d":
				dfIndex = showDfInfo(dfIndex + 1)
			case "f":
				dfIndex = showDfInfo(dfIndex - 1)
			}
		}
		for _, l := range lists {
			ui.Render(l)
		}
	}
	return 0
}


/*!
 * Entry-point
 */
func main() {
	initUi()
}