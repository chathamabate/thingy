package main

import (
	"time"

	"github.com/chathamabate/thingy/tui"
)


func main() {
    d := tui.CompDefault{
    }

    tui.RunTUI(&d, time.Duration(50 * int64(time.Millisecond)))
}
