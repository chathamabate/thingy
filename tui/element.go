package tui

import (

	"github.com/gdamore/tcell/v2"
)

type Element interface {

    Resize(ectx *ElementContext, r, c int, rows, cols int)

    // All non-resize events are forwarded through this call.
    // NOTE: This call should never block or result in an infinite looping of
    // events.
    // 
    // Returns true to request a redraw.
    HandleEvent(ectx *ElementContext, ev tcell.Event) error
    
    // Every element should have a boolean "drawFlag" when this flag is true,
    // there needs to be a redraw of THIS element. (NOT NECESSARILY its children)
    SetDrawFlag(f bool)

    // If this is a container element, this SHOULD NOT return just the draw flag.
    // It should recursively OR together the draw flag with RedrawNeeded() from the 
    // child elements.
    RedrawNeeded() bool

    // NOTE: Drawing should be recursive. This should ALWAYS redraw this element, and
    // recursively/conditionally redraw the child elements. (Depending on how this element is intended
    // to work)
    Draw(s tcell.Screen)

    // When an element is deregistered, stop is called.
    //
    // NOTE: This should not be recursive.
    // Stop will automatically be called on all child elements as well
    // during a deregister.
    Stop()
}


