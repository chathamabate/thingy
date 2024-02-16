package tui

import "github.com/gdamore/tcell/v2"

type TextElement struct {
    *DefaultElement

    style tcell.Style
    text string
}

func NewTextElement(s tcell.Style, t string) *TextElement {
    return &TextElement{
        DefaultElement: NewDefaultElement(),
        style: s,
        text: t,
    }
}

func TextElementF(s tcell.Style, t string) ElementFactory {
    return func (env *Environment) (ElementID, error) {
        return env.Register(NewTextElement(s, t))
    }
}

func (bt *TextElement) Draw(s tcell.Screen) {
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

func NewBorderedElement(t string, ts tcell.Style, bs tcell.Style) *BorderedElement {
    return &BorderedElement{
        DefaultElement: NewDefaultElement(),
        title: t,
        titleStyle: ts,
        borderStyle: bs,
    }
}

func BorderedElementF(t string, ts tcell.Style, bs tcell.Style, ef ElementFactory) ElementFactory {
    return func (env *Environment) (ElementID, error) {
        cid, err := ef(env) 
        if err != nil {
            return -1, err
        }

        eid, err := env.Register(NewBorderedElement(t, ts, bs))
        if err != nil {
            return -1, err
        }

        // This will never error.
        env.Attach(eid, cid)
        return eid, nil
    }
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

// This omits left and right endpoints/corners.
func (be *BorderedElement) drawTitleLine(s tcell.Screen) {
    if len(be.title) == 0 || be.GetCols() < 7 {
        for c := 1; c < be.GetCols(); c++ {
            s.SetContent(be.GetC() + c, be.GetR(), horiz, nil, be.borderStyle)
        }

        return
    }

    // Otherwise we actually draw the title!
    s.SetContent(be.GetC() + 1, be.GetR(), horiz, nil, be.borderStyle)
    s.SetContent(be.GetC() + 2, be.GetR(), ' ', nil, be.borderStyle)

    linePos := 3
    for _, ru := range be.title {
        s.SetContent(be.GetC() + linePos, be.GetR(), ru, nil, be.titleStyle) 
        
        linePos++
        if linePos == be.GetCols() - 3 {
            break
        }
    }

    s.SetContent(be.GetC() + linePos, be.GetR(), ' ', nil, be.borderStyle)
    linePos++

    for ; linePos < be.GetCols() - 1; linePos++ {
        s.SetContent(be.GetC() + linePos, be.GetR(), horiz, nil, be.borderStyle)
    }
}

// We are just going to redraw the border here...
func (be *BorderedElement) Draw(s tcell.Screen) {
    if be.GetRows() == 0 || be.GetCols() == 0 {
        return
    }

    // The title will be bordered on the left and right by spaces.
    // These spaces will have the border's style tho.
    // If the title is empty, no spaces will be added to the border.
    // 
    // The title will be fit in a space of size:
    // cols - 6.
    // leftcorner + (1 horiz) + space + title + (n horiz) + rightcorner = cols.

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
        be.drawTitleLine(s)
        s.SetContent(be.GetC() + be.GetCols() - 1, be.GetR(), rightend, nil, be.borderStyle)

        return
    }

    // be.GetRows() >= 2 && be.GetCols() >= 2

    // Title Line
    be.drawTitleLine(s)

    // Bottom Horizontal borders.
    for c := 1; c < be.GetCols() - 1; c++ {
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


