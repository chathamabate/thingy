package tui

import (
	"errors"

	"github.com/gdamore/tcell/v2"
)


type Viewport struct {
    // Coordinates of where the viewport should be 
    // rendered on screen.
    x, y int

    // Offset into the element which is being rendered.
    xOff, yOff int

    // Dimmensions of the viewport.
    width, height int
}

// An element factory is meant to create an element then register
// it into the environment.
// Returning the element's new ID.
//
// NOTE: I am using the factory design pattern here to encourage the
// user to never have direct access to an element which is part of
// an environment.
type ElementFactory func (*Environment) (ElementID, error)

type Element interface {
    // When an element is registered, start is called.
    //
    // NOTE: This should NOT be recursive.
    // When child elements are registered, start will be called for them.
    // This call should only deal with what must occur for this 
    // element alone, NOT its children.
    Start()

    // This should resize the given element.
    //

    // Every element will have a width and height.
    // SetWidth/Height may throw errors, since some elements may not allow
    // for such operations.
    // For example an element which grows downward based on the elements it
    // contains may not allow for the user to randomly set its height.
    // Using SetWidth/Height always triggers the draw flag (If there is no error that is)

    // NOTE: The below 3 Set calls will all automatically trigger a redraw for this
    // element. They SHOULD be recursive, if child elements must also be resized somehow.

    GetWidth(ectx *ElementContext) int
    SetWidth(ectx *ElementContext, w int) error

    GetHeight(ectx *ElementContext) int
    SetHeight(ectx *ElementContext, h int) error

    SetViewport(ectx *ElementContext, vp Viewport) error

    // All non-resize events are forwarded through this call.
    // NOTE: This call should never block or result in an infinite looping of
    // events.
    HandleEvent(ectx *ElementContext, ev tcell.Event) error
    
    // Every element should have a boolean "drawFlag" when this flag is true,
    // there needs to be a redraw of THIS element. (NOT NECESSARILY its children)
    SetDrawFlag(f bool)
    GetDrawFlag() bool

    // This always draws just THIS element. Nothing recursive here.
    // If Drawing this element requires the redrawing of children or vice versa,
    // make sure to declare this explicitly in handle event somewhere.
    //
    // NOTE: Do not mess with the draw flag here... that is handled entirely
    // by the environment.
    // 
    // NOTE: See environment.Draw.
    Draw(s tcell.Screen)

    // When an element is deregistered, stop is called.
    //
    // NOTE: This should not be recursive.
    // Stop will automatically be called on all child elements as well
    // during a deregister.
    Stop()
}

type DefaultElement struct {
    r, c int
    rows, cols int
    drawFlag bool
}

func NewDefaultElement() *DefaultElement {
    return &DefaultElement{
        r: 0, 
        c: 0,
        rows: 0,
        cols: 0,
        drawFlag: false,
    }
}

func (de *DefaultElement) Start() {
}

func (de *DefaultElement) Resize(ectx *ElementContext, r, c int, rows, cols int) error {
    if r < 0 || c < 0 || rows < 0 || cols < 0 {
        return errors.New("Resize: Negaitive dimmension given")
    }

    de.r = r
    de.c = c
    de.rows = rows
    de.cols = cols

    return nil
}

func (de *DefaultElement) GetR() int {
    return de.r
}

func (de *DefaultElement) GetC() int {
    return de.c
}

func (de *DefaultElement) GetRows() int {
    return de.rows
}

func (de *DefaultElement) GetCols() int {
    return de.cols
}

func (de *DefaultElement) HandleEvent(ectx *ElementContext, ev tcell.Event) error {
    return nil
}

func (de *DefaultElement) SetDrawFlag(f bool) {
    de.drawFlag = f
}

func (de *DefaultElement) GetDrawFlag() bool {
    return de.drawFlag
}

func (de *DefaultElement) Draw(s tcell.Screen) {
    for i := de.r; i < de.r + de.rows; i++ {
        for j := de.c; j < de.c + de.cols; j++ {
            s.SetContent(j, i, ' ', nil, tcell.StyleDefault)
        }
    }
}

func (de *DefaultElement) Stop() {
}



