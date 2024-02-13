package tui

import "github.com/gdamore/tcell/v2"


type CompText struct {
    *CompDefault

    label string
    style tcell.Style

    spaces bool
}

func NewCompText(l string, s tcell.Style, sps bool) *CompText {
    return &CompText{
        CompDefault: NewCompDefualt(),
        label: l,
        style: s,
        spaces: sps,
    }
}

func (ct *CompText) GetLabel() string {
    return ct.label
}

func (ct *CompText) SetLabel(l string) {
    ct.label = l
    ct.CompDefault.SetRedrawNeeded(true)
}

func (ct *CompText) GetStyle() tcell.Style {
    return ct.style
}

func (ct *CompText) SetStyle(s tcell.Style) {
    ct.style = s
    ct.CompDefault.SetRedrawNeeded(true)
}

func (ct *CompText) GetSpaces() bool {
    return ct.spaces
}

func (ct *CompText) SetSpaces(sps bool) {
    if sps == ct.spaces {
        return
    }

    ct.spaces = sps
    ct.CompDefault.SetRedrawNeeded(true)
}

func (ct *CompText) Draw(s tcell.Screen) bool {
    if !ct.PopRedrawNeeded() || ct.GetRows() <= 0 || ct.GetCols() <= 0 {
        return false
    }

    if !ct.spaces {

        return true
    }


    // Do some shit here idk i feel sick.

    // Draw our text...  
    // Should we break
    
    return true
}




