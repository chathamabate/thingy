package tui
import (
	tcell "github.com/gdamore/tcell/v2"
)

// Example Plain Text Component.
type CompPlainText struct {
    CompDefault
    label string
}

func NewCompPlainText(l string) *CompPlainText {
    return &CompPlainText{label: l} 
}

func (cpt *CompPlainText) Draw(s tcell.Screen) bool {
    if !cpt.RedrawNeeded || cpt.Rows <= 0 || cpt.Cols <= 0 {
        return false
    }

    st := tcell.StyleDefault

    for r := 0; r < cpt.Rows; r++ {
        for c := 0; c <  cpt.Cols; c++ {
            s.SetContent(c + cpt.C, r + cpt.R, ' ', nil, st)
        }
    }

    runeNum := 0
    for _, ru := range cpt.label {
        r := runeNum / cpt.Cols
        if r >= cpt.Rows {
            break
        }

        c := runeNum % cpt.Cols
        
        s.SetContent(c + cpt.C, r + cpt.R, ru, nil, st)
        runeNum++
    }

    cpt.RedrawNeeded = false
    return true
}


// Example Bordered Component.

type CompBordered struct {
    CompDefault
    borderStyle tcell.Style
    child Comp
}

func NewCompBordered(bs tcell.Style, c Comp) *CompBordered {
    return &CompBordered{
        borderStyle: bs,
        child: c,
    }
}

func (cb *CompBordered) Resize(r, c int, rows, cols int) { 
    cb.CompDefault.Resize(r, c, rows, cols)
    cb.child.Resize(r+1, c+1, rows-2, cols-2)
}

func (cb *CompBordered) ForwardEvent(ev tcell.Event) error { 
    return cb.child.ForwardEvent(ev)
}

func (cb *CompBordered) Init() error {
    return cb.child.Init()
}

func (cb *CompBordered) Update() error {
    return cb.child.Update()
}

func (cb *CompBordered) Draw(s tcell.Screen) bool {
    if cb.Rows <= 0 || cb.Cols <= 0 {
        return false
    }

    drew := false

    // Render our border.
    if cb.RedrawNeeded {
        for c := 0; c < cb.Cols; c++ {
            s.SetContent(cb.C + c, cb.R, ' ', nil, cb.borderStyle)
            s.SetContent(cb.C + c, cb.R + cb.Rows - 1, ' ', nil, cb.borderStyle)
        }

        for r := 1; r < cb.Rows-1; r++ {
            s.SetContent(cb.C, cb.R + r, ' ', nil, cb.borderStyle)
            s.SetContent(cb.C + cb.Cols - 1, cb.R + r, ' ', nil, cb.borderStyle)
        }

        cb.RedrawNeeded = false
        drew = true 
    }

    if cb.child.Draw(s) {
        drew = true
    }

    return drew
}

func (cb *CompBordered) Cleanup() {
    cb.child.Cleanup()
}
