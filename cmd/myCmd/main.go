package main

import (
	"log"
	"time"

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
    defer s.Fini()

    env := tui.NewEnvironment(s, 2, time.Duration(100 * time.Millisecond))

    pg := "Hello World, wooo, wooo, wooo, wooo, wooo, wooo"

    eid, err := env.CreateAndRegister(
        tui.BorderedElementF(
            "My Element",
            tcell.StyleDefault,
            tcell.StyleDefault.Foreground(tcell.ColorLightCyan), 
            tui.TextElementF(
                tcell.StyleDefault,
                pg,
            ),
        ),
    )

    if err != nil {
        log.Fatal(err)
    }

    err = env.MakeRoot(eid)
    if err != nil {
        log.Fatal(err)
    }

    err = env.Run()
    if err != nil {
        log.Fatal(err)
    }
}
