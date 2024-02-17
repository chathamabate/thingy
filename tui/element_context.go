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

    children []ChildContext
}

type ChildContext struct {
    id ElementID
     
    // Each child can have arbitrary attributes?
    attrs map[string]interface{}
}

func (ectx *ElementContext) Parent() (*ElementContext, error) {
    if ectx.parentID == NULL_EID {
        return nil, errors.New("GetParent: Element has no parent")
    }

    // This should NEVER error as parent is a valid ID.
    return ectx.env.GetElementContext(ectx.selfID)
}

func (ectx *ElementContext) Child(index int) (*ElementContext, error) {
    if index < 0 || len(ectx.children) <= index {
        return nil, fmt.Errorf("GetChild: Bad index: %d", index)
    }

    return ectx.env.GetElementContext(ectx.children[index].id)
}

func (ectx *ElementContext) GetChildAttr(index int, key string) (interface{}, error) {
    if index < 0 || len(ectx.children) <= index {
        return nil, fmt.Errorf("GetChildAttr: Bad index: %d", index)
    }

    val, ok := ectx.children[index].attrs[key]
    if !ok {
        return nil, fmt.Errorf("GetChildAttr: Bad key: %s", key) 
    }

    return val, nil
}

func (ectx *ElementContext) SetChildAttr(index int, key string, value interface{}) error {
    if index < 0 || len(ectx.children) <= index {
        return fmt.Errorf("SetChildAttr: Bad index: %d", index)
    }

    ectx.children[index].attrs[key] = value

    return nil
}

func (ectx *ElementContext) DeleteChildAttr(index int, key string) error {
    if index < 0 || len(ectx.children) <= index {
        return fmt.Errorf("DeleteChildAttr: Bad index: %d", index)
    }
    
    delete(ectx.children[index].attrs, key)

    return nil
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
    return len(ectx.children)
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




