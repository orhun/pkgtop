package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"os"
	"os/exec"
	"strconv"
	str "strings"
)

var termGrid, dfGrid, pkgGrid *ui.Grid      /* Grid widgets for the layout */
var pkgText, sysInfoText *widgets.Paragraph /* Paragraph widgets for showing text */
var cmdList *widgets.List                   /* List widget for the executed commands. */
var dfIndex, pkgIndex = 0, 0                /* Index value for the disk usage widgets & package list */
var showInfo = false                        /* Switch to the package information page */
var searchMode = false                      /* Boolean value for enabling/disabling the search mode */
var searchQuery, searchSuffix = "", ""      /* List title suffix & search query value */
var cmdPrefix = " λ ~ "                     /* Prefix for prepending to the commands */
var cmdConfirm = " [y] "                    /* Confirmation string for commands to execute */
var osIdCmd = "awk -F '=' '/^ID=/ " +       /* Print the OS ID information (for distro checking) */
							"{print tolower($2)}' /etc/*-release"
var sysInfoCmd = "printf \"Hostname: $(uname -n)\\n" + /* Print the system information with 'uname' */
	" Kernel: $(uname -s)\\n" +
	" Kernel Release: $(uname -r)\\n" +
	" Kernel Version: $(uname -v)\\n" +
	" Processor Type: $(uname -p)\\n" +
	" Hardware: $(uname --m)\\n" +
	" Hardware Platform: $(uname -i)\\n" +
							" OS: $(uname -o)\\n\""
var dfCmd = "df -h | awk '{$1=$1};1 {if(NR>1)print}'" /* Print the disk usage with 'df' */
var pkgsCmd = map[string]string{                      /* Commands for listing the installed packages */
	"arch": "pacman -Qi | awk '/^Name/{name=$3} " +
		"/^Version/{ver=$3} " +
		"/^Description/{desc=substr($0,index($0,$3))} " +
		"/^Installed Size/{size=$4$5; " +
		"print name \"~\" ver \"~\" size \"~\" desc}' " +
		"| sort -h -r -t '~' -k3 " +
		"&& echo \"pacman -Qi %s | sed -e 's/^/  /'~pacman -Rcns --noconfirm %s\"" +
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
	/* Find the optimal widget count for the Grid. */
	dfCount := (sysInfoText.Block.Inner.Max.Y + 1) / 3
	/* Execute the 'df' command and split the output by newline. */
	dfOutput := str.Split(execCmd("sh", "-c", dfCmd), "\n")
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

/*!
 * Parse the 'packages' command output as List widgets (GridItem) for Grid.
 *
 * \param pkgs (output lines)
 * \return pkgls, entries, optCmds
 */
func getPkgListEntries(pkgs []string) ([]*widgets.List,
	[]interface{}, []string) {
	/* Create a slice of List widgets. */
	var pkgls []*widgets.List
	/* Create the title and option command slices from the last lines. */
	titles, optCmds := str.Split(pkgs[len(pkgs)-1], "|"),
		str.Split(pkgs[len(pkgs)-2], "~")
	/* Loop through the lines for creating GridItems that contain List widget. */
	entries := make([]interface{}, len(titles))
	for i := 0; i < len(titles); i++ {
		/* Parse the line for package details and append to the 'rows'. */
		var rows []string
		for _, pkg := range pkgs {
			/* Pass the lines that have insufficient length. */
			if len(str.Split(pkg, "~")) != len(titles) {
				continue
			}
			rows = append(rows, " "+str.Split(pkg, "~")[i])
		}
		/* Create a List widget and initialize with the parsed values. */
		pkgl := widgets.NewList()
		pkgl.Title = titles[i]
		pkgl.Rows = rows
		pkgl.WrapText = false
		pkgl.Border = false
		pkgl.TextStyle = ui.NewStyle(ui.ColorBlue)
		/* Add List widget to the GridItem slice. */
		entries[i] = ui.NewCol(1.0/float64(len(titles)), pkgl)
		pkgls = append(pkgls, pkgl)
	}
	return pkgls, entries, optCmds
}

/*!
 * Scroll and render a slice of List widgets.
 *
 * \param lists
 * \param amount
 * \param row
 * \param force
 * \return 0 on success
 */
func scrollLists(lists []*widgets.List, amount int, 
	row int, force bool) int {
	for _, l := range lists {
		if row != -1 {
			l.SelectedRow = row
		} else {
			l.ScrollAmount(amount)
		}
		if len(l.Rows) != 0 || force {
			ui.Render(l)
		}
	}
	return 0
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
	return str.TrimSpace(string(out))
}

/*!
 * Initialize, execute, render and handle.
 *
 * \param osId (Operating system identity)
 * \return 0 on exit
 */
func start(osId string) int {
	/* Initialize the termui library. */
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	/* Close the UI on function exit. */
	defer ui.Close()
	/* Initialize the widgets. */

	// TODO: Add text to pkgText widget.
	// TODO: Set color of widgets.

	termGrid, dfGrid, pkgGrid =
		ui.NewGrid(),
		ui.NewGrid(),
		ui.NewGrid()
	pkgText, sysInfoText =
		widgets.NewParagraph(),
		widgets.NewParagraph()
	cmdList = widgets.NewList()
	cmdList.WrapText = false
	cmdList.TextStyle = ui.NewStyle(ui.ColorBlue)
	/* Update the commands list. */
	cmdList.Rows = []string{cmdPrefix + pkgsCmd[osId],
		cmdPrefix + osIdCmd}
	/* Retrieve packages with the OS command. */
	pkgs := str.Split(execCmd("sh", "-c", pkgsCmd[osId]), "\n")
	/* Check the packages count. */
	if len(pkgs) < 2 {
		ui.Close()
		log.Fatalf("Failed to retrieve package list. (OS: '%s')", osId)
	}
	/* Initialize and render the widgets for showing the package list. */
	lists, pkgEntries, optCmds := getPkgListEntries(pkgs)
	pkgGrid.Set(ui.NewRow(1.0, pkgEntries...))
	/* Show the OS information. */
	cmdList.Rows = append([]string{cmdPrefix + sysInfoCmd}, cmdList.Rows...)
	sysInfoText.Text = " " + execCmd("sh", "-c", sysInfoCmd)
	/* Configure and render the main grid layout. */
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
			ui.NewCol(1.0, cmdList),
		),
	)
	ui.Render(pkgGrid, termGrid)
	/* Show the disk usage information. (post-render) */
	dfIndex = showDfInfo(dfIndex)

	// TODO: Add new key events (copy...)

	/* Get events from termui. */
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			/* Allow typing to the search area if the search mode is on. */
			if searchMode && (len(str.ToLower(e.ID)) == 1 || 
				str.ToLower(e.ID) == "<backspace>") {
				/* Delete the last char from query on the backspace key press. */
				if len(searchQuery) != 0 && str.ToLower(e.ID) == "<backspace>" {
					searchQuery = searchQuery[:len(searchQuery)-1]
				/* Append key to the query. */
				} else if str.ToLower(e.ID) != "<backspace>" {
					searchQuery += str.ToLower(e.ID)
				}
				/* Create lists again for searching. */
				searchLists, _, _ := getPkgListEntries(pkgs)
				/* Empty the current list rows. */
				for _, l := range lists {
					l.Rows = nil
				}
				/* Loop through the first list, compare the query and show results. */
				for s, name := range searchLists[0].Rows {
					if str.Contains(name, searchQuery) {
						for i, l := range searchLists {
							lists[i].Rows = append(lists[i].Rows, l.Rows[s])
						}
					}
				}
				/* Update the search area. */
				lists[0].Title = searchSuffix + searchQuery
				/* Scroll and (force) render the lists. */
				scrollLists(lists, -1, 0, true)
				break
			}
			switch str.ToLower(e.ID) {
			/* Exit search mode or quit. */
			case "q", "<c-c>", "<c-d>":
				if !searchMode {
					return 0
				} 
				searchMode = false
			/* Terminal resize. */
			case "<resize>":
				payload := e.Payload.(ui.Resize)
				termGrid.SetRect(0, 0,
					payload.Width, payload.Height)
				ui.Clear()
				ui.Render(termGrid)
				dfIndex = showDfInfo(dfIndex)
				scrollLists(lists, -1, lists[0].SelectedRow, false)
			/* Go back from information page. */
			case "<backspace>":
				showInfo = true
				fallthrough
			/* Show package information. */
			case "i", "<enter>", "<space>":
				if !showInfo && len(lists[0].Rows) != 0 {
					/* Parse the 'package info' command output after execution,
					 * use first list for showing the information.
					 */
					selectedPkg := lists[0].Rows[lists[0].SelectedRow]
					pkgInfoCmd := fmt.Sprintf(optCmds[0], selectedPkg)
					cmdList.Rows = append([]string{cmdPrefix + pkgInfoCmd}, cmdList.Rows...)
					cmdList.ScrollTop()
					/* Prepare the list widget. */
					lists = lists[:1]
					lists[0].Title = ""
					lists[0].WrapText = !showInfo
					lists[0].Rows = []string{"  " + execCmd("sh", "-c", pkgInfoCmd)}
					/* Set the Grid entries. */
					pkgEntries = nil
					pkgEntries = append(pkgEntries, ui.NewCol(1.0, lists[0]))
					pkgGrid.Set(ui.NewRow(1.0, pkgEntries...))
					searchMode = false
				} else {
					/* Parse the packages with previous command output and show. */
					lists[0].Rows = nil
					lists[0].WrapText = false
					lists, pkgEntries, optCmds = getPkgListEntries(pkgs)
					pkgGrid.Set(ui.NewRow(1.0, pkgEntries...))
				}
				/* Set the flags for showing info and searching package. */
				showInfo = !showInfo
				ui.Render(pkgGrid, cmdList)
				scrollLists(lists, pkgIndex, -1, false)
			/* Scroll down. (packages) */
			case "j", "<down>":
				scrollLists(lists, 1, -1, false)
			/* Scroll to bottom. (packages) */
			case "<c-j>":
				scrollLists(lists, -1,
					len(lists[0].Rows)-1, false)
			/* Scroll up. (packages) */
			case "k", "<up>":
				scrollLists(lists, -1, -1, false)
			/* Scroll to top. (packages) */
			case "<c-k>":
				scrollLists(lists, -1, 0, false)
			/* Scroll down. (disk usage) */
			case "l", "<right>":
				dfIndex = showDfInfo(dfIndex + 1)
			/* Scroll up. (disk usage) */
			case "h", "<left>":
				dfIndex = showDfInfo(dfIndex - 1)
			/* Scroll executed commands list. */
			case "c":
				if cmdList.SelectedRow < len(cmdList.Rows)-1 {
					cmdList.ScrollDown()
				} else {
					cmdList.ScrollTop()
				}
				ui.Render(cmdList)
			case "s":
				/* Allow searching if not showing any package information. */
				if !showInfo {
					/* Set variables for the package searching. */
					searchMode, searchQuery = true, ""
					/* Use the first lists title for the search. */
					if !str.Contains(searchSuffix, "search") {
						searchSuffix = lists[0].Title + " > search: "
					}
					lists[0].Title = searchSuffix
					ui.Render(lists[0])
				}
			/* Remove package. */
			case "r":
				/* Break if no packages found to remove or showing information. */
				if len(lists[0].Rows) == 0 || showInfo {
					break
				}
				/* Add the 'remove' command to command list with confirmation prefix. */
				selectedPkg := lists[0].Rows[lists[0].SelectedRow]
				pkgRemoveCmd := fmt.Sprintf(optCmds[1], selectedPkg)
				cmdList.Rows = append([]string{cmdConfirm + pkgRemoveCmd},
					cmdList.Rows...)
				cmdList.ScrollTop()
				ui.Render(cmdList)
			/* Confirm and execute the command. */
			case "y":
				selectedCmdRow := cmdList.Rows[cmdList.SelectedRow]
				if str.Contains(selectedCmdRow, cmdConfirm) {
					/* Close the UI, execute the command and show output. */
					ui.Close()
					cmd := exec.Command("sh", "-c",
						str.Replace(selectedCmdRow, cmdConfirm, "", -1))
					cmd.Stderr = os.Stderr
					cmd.Stdout = os.Stdout
					err := cmd.Run()
					/* Show the UI again if the execution is successful. */
					if err == nil {
						start(osId)
					}
				}
			}
		}
	}
}

/*!
 * Entry-point
 */
func main() {
	start(execCmd("sh", "-c", osIdCmd))
}
