package tui

import (

	"github.com/gdamore/tcell/v2"
)

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
    // NOTE: if this is a container element, a resize should be
    // called recursively for child elements.
    // 
    // No need to set the draw flag in this function, this is automatically
    // handled by the environment when a resize is forwarded.
    Resize(ectx *ElementContext, r, c int, rows, cols int) error

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


