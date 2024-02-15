package main

import (
	"log"

	"github.com/chathamabate/thingy/tui"
	"github.com/gdamore/tcell/v2"
)


func main() {
    s, err := tcell.NewScreen()
    if err != nil {
        log.Fatal(err)
    }

    err = s.Init()
    if err != nil {
        log.Fatal(err)
    }

}
