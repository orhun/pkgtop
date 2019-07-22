package main
import (
	"testing"
	"math/rand"
)

func TestBtoi(t *testing.T) {
	i := 1
	if btoi(i == 1) != 1 ||  btoi(i != 1) != 0 {
		t.Errorf("Unable to convert boolean to integer value")
	}
}

func TestMaxValMap(t *testing.T) {
	m := map[string]int{
		"1": rand.Intn(99),
		"2": 100,
		"3": rand.Intn(99),
		"4": rand.Intn(99),
	}
	val := maxValMap(m)
	if val != "2" {
		t.Errorf("Expected '2', got '%s'", val)
	}
}