package main

import (
	"fmt"
	"log"
	ui "github.com/gizak/termui/v3"
)

func main() {
	fmt.Println("pkgtop")
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
}