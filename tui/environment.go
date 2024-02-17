package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
    "time"
)


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

    // The screen this Environment draws to.
    screen tcell.Screen

    // An 1 update tick event will be sent to the root
    // per updateDur.
    updateDur time.Duration

    // When this is set to true, the environment will exit its
    // run call this cycle.
    exitRequested bool
}

func NewEnvironment(s tcell.Screen, mc int, ud time.Duration) *Environment {
    return &Environment{
        elements: make([]*EnvEntry, 10),
        maxCapacity: mc,
        fill: 0,
        ptrID: 0,
        rootID: NULL_EID,
        screen: s,
        updateDur: ud,
        exitRequested: false,
    }
}

func (env *Environment) CreateAndRegister(f ElementFactory) (ElementID, error) {
    eid, err := f(env)

    if err != nil {
        return -1, fmt.Errorf("CreateAndRegister: %w", err)
    }

    return eid, nil
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
            children: make([]ChildContext, 0),
        },
        e: e,
    }

    env.fill++
    
    e.Start()

    return eid, nil
}

func (env *Environment) RequestExit() {
    env.exitRequested = true
}

func (env *Environment) getEnvEntry(eid ElementID) (*EnvEntry, error) {
    if eid < 0 || len(env.elements) <= int(eid) {
        return nil, fmt.Errorf("getEnvEntry: Out of bounds element id: %d", eid)
    }

    ee := env.elements[eid] 
    if ee == nil {
        return nil, fmt.Errorf("getEnvEntry: Unknown element id: %d", eid)
    }

    return ee, nil
}

func (env *Environment) GetElementContext(eid ElementID) (*ElementContext, error) {
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return nil, fmt.Errorf("GetElementContext: %w", err)
    }

    return ee.ectx, nil
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

    cctx := ChildContext{
        id: eid,
        attrs: make(map[string]interface{}),
    }

    // Map parent to its new child (at the right index)
    cLen := len(pctx.children)

    if !at {
        pctx.children = append(pctx.children, cctx)

        return cLen, nil
    }

    // Otherwise we use index!
    if index < 0 || cLen < index {
        return 0, fmt.Errorf("attach: Bad index given: %d", index)
    }

    pctx.children = append(pctx.children, ChildContext{})
    for i := cLen; i > index ; i-- {
        pctx.children[i] = pctx.children[i-1]
    }
    pctx.children[index] = cctx

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
    parentChildren := pctx.children

    // eid must be in the child array of parent.
    // Find its index, then perform the removal.
    i := 0
    for ; parentChildren[i].id != eid; i++ {
    }

    pctx.children = append(parentChildren[:i], parentChildren[i+1:]...)

    return nil
}

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

    // When something is made a root, it must be resized to fit the current
    // screen.
    cols, rows := env.screen.Size() 
    err = env.ForwardResize(env.rootID, 0, 0, rows, cols)

    return nil
}

// Child Attribute Functions.

func (env *Environment) getChildAttrs(eid ElementID, childIndex int) (map[string]interface{}, error) {
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return nil, fmt.Errorf("getChildAttrs: %w", err)
    }

    cLen := len(ee.ectx.children)
    if childIndex < 0 || cLen <= childIndex {
        return nil, fmt.Errorf("getChildAttrs: Bad childIndex: %d", childIndex)
    }

    return ee.ectx.children[childIndex].attrs, nil
}

func (env *Environment) GetChildAttr(eid ElementID, childIndex int, key string) (interface{}, error) {
    attrs, err := env.getChildAttrs(eid, childIndex) 
    if err != nil {
        return nil, fmt.Errorf("GetChildAttr: %w", err)
    }

    attr, ok := attrs[key]
    if !ok {
        return nil, fmt.Errorf("GetChildAttr: Key not found %s", key)
    }

    return attr, nil
}

func (env *Environment) SetChildAttr(eid ElementID, childIndex int, key string, val interface{}) error {
    attrs, err := env.getChildAttrs(eid, childIndex) 
    if err != nil {
        return fmt.Errorf("SetChildAttr: %w", err)
    }

    attrs[key] = val

    return nil
}

func (env *Environment) DeleteChildAttr(eid ElementID, childIndex int, key string) error {
    attrs, err := env.getChildAttrs(eid, childIndex) 
    if err != nil {
        return fmt.Errorf("DeleteChildAttr: %w", err)
    }

    delete(attrs, key)

    return nil
}

// Sizing Stuff.

func (env *Environment) GetWidth(eid ElementID) (int, error) {
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return 0, fmt.Errorf("GetWidth: %w", err)
    }

    return ee.e.GetWidth(ee.ectx), nil
}

func (env *Environment) SetWidth(eid ElementID, w int) error {
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return fmt.Errorf("SetWidth: %w", err)
    }

    err = ee.e.SetWidth(ee.ectx, w)

    if err != nil {
        return fmt.Errorf("SetWidth: %w", err)
    }

    ee.e.SetDrawFlag(true)
    return nil
}

func (env *Environment) GetHeight(eid ElementID) (int, error) {
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return 0, fmt.Errorf("GetHeight: %w", err)
    }

    return ee.e.GetHeight(ee.ectx), nil
}

func (env *Environment) SetHeight(eid ElementID, h int) error {
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return fmt.Errorf("SetHeight: %w", err)
    }

    err = ee.e.SetHeight(ee.ectx, h)

    if err != nil {
        return fmt.Errorf("SetHeight: %w", err)
    }

    ee.e.SetDrawFlag(true)
    return nil
}

func (env *Environment) SetViewport(eid ElementID, vp Viewport) error {
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return fmt.Errorf("SetViewport: %w", err)
    }

    err = ee.e.SetViewport(ee.ectx, vp)

    if err != nil {
        return fmt.Errorf("SetViewport: %w", err)
    }

    ee.e.SetDrawFlag(true)
    return nil
}

// Event Forwarding Functions.

func (env *Environment) ForwardEvent(eid ElementID, ev tcell.Event) error {
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return fmt.Errorf("ForwardEvent: %w", err)
    }

    err = ee.e.HandleEvent(ee.ectx, ev)
    if err != nil {
        return fmt.Errorf("ForwardEvent: %w", err)
    }

    return nil
}

func (env *Environment) SetDrawFlag(eid ElementID) error {
    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return fmt.Errorf("SetDrawFlag: %w", err)
    }

    ee.e.SetDrawFlag(true)
    return nil
}

// This returns true if and only if Draw was called on at least one element.
// This begins drawing starting at the root. Then going down.
func (env *Environment) Draw() bool {
    if env.rootID != NULL_EID {
        return env.draw(env.rootID)
    }

    return false
}

// Draw recursive helper.
func (env *Environment) draw(eid ElementID) bool {
    ee := env.elements[eid]

    drawOccured := false    

    // Draw parent first.
    if ee.e.GetDrawFlag() {
        ee.e.Draw(env.screen)
        ee.e.SetDrawFlag(false)
        drawOccured = true
    }

    // Next draw children.
    for _, cctx := range ee.ectx.children {
        drawOccured = env.draw(cctx.id) || drawOccured
    }

    return drawOccured
}

// NOTE: Unlike register, this Deregister is recursive.
// You cannot Deregister an element which has a parent!
func (env *Environment) Deregister(eid ElementID) error {
    if env.rootID == eid {
        return fmt.Errorf("Deregister: Cannot deregister root: %d", eid)
    }

    ee, err := env.getEnvEntry(eid)
    if err != nil {
        return fmt.Errorf("Deregister: %w", err)
    }
    ectx := ee.ectx

    if ectx.parentID != NULL_EID {
        return fmt.Errorf("Deregister: Cannot deregsiter element with parent: %d, %d", 
            ectx.parentID, eid)
    }

    // Now, we first deregister all children.
    for _, cctx := range ectx.children {
        cid := cctx.id

        // sever parent tie. 
        env.elements[cid].ectx.parentID = NULL_EID
        env.Deregister(cid) // this should always succeed.
    }

    ee.e.Stop()

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

// This deregisters all elements in the Environment!
// Essenstially a clean up call.
func (env *Environment) DeregisterAll() error {
    // Make sure to clear the root.
    env.rootID = NULL_EID

    for i := range env.elements {
        ee := env.elements[i]

        // Only deregister elements with no parents.
        // Others will be dealt with recursively.
        if ee.ectx.parentID == NULL_EID {
            err := env.Deregister(ElementID(i))

            // A single error stops the whole thing.
            if err != nil {
                return fmt.Errorf("DeregisterAll: %w", err)
            }
        }
    }

    return nil
}

type UpdateTickEvent struct {
    at time.Time
}

func NewUpdateTickEvent() *UpdateTickEvent {
    return &UpdateTickEvent {
        at: time.Now(),
    }
}

func (u UpdateTickEvent) When() time.Time {
    return u.at
}

// NOTE: UI Loop organization:
//
// 1) Store the duration of the previous iteration
//    including sleep time if needed.
//
// 2) Synchronosly process all queued tcell events through
//    through the root.
//
// 3) Calculate how many ticks occured during the elapsed time of the 
//    last iteration. Send that many update events through the root.
//
// 4) Draw!

func (env *Environment) Run() error {
    // Clear our screen before doing anything else.
    env.screen.Clear()
    env.screen.Show()

    env.exitRequested = false
    var err error

    iterStartTime := time.Now()
    expIterDur := env.updateDur.Milliseconds()

    for {
        // This is the time it took to complete the last iteration.
        lastIterDur := time.Since(iterStartTime).Milliseconds()
        extraTime := expIterDur - lastIterDur

        // If our iteration finished faster than expected.
        // Let's sleep!
        if extraTime > 0 {
            time.Sleep(time.Duration(extraTime * int64(time.Millisecond)))

            // Recalc iteration duration after sleeping.
            lastIterDur = time.Since(iterStartTime).Milliseconds()
        }

        // Now we start the next iteration!
        iterStartTime = time.Now()

        // First poll for system events.
        for env.screen.HasPendingEvent() {
            e := env.screen.PollEvent()

            switch ev := e.(type) {
            case *tcell.EventResize:
                cols, rows := ev.Size()
                err = env.ForwardResize(env.rootID, 0, 0, rows, cols)
                break

            case *tcell.EventKey:
                if ev.Key() == tcell.KeyCtrlC {
                    env.exitRequested = true
                    break
                }

                err = env.ForwardEvent(env.rootID, ev)
                break
            default:
                err = env.ForwardEvent(env.rootID, e)
                break
            }

            if env.exitRequested {
                return nil
            }

            if err != nil {
                return fmt.Errorf("Run: %w", err)
            }
        }

        // Now let's send our update ticks.
        ticksPassed := int(lastIterDur / expIterDur)
        for i := 0; i < ticksPassed; i++ {
            err = env.ForwardEvent(env.rootID, NewUpdateTickEvent())
            if err != nil {
                return fmt.Errorf("Run: %w", err)
            }
        }

        // Finally, time to draw!
        if env.Draw() {
            env.screen.Show()
        }
    }
}

