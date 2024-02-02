package tui

import "github.com/gdamore/tcell"


func DisplayLine(sc tcell.Screen, 
    x, y int, width, height int, 
    msg string, s tcell.Style) {

    row := 0
    col := 0

    for _, r := range msg {
        sc.SetContent(x + col, y + row, r, nil, s)

        col++
        if col == width {
            col = 0
            row++
        }

        if row == height {
            break
        }
    }
}
