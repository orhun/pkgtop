package main
import (
	"testing"
)

func TestGetDfEntries(t *testing.T){
	dfTest := []string {
		"x~14",
		"y~40",
		"z~25",
	}
	if getDfEntries(dfTest) == nil || 
		len(getDfEntries(dfTest)) == 0 {
		t.Error("Error occurred while parsing the values")
	}
}