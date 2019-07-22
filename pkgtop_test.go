package main
import (
	"testing"
)

func TestBtoi(t *testing.T){
	i := 1
	if btoi(i == 1) != 1 ||  btoi(i != 1) != 0 {
		t.Errorf("Unable to convert boolean to integer value")
	}
}