package main

import (
	"time"

	"github.com/chathamabate/thingy/tui"
)


func main() {
    cpt := tui.NewCompPlainText("Hello World")

    tui.RunTUI(cpt, time.Duration(50 * int64(time.Millisecond)))
}
