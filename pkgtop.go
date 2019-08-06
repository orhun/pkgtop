package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"os/exec"
	"strconv"
	str "strings"
)

var termGrid, dfGrid, pkgGrid *ui.Grid                /* Grid widgets for the layout */
var pkgText, sysInfoText *widgets.Paragraph           /* Paragraph widgets for showing text */
var dfCount, dfIndex = 4, 0                           /* Index and count values for the disk usage widgets */
var sysInfoCmd = "printf \"Hostname: $(uname -n)\n" + /* Print the system information with 'uname' */
	"Kernel: $(uname -s)\n" +
	"Kernel Release: $(uname -r)\n" +
	"Kernel Version: $(uname -v)\n" +
	"Processor Type: $(uname -p)\n" +
	"Hardware: $(uname --m)\n" +
	"Hardware Platform: $(uname -i)\n" +
							"OS: $(uname -o)\n\""
var dfCmd = "df -h | awk '{$1=$1};1 {if(NR>1)print}'" /* Print the disk usage with 'df' */
var pkgsCmd = map[string]string{                      /* Commands for listing the installed packages */
	"arch": "pacman -Qi | awk '/^Name/{name=$3} " +
		"/^Version/{ver=$3} " +
		"/^Description/{desc=substr($0,index($0,$3))} " +
		"/^Installed Size/{size=$4$5; " +
		"print name \"~\" ver \"~\" size \"~\" desc}' " +
		"| sort -h -r -t '~' -k3 " +
		"&& echo 'Name|Version|Installed Size|Description'",
}

/*!
 * Parse the 'df' command output as Gauge and GridItem.
 *
 * \param diskUsage (array of 'df' lines)
 * \param s (starting index)
 * \param n (n * entry)
 * \return gauges, entries
 */
func getDfEntries(diskUsage []string, s int, n int) ([]*widgets.Gauge,
	[]interface{}) {
	/* Use the length of 'df' array if "n"
	 * (entry count to show) is greater.
	 */
	if len(diskUsage) < n {
		n = len(diskUsage)
	}
	entries := make([]interface{}, n)
	var gauges []*widgets.Gauge
	for i := s; i < s+n; i++ {
		/* Pass the insufficient lines. */
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
		entries[i-s] = ui.NewRow(
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
	/* Return the maximum index on overflow. */
	if len(dfOutput)-dfIndex < dfCount && len(dfOutput) > dfCount {
		return len(dfOutput) - dfCount
		/* Use the first index on invalid entry count. */
	} else if len(dfOutput) <= dfCount {
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

// TODO: Update the package parser & unit test
func getPkgListEntries(pkgs []string) ([]*widgets.List,
	[]interface{}) {
	var pkgls []*widgets.List
	if len(pkgs) > 0 && len(pkgs[len(pkgs)-1]) < 5 {
		pkgs = pkgs[:len(pkgs)-1]
	}
	titles := str.Split(pkgs[len(pkgs)-1], "|")
	entries := make([]interface{}, len(titles))
	for i := 0; i < len(titles); i++ {
		var rows []string
		for _, pkg := range pkgs {
			if len(str.Split(pkg, "~")) != len(titles) {
				continue
			}
			rows = append(rows, " "+str.Split(pkg, "~")[i])
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

func scrollLists(lists []*widgets.List, amount int, row int) {
	for _, l := range lists {
		if row != -1 {
			l.SelectedRow = row
		}else {
			l.ScrollAmount(amount)
		}
		ui.Render(l)
	}
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
 * return 0 on exit
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

	// TODO: Parse the package list according to the distribution
	// awk -F '=' '/^ID=/ {print tolower($2)}' /etc/*-release
	lists, pkgEntries := getPkgListEntries(
		str.Split(execCmd("sh", "-c", pkgsCmd["arch"]), "\n"))
	pkgGrid.Set(ui.NewRow(1.0, pkgEntries...))
	ui.Render(pkgGrid)

	/* Show the disk usage information */
	dfIndex = showDfInfo(dfIndex)
	/* Show the OS information */
	sysInfoText.Text = execCmd("sh", "-c", sysInfoCmd)
	/* Configure and render the main grid layout */
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
	ui.Render(termGrid)

	// TODO: Improve the UI key events
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
			case "<Enter>", "<Space>":
				// TODO: Show package information
			case "j", "<Down>":
				scrollLists(lists, 1, -1)
			case "<C-j>":
				scrollLists(lists, -1, 
					len(lists[0].Rows) - 1)
			case "k", "<Up>":
				scrollLists(lists, -1, -1)
			case "<C-k>":
				scrollLists(lists, -1, 0)
			case "l", "<Right>":
				dfIndex = showDfInfo(dfIndex + 1)
			case "h", "<Left>":
				dfIndex = showDfInfo(dfIndex - 1)
			}
		}
	}
}

/*!
 * Entry-point
 */
func main() {
	initUi()
}
