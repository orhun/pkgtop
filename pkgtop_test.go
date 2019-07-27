package main
import (
	"testing"
)

func TestGetDfEntries(t *testing.T) {
	df := []string {
		"x~14",
		"y~40",
		"z~25",
	}
	if getDfEntries(df) == nil || 
		len(getDfEntries(df)) == 0 {
		t.Error("Error occurred while parsing the values")
	}
}
