package main
import (
	"testing"
	str "strings"
)

func TestGetDfEntries(t *testing.T) {
	df := []string {
		"x~14",
		"y~40",
		"z~25",
	}
	if getDfEntries(df) == nil || 
		len(getDfEntries(df)) == 0 {
		t.Error("Error occurred while parsing the df values")
	}
}

func TestGetPkgListEntries(t *testing.T) {
	pkgs := []string {
		"val0~10~x~test1",
		"val1~20~y~test2",
	}
	titles := []string{"1", "2",}
	lists, entries := getPkgListEntries(pkgs, titles)
	if(len(lists) != len(titles) || 
		len(entries) != len(titles)) {
		t.Errorf("Error occurred while parsing the pkg values")
		t.Errorf("Expected length %d, got %d-%d", 
			len(titles), len(lists), len(entries))
	}
}

func TestGetOsInfoText(t *testing.T) {
	info := []string{
		"1~x", 
		"2~y", 
	}
	if len(getOsInfoText(info)) == 0 || 
		!str.Contains(getOsInfoText(info), str.Split(info[0], "~")[0]) {
		t.Errorf("Error occurred while parsing the OS info values")
	}
}