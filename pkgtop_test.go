package main
import (
	"testing"
	"strings"
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
		t.Errorf("Error occurred while parsing the 'df' values. " + 
			"Expected length %d, got %d-%d", 
			dfCount, len(gauges), len(dfEntries))
	}
}

func TestGetPkgListEntries(t *testing.T) {
	pkgs := []string {
		"val0~10~x~test1",
		"val1~20~y~test2",
		"[1]|[2]",
	}
	titles := strings.Split(pkgs[len(pkgs)-1], "|")
	lists, entries := getPkgListEntries(pkgs)
	if(len(lists) != len(titles) || 
		len(entries) != len(titles)) {
		t.Errorf("Error occurred while parsing the 'pkg' values. " + 
			"Expected length %d, got %d-%d", 
			len(titles), len(lists), len(entries))
	}
}

func TestExecCmd(t *testing.T) {
	printfCmd := execCmd("printf", "test")
	testCmd := execCmd("sh", "-c", "test 10 -eq 10 && printf \"true\"")
	if printfCmd != "test" || testCmd != "true" {
		t.Errorf("Expected 'test-True', got '%s-%s'", printfCmd, testCmd)
	}
}