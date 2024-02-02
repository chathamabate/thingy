package tui

import (
	"fmt"
	"time"

	tcell "github.com/gdamore/tcell/v2"
)

type EventQueue interface {
    PostEvent(ev tcell.Event) error 
}

type Comp interface {
    GetParent() Comp

    // Life cycle functions.

    // This should only ever be called once.
    // It should set the parent component and initialize child components.
    // NOTE: This component must only ever belong to a SINGLE parent.
    //
    // It should NEVER be appended to two different components.
    Init(p Comp) error

    // NOTE: Resize is promised to be called once right after init.
    // Every time the terminal is resized afterwards, this will be called.
    Resize(rows, cols int)

    // ev will never be a resize event.
    // This will always be passed to the above resize function.
    ForwardEvent(ev tcell.Event) error
    Update() error

    // NOTE: Rendering of any kind should ONLY be done in this funciton.
    // We return true if the render actually did anything.
    Render(r, c int) bool

    // This is called at the end of a Components lifetime.
    // Either when the User exits the program, or when this component goes out of scope.
    Cleanup()
}

func RunTUI(root Comp, updateDur time.Duration) error {
    s, err := tcell.NewScreen()
    if err != nil {
        return err
    }

    err = s.Init()
    if err != nil {
        return err
    }
    defer s.Fini()

    // Initialize our components.
    err = root.Init(nil)
    if err != nil {
        return err
    }
    defer root.Cleanup()

    // We make sure our screen is cleared to start with.
    s.Clear()  
    s.Show()

    tp := time.Now()

    for {
        for s.HasPendingEvent() {
            e := s.PollEvent()
            switch ev := e.(type) {
            case *tcell.EventResize:
                rows, cols := ev.Size()
                root.Resize(rows, cols)
                break
            case *tcell.EventKey:
                if ev.Key() == tcell.KeyCtrlC {
                    return nil
                }

                err = root.ForwardEvent(e)
                break
            default:
                err = root.ForwardEvent(e)
                break
            }

            if err != nil {
                return err
            }
        }

        err = root.Update()
        if err != nil {
            return err
        }

        if root.Render(0, 0) {
            s.Show()
        }

        iterationTime := time.Since(tp)
        extraTime := updateDur.Milliseconds() - iterationTime.Milliseconds()
        if extraTime > 0 {
            time.Sleep(time.Duration(extraTime * int64(time.Millisecond)))
        }

        tp = time.Now()
    }
}

// Default component implementation.
// Has no children, displays nothing.
type CompDefault struct {
    Parent Comp
    Rows int
    Cols int
}

func (c *CompDefault) ForwardEvent(ev tcell.Event) error { 
    return nil
}

func (c *CompDefault) Resize(rows, cols int) { 
    c.Rows = rows
    c.Cols = cols
}

func (c *CompDefault) GetParent() Comp {
    return c.Parent
}

func (c *CompDefault) Init(parent Comp) error {
    c.Parent = parent

    return nil
}

func (c *CompDefault) Update() error {
    return nil
}

func (c *CompDefault) Render(int, int) bool {
    return false
}

func (c *CompDefault) Cleanup() {

}



