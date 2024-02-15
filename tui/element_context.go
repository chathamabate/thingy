package tui

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type ElementContext struct {
    env *Environment
    
    parentID ElementID
    selfID ElementID 
    childIDs []ElementID
}

func (ectx *ElementContext) Parent() (*ElementContext, error) {
    if ectx.parentID == NULL_EID {
        return nil, errors.New("GetParent: Element has no parent")
    }

    // This should NEVER error as parent is a valid ID.
    return ectx.env.GetElementContext(ectx.selfID)
}

func (ectx *ElementContext) Child(index int) (*ElementContext, error) {
    if index < 0 || len(ectx.childIDs) <= index {
        return nil, fmt.Errorf("GetChild: Bad index: %d", index)
    }

    return ectx.env.GetElementContext(ectx.childIDs[index])
}

func (ectx *ElementContext) CreateRegisterAndAttach(f ElementFactory) (int, error) {   
    eid, err := ectx.env.CreateAndRegister(f)
    if err != nil {
        return -1, err
    }

    return ectx.env.Attach(ectx.selfID, eid)
}

func (ectx *ElementContext) CreateRegisterAndAttachAt(index int, f ElementFactory) error {
    eid, err := ectx.env.CreateAndRegister(f)
    if err != nil {
        return err
    }

    return ectx.env.AttachAt(ectx.selfID, eid, index)
}

func (ectx *ElementContext) NumChildren() int {
    return len(ectx.childIDs)
}

func (ectx *ElementContext) RequestExit() {
    ectx.env.RequestExit()
}

func (ectx *ElementContext) ForwardResize(r, c int, rows, cols int) error {
    return ectx.env.ForwardResize(ectx.selfID, r, c, rows, cols)
}

func (ectx *ElementContext) ForwardEvent(ev tcell.Event) error {
    return ectx.env.ForwardEvent(ectx.selfID, ev)
}

func (ectx *ElementContext) SetDrawFlag() {
    ectx.env.SetDrawFlag(ectx.selfID)
}

func (ectx *ElementContext) DetachAndDeregister() error {
    err := ectx.env.Detach(ectx.selfID)
    if err != nil {
        return err
    }

    return ectx.env.Deregister(ectx.selfID)
}




