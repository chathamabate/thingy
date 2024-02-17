package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// -------------------------------------- Text Element --------------------------------------

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

// -------------------------------------- Bordered Element --------------------------------------

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

// -------------------------------------- Divided Element --------------------------------------

type DivisionSpec interface {
    IsFixed() bool
    FixedSize() int

    IsFlexible() bool
    FlexFactor() int
}

type FixedSpec struct {
    fixedSize int
}

func (fs FixedSpec) IsFixed() bool {
    return true
}

func (fs FixedSpec) FixedSize() int {
    return fs.fixedSize
}

func (fs FixedSpec) IsFlexible() bool {
    return false
}

func (fs FixedSpec) FlexFactor() int {
    return 0
}

type FlexSpec struct {
    flexFactor int
}

func (fs FlexSpec) IsFixed() bool {
    return false
}

func (fs FlexSpec) FixedSize() int {
    return 0
}

func (fs FlexSpec) IsFlexible() bool {
    return true
}

func (fs FlexSpec) FlexFactor() int {
    return fs.flexFactor
}

// Helper function for calculating dimmensions of child elements.
func mapToDims(totalDim int, specs []DivisionSpec) ([]int, error) { 
    totalFixedDim := 0
    totalFlexFactor := 0
    for i, spec := range specs {
        if spec.IsFixed() {
            if spec.FixedSize() < 0 {
                return nil, fmt.Errorf("mapToDims: Negative fixed dim: %d", i)
            }
            
            totalFixedDim += spec.FixedSize()
        } else {
            if spec.FlexFactor() < 0 {
                return nil, fmt.Errorf("mapToDims: Negative flex factor: %d", i)
            }

            totalFlexFactor += spec.FlexFactor()
        }
    }

    flexUnit := 0
    if totalFlexFactor > 0 && totalFixedDim < totalDim {
        flexUnit := (totalDim - totalFixedDim) / totalFlexFactor
    }

    dims := make([]int, len(specs)) 

    pos := 0
    for i, spec := range specs {
        spaceLeft := totalDim - pos 

        dims[i] = 0

        // NOTE: if we reach a fixed size division which cannot be drawn,
        // Everything is given size 0! There is simply not enough room!
        if spec.IsFixed() {
            if spec.FixedSize() > spaceLeft {
                return make([]int, len(specs)), nil
            }
            
            dims[i] = spec.FixedSize()
        } else {
            dims[i] = spec.FlexFactor() * flexUnit
        }

        // Advance our position.
        pos += dims[i]
    }

    // NOTE: If there is space left after calculating the lengths
    // Just add it to the first flex division found.
    // If there are no flex divisions, this area will be left blank.
    // Fixed divisions NEVER resize.
    if pos < totalDim {
        spaceLeft := totalDim - pos 
        for i, spec := range specs {
            if spec.IsFlexible() {
                dims[i] += spaceLeft
            }
        }
    }
}

// A divided element holds a variable number of child elements. 
// A divided element CAN hold zero elements.
// A divided element can hold column divisions or row division, but not both.
//
// A division's size can be variable (with a flex coeficient)
// Or a division's size can be fixed (with an exact row or column amount)

type DividedElement struct {
    *DefaultElement
    
    // true = column divisions.
    // false = row divisions.
    columnDivisions bool   

    // If true, there were will line dividers between each division.
    dividers bool

    // NOTE: Since a divided element can hold no elements and/or
    // not be entirely filled by its children, extra space (and dividers) will
    // take this style.
    style tcell.Style
}

// NOTE:
// All children need a "div-spec" attribute which maps to a DivisionSpec.
// This attribute determines how the child will be resized.
//
// A fixed division is either displayed at its specified size, or not displayed at all.
// The moment a fixed division cannot be displayed, drawing stops, no subsequent flex divisions
// will be displayed. flex divisions are only displayed if there is enough room including all fixed 
// divisions.

func (de *DividedElement) Resize(ectx *ElementContext, r, c int, rows, cols int) error {
    de.DefaultElement.Resize(ectx, r, c, rows, cols)

    numChildren := ectx.NumChildren()

    totalFixed := 0
    totalFlex := 0
    
    var totalDim int
    if de.columnDivisions {
        totalDim = cols
    } else {
        totalDim = rows
    }

    // First, let's extract the division specs.
    specs := make([]DivisionSpec, numChildren)   
    for i := 0; i < numChildren; i++ {
        val, err := ectx.GetChildAttr(i, "div-spec")
        if err != nil {
            return err
        }

        ds, ok := val.(DivisionSpec)
        if !ok {
            return fmt.Errorf("Resize: div-spec has incorred type: %d", i)
        }

        specs[i] = ds

        if ds.IsFixed() {
            totalFixed += ds.FixedSize()
        } else {
            totalFlex += ds.FlexFactor()
        }
    }

    // flexUnit = (rows or columns) per unit of flex.
    flexUnit := 0
    if totalFixed < totalDim && totalFlex > 0 {
        flexUnit = (totalDim - totalFixed) / totalFlex
    }

    pos := 0
    index := 0
    for ; index < numChildren; index++ {
        spaceLeft := totalDim - pos

        var dim int   
        if spec.IsFixed() {
            dim = spec.FixedSize()
        } else {
            dim = spec.FlexFactor() * flexUnit
        }

        if dim > spaceLeft {
            dim = 0
        }

        
    }

    for i, spec := range specs {
        spaceLeft := totalDim - pos

        var dim int   
        if spec.IsFixed() {
            dim = spec.FixedSize()
        } else {
            dim = spec.FlexFactor() * flexUnit
        }

        if dim == 0 || {

    }

    return nil
}



