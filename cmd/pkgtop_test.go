package main

import (
	"github.com/gizak/termui/v3/widgets"
	"strings"
	"testing"
)

func TestGetDfEntries(t *testing.T) {
	dfOutput := []string{
		"dev dev 0 1.9G 0% /dev",
		"run run 996K 1.9G 1% /run",
		"/dev/sda1 /dev/sda1 51G 18G 75% /",
	}
	dfIndex := 0
	dfCount := 3
	gauges, dfEntries := getDfEntries(
		dfOutput,
		dfIndex,
		dfCount)
	if dfEntries == nil || len(dfEntries) != dfCount ||
		len(gauges) != dfCount {
		t.Errorf("Error occurred while parsing the 'df' values. "+
			"Expected length %d, got %d-%d",
			dfCount, len(gauges), len(dfEntries))
	}
}

func TestGetPkgListEntries(t *testing.T) {
	pkgs := []string{
		"val0;10;x;test1",
		"val1;20;y;test2",
		"test;echo 'y'",
		"[1]|[2]",
	}
	titles := strings.Split(pkgs[len(pkgs)-1], "|")
	lists, entries, optCmds := getPkgListEntries(pkgs)
	if len(lists) != len(titles) ||
		len(entries) != len(titles) {
		t.Errorf("Error occurred while parsing the 'pkg' values. "+
			"Expected length %d, got %d-%d",
			len(titles), len(lists), len(entries))
	} else if optCmds[1] != "echo 'y'" {
		t.Errorf("Error occurred while parsing the 'pkg' values. "+
			"Expected \"echo 'y'\", got \"%s\"", optCmds[1])
	}
}

func TestExecCmd(t *testing.T) {
	echoCmd := execCmd("echo", "test")
	testCmd := execCmd("sh", "-c", "test 10 -eq 10 && echo \"true\"")
	if echoCmd != "test" || testCmd != "true" {
		t.Errorf("Expected 'test-true', got '%s-%s'", echoCmd, testCmd)
	}
}

func TestScrollLists(t *testing.T) {
	l1, l2 := widgets.NewList(), widgets.NewList()
	l1.Rows, l2.Rows = []string{"1", "2", "3"},
		[]string{"x", "y", "z"}
	lists := []*widgets.List{l1, l2}
	if scrollLists(lists, -1, 1, false) != 0 ||
		l1.SelectedRow != 1 {
		t.Errorf("Failed to scroll the List widgets. [row: %d]",
			lists[1].SelectedRow)
	}
}
