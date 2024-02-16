package tui

import "github.com/gdamore/tcell/v2"

type BasicTextElement struct {
    *DefaultElement

    style tcell.Style
    text string
}

func NewTextElement(s tcell.Style, t string) *BasicTextElement {
    return &BasicTextElement{
        DefaultElement: NewDefaultElement(),
        style: s,
        text: t,
    }
}

func TextElement(s tcell.Style, t string) ElementFactory {
    return func (env *Environment) (ElementID, error) {
        return env.Register(NewTextElement(s, t))
    }
}

func (bt *BasicTextElement) Draw(s tcell.Screen) {
    if bt.GetRows() == 0 || bt.GetCols() == 0 {
        return
    }

    r := 0
    c := 0

    for _, ru := range bt.text {
        s.SetContent(bt.GetC() + c, bt.GetR() + r, 
            ru, nil, bt.style)     
    
        c++
        if c == bt.GetCols() {
            c = 0
            r++

            if r == bt.GetRows() {
                break
            }
        }
    }


    for ; r < bt.GetRows(); r++ {
        for ; c < bt.GetCols(); c++ {
            s.SetContent(bt.GetC() + c, bt.GetR() + r, 
                ' ', nil, bt.style)     
        }

        // After completing the last row with text,
        // We contiue on starting at the beginning of
        // remaining rows.
        c = 0
    }
}

// A bordered element has a title and border.
// If the title is an empty string, just a border.
// A bordered element MUST have one and only one child element.
type BorderedElement struct {
    *DefaultElement

    title string

    titleStyle tcell.Style 
    borderStyle tcell.Style
}

func (be *BorderedElement) Resize(ectx *ElementContext, r, c int, rows, cols int) error {
    err := be.DefaultElement.Resize(ectx, r, c, rows, cols)
    if err != nil {
        return err
    }

    cctx, _ := ectx.Child(0) 
    err = cctx.ForwardResize(r+1, c+1, max(rows-2, 0), max(cols-2, 0))
    return err
}

func (be *BorderedElement) HandleEvent(ectx *ElementContext, ev tcell.Event) error {
    cctx, _ := ectx.Child(0)
    return cctx.ForwardEvent(ev)
}

const (
    horiz = '\u2500' 
    vert = '\u2502'

    topleft = '\u250C'
    topright = '\u2510'
    bottomleft = '\u2514'
    bottomright = '\u2518'
    
    leftend = '\u2576'
    rightend = '\u2574'

    bottomend = '\u2575'
    topend = '\u2577'
)

// We are just going to redraw the border here...
//
func (be *BorderedElement) Draw(s tcell.Screen) {
    if be.GetRows() == 0 || be.GetCols() == 0 {
        return
    }

    // Single cell, do nothing.
    if be.GetCols() == 1 && be.GetRows() == 1 {
        return
    }

    // Single column. rows >= 2.
    if be.GetCols() == 1 {
        s.SetContent(be.GetC(), be.GetR(), topend, nil, be.borderStyle)

        for r := 1; r < be.GetRows() - 1; r++ {
            s.SetContent(be.GetC(), be.GetR() + r, vert, nil, be.borderStyle)
        }

        s.SetContent(be.GetC(), be.GetR() + be.GetRows() - 1, bottomend, nil, be.borderStyle)

        return
    }

    // Single row. cols >= 2
    if be.GetRows() == 1 {
        s.SetContent(be.GetC(), be.GetR(), leftend, nil, be.borderStyle)

        for c := 1; c < be.GetCols() - 1; c++ {
            s.SetContent(be.GetC() + c, be.GetR(), horiz, nil, be.borderStyle)
        }

        s.SetContent(be.GetC() + be.GetCols() - 1, be.GetR(), rightend, nil, be.borderStyle)

        return
    }

    // be.GetRows() >= 2 && be.GetCols() >= 2

    // Horizontal borders.
    for c := 1; c < be.GetCols() - 1; c++ {
        s.SetContent(be.GetC() + c, be.GetR(), horiz, nil, be.borderStyle)
        s.SetContent(be.GetC() + c, be.GetR() + be.GetRows() - 1, horiz, nil, be.borderStyle)
    }

    // Vertical borders.
    for r := 1; r < be.GetRows() - 1; r++ {
        s.SetContent(be.GetC(), be.GetR() + r, vert, nil, be.borderStyle)
        s.SetContent(be.GetC() + be.GetCols() - 1, be.GetR() + r, vert, nil, be.borderStyle)
    }

    // Corners.
    s.SetContent(be.GetC(), be.GetR(), topleft, nil, be.borderStyle)
    s.SetContent(be.GetC() + be.GetCols() - 1, be.GetR(), topright, nil, be.borderStyle)
    s.SetContent(be.GetC() + be.GetCols() - 1, be.GetR() + be.GetRows() - 1, bottomright, nil, be.borderStyle)
    s.SetContent(be.GetC(), be.GetR() + be.GetRows() - 1, bottomleft, nil, be.borderStyle)
}


