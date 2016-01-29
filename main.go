package main

import (
	"github.com/xackery/discordeq/menu"
)

func main() {
	isFirstRun := true
	for {
		err := menu.ShowMenu()
		if isFirstRun && err == nil { //Keep looping menu until no errors found
			break
		}
	}
}
