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
