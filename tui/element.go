package tui

import "github.com/gdamore/tcell/v2"

type Element interface {
    // NOTE: an element can be created before it is tied to 
    // an environment. This call is used to add this element to an
    // environment. 
    //
    // If this element contians other elements, it is important those
    // elements recursively join in this call as well.
    Join(parentID ElementID, env *Environment) ElementID

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

    // NOTE: Drawing should be recursive. This should ALWAYS redraw this element, but
    // conditionally redraw the child elements. (If this is a container that is)
    Draw(s tcell.Screen)

    // This function is called when it is time for this element to 
    // leave an environemnt. Like above, make sure all children leave as well!
    //
    // NOTE: We must free our ID in this call!
    Leave(ectx *ElementContext)
}

// Inside an environment, every element has a unique ElementID.
type ElementID int64
const NO_PARENT = -1

type Environment struct {
    elements map[ElementID]*ElementContext

    rootID ElementID 
}

// This function detaches and element from its parent (if it has one)
func (env *Environment) Detach(eid ElementID) {
    ectx, ok := env.elements[eid]
    if !ok {
        return
    }

    // Already detached!
    pid := ectx.parentID
    if pid == NO_PARENT {
        return
    }

    ectx.parentID = NO_PARENT 

    // Our element has no parent pointer now.
    // Now are old parent must have no record of our
    // element.

    parent := env.elements[pid]
    parentCIDs := parent.childIDs

    // eid must be in the child array of parent.
    // Find its index, then perform the removal.
    i := 0
    for ; parentCIDs[i] != eid; i++ {
    }

    parent.childIDs = append(parentCIDs[:i], parentCIDs[i+1:]...)

    parent.self.SetDrawFlag(true)
}

func (env *Environment) MakeRoot(eid ElementID) {
    // Already a root, do nothing.
    if eid == env.rootID {
        return
    }

    ectx, ok := env.elements[eid]
    if !ok {
        return
    }

    env.Detach(eid)

    env.rootID = eid
    ectx.self.SetDrawFlag(true)
}

func (env *Environment) ReserveID() ElementID {
    return 0
}

func (env *Environment) Register(pid ElementID, eid ElementID, cids []ElementID, e Element) {
}

func (env *Environment) Free(eid ElementID) {
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
    if index >= len(e.childIDs) || index < 0 {
        return 0, false
    }

    return e.childIDs[index], true
}

