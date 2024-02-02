package main

import (
	"log"

	"github.com/chathamabate/thingy/internal/tui"
	"github.com/gdamore/tcell"
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

    s.Fill(' ', tcell.StyleDefault)

    st := tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorRed)
    tui.DisplayLine(s, 0, 0, 10, 2, "Hello World asdfasdfasdfas", st)

    s.Show()

    exit := false
    for !exit {
        rawE := s.PollEvent()
        switch e := rawE.(type) {
        case *tcell.EventKey:
            kc := e.Key()
            if kc == tcell.KeyEnter {
                exit = true
            }
            break
        }
    }

    s.Fini()
}
