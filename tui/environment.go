package tui

import "fmt"

// An environment is kinda like the DOM in javascript.
// It is a single object which stores the relationships between all 
// elements in the UI.
//
// NOTE: By design an environment is in no way thread-safe.
// All actions on an environment are supposed to be synchronous and 
// non-blocking.

type ElementID int
const NULL_EID = -1

type EnvEntry struct {
    ectx *ElementContext
    e Element
}

type Environment struct {
    // Map containing all elements in the entire environment.
    elements []*EnvEntry

    // Size of the elements map will never outgrow maxCapacity.
    maxCapacity int

    // The number of non-nil entries in the map.
    fill int

    // This simply an arbitrary id value, we start our search
    // at this value when looking for a new ID.
    // If fill < len(elements) it is gauranteed that ptrID
    // is a valid index into the elements map.
    //
    // NOTE: see Register.
    ptrID ElementID

    rootID ElementID 
}

// Register adds an element to an environment and returns
// its corresponding Element ID.
func (env *Environment) Register(e Element) (ElementID, error) {
    if env.fill == env.maxCapacity {
        return -1, fmt.Errorf("Register: Environment at max capacity: %d", env.maxCapacity)
    }

    var eid ElementID

    if env.fill == len(env.elements) {
        // In this case our elements map is full, but not yet at
        // it's max capcity. Add an extra spot, make that the id.
        eid = ElementID(len(env.elements))
        env.elements = append(env.elements, nil)
    } else {
        // Otherwise, there is a spot in the map, we just need to 
        // find it.
        for ; env.elements[env.ptrID] != nil; env.ptrID++ {
            if int(env.ptrID) == len(env.elements) {
                env.ptrID = -1  // Will be zero of post action.
            }
        }

        eid = env.ptrID
    }

    env.elements[eid] = &EnvEntry{
        ectx: &ElementContext{
            env: env,
            parentID: NULL_EID,
            selfID: eid,
            childIDs: make([]ElementID, 0),
        },
        e: e,
    }
    
    env.fill++

    return eid, nil
}

func (env *Environment) getEnvEntry(eid ElementID) (*EnvEntry, error) {
    ee := env.elements[eid] 
    if ee == nil {
        return nil, fmt.Errorf("getEnvEntry: Unknown element id: %d", eid)
    }

    return ee, nil
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
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return 0, fmt.Errorf("attach: %w", err)
    }
    ectx := ee.ectx

    if ectx.parentID != NULL_EID {
        return 0, fmt.Errorf("attach: Element already has parent: %d, %d", 
            ectx.parentID, eid)
    }

    pe, err := env.getEnvEntry(pid)
    if err != nil {
        return 0, fmt.Errorf("attach: %w", err)
    }
    pctx := pe.ectx

    // Map our element to its new parent.
    ectx.parentID = pid

    // Map parent to its new child (at the right index)
    cidsLen := len(pctx.childIDs)

    if !at {
        pctx.childIDs = append(pctx.childIDs, eid)

        return cidsLen, nil
    }

    // Otherwise we use index!
    if index < 0 || cidsLen < index {
        return 0, fmt.Errorf("attach: Bad index given: %d", index)
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
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return fmt.Errorf("Detach: %w", err)
    }
    ectx := ee.ectx

    // Already detached!
    pid := ectx.parentID
    if pid == NULL_EID {
        return fmt.Errorf("Detach: Element has no parent: %d", 
            eid)
    }

    // Perform detach!

    ectx.parentID = NULL_EID 

    // Our element has no parent pointer now.
    // Now are old parent must have no record of our
    // element.

    pctx := env.elements[pid].ectx
    parentCIDs := pctx.childIDs

    // eid must be in the child array of parent.
    // Find its index, then perform the removal.
    i := 0
    for ; parentCIDs[i] != eid; i++ {
    }

    pctx.childIDs = append(parentCIDs[:i], parentCIDs[i+1:]...)

    return nil
}

func (env *Environment) ClearRoot()

func (env *Environment) MakeRoot(eid ElementID) error {
    // Clearing the root always works!
    if eid == NULL_EID {
        env.rootID = eid
        return nil
    }

    if env.rootID == eid {
        return fmt.Errorf("MakeRoot: Element already root: %d", eid)
    }

    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return fmt.Errorf("MakeRoot: %w", err)
    }

    if ee.ectx.parentID != NULL_EID {
        return fmt.Errorf("MakeRoot: Element has parent: %d, %d", ee.ectx.parentID, eid)
    }

    env.rootID = eid

    return nil
}

// NOTE: Unlike register, this Deregister is recursive.
// You cannot Deregister an element which has a parent!
func (env *Environment) Deregister(eid ElementID) error {
    if env.rootID == eid {
        return fmt.Errorf("Deregister: Cannot deregister root: %d", eid)
    }

    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return err
    }
    ectx := ee.ectx

    if ectx.parentID != NULL_EID {
        return fmt.Errorf("Deregister: Cannot deregsiter element with parent: %d, %d", 
            ectx.parentID, eid)
    }

    // Now, we first deregister all children.
    for _, cid := range ectx.childIDs {
        // sever parent tie. 
        env.elements[cid].ectx.parentID = NULL_EID
        env.Deregister(cid) // this should always succeed.
    }

    // finally, remove this guy from the env
    env.elements[eid] = nil

    env.fill--

    // If our map was previously full, let's set our
    // pointer id to eid to increase speed of next register.
    if env.fill == len(env.elements) - 1 {
        env.ptrID = eid
    }

    return nil
}

