package tui

import "github.com/gdamore/tcell"

// NOTE:
// All Panes will either be Horizontal or Vertical flex
// displays.
// 

type Inheritable[T interface{}] interface {
    Inherit() bool 

    // If inherit is 
    GetVal() T
}

type ComponentStyle struct {
    TextStyle tcell.Style 
    BGColor tcell.Color

    Outlined bool
    OutlineStyle tcell.Style

    // How flexible this component is.
    Flex uint8

    // The direction child components will flex.
    FlexDirection bool
}

type Pane interface {
    // Returns nil if this is the root pane.
    GetParent() Pane

    GetFlex()

    Resize(width, height int) error
}

func StartTUI() error {
    s, err := tcell.NewScreen()
    if err != nil {
        return err
    }

    err = s.Init()
    if err != nil {
        return err
    }

    exit := false
    for !exit {
        
    }
}

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
