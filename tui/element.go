package tui

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
)

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

    // NOTE: Drawing should be recursive. This should ALWAYS redraw this element, and
    // recursively/conditionally redraw the child elements. (Depending on how this element is intended
    // to work)
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

func (env *Environment) GetElement(eid ElementID) (*ElementContext, error) {
    ectx, ok := env.elements[eid] 
    if !ok {
        return nil, fmt.Errorf("Unknown element id: %d", eid)
    }

    return ectx, nil
}

// Attach on element to another as a parent/child.
// This by default is at the end, returns index of child in 
// parents children array.
// 
// If eid points to an element with a parent, this function does nothing.
func (env *Environment) Attach(pid ElementID, eid ElementID) (int, error) {
    return env.attach(pid, eid, false, 0)
}

// Function for attaching at a specific index.
func (env *Environment) AttachAt(pid ElementID, eid ElementID, index int) error {
    _, err := env.attach(pid, eid, true, index)
    return err
}

func (env *Environment) attach(pid ElementID, eid ElementID, at bool, index int) (int, error) {
    ectx, err := env.GetElement(eid)
    if err != nil {
        return 0, err
    }

    if ectx.parentID != NO_PARENT {
        return 0, fmt.Errorf("Element already has parent: %d, %d", 
            ectx.parentID, eid)
    }

    pctx, err := env.GetElement(pid)
    if err != nil {
        return 0, err
    }

    cidsLen := len(pctx.childIDs)

    if !at {
        // Perform Attach.
        ectx.parentID = pid
        pctx.childIDs = append(pctx.childIDs, eid)

        return cidsLen, nil
    }

    // Otherwise we use index!
    if index < 0 || cidsLen < index {
        return 0, fmt.Errorf("Bad index given: %d", index)
    }

    pctx.childIDs = append(pctx.childIDs, 0)
    for i := cidsLen; i > index ; i-- {
        pctx.childIDs[i] = pctx.childIDs[i-1]
    }
    pctx.childIDs[index] = eid

    return index, nil
}

// This function detaches and element from its parent (if it has one)
func (env *Environment) Detach(eid ElementID) error {
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

