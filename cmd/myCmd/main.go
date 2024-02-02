package main

import (
	"time"

	"github.com/chathamabate/thingy/tui"
	"github.com/gdamore/tcell/v2"
)


func main() {
    root := tui.NewCompBordered(
        tcell.StyleDefault.Background(tcell.ColorRed),
        tui.NewCompPlainText("Hello World"),
    )

    tui.RunTUI(root, time.Duration(50 * int64(time.Millisecond)))
}
