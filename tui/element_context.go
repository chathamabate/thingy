package tui

type ElementContext struct {
    env *Environment
    
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
