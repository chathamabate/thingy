package tui

import "github.com/gdamore/tcell/v2"



type Element interface {
    Resize(ectx *ElementContext, r, c int, rows, cols int)

    // All non-resize events are forwarded through this call.
    // NOTE: This call should never block or result in an infinite looping of
    // events.
    // 
    // Returns true to request a redraw.
    HandleEvent(ectx *ElementContext, ev tcell.Event) (bool, error) 
}

// Inside an environment, every element has a unique ElementID.
type ElementID uint64

type Environment struct {
    elements map[ElementID]*ElementContext
    rootID ElementID 
}

type ElementContext struct {
    env *Environment
    self Element
    
    parentID ElementID
    selfID ElementID 
    childIDs []ElementID
}

func (e *ElementContext) Env() *Environment {
    return e.env
}

func (e *ElementContext) ID() ElementID {
    return e.selfID
}

func (e *ElementContext) ParentID() ElementID {
    return e.parentID
}

func (e *ElementContext) ChildID(index int) (ElementID, bool) {
    if index >= len(e.childIDs)  || index < 0 {
        return 0, false
    }

    return e.childIDs[index], true
}

