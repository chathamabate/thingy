package tui

import (
	"time"

	tcell "github.com/gdamore/tcell/v2"
)

type EventQueue interface {
    PostEvent(ev tcell.Event) error 
}

type Comp interface {
    // Life cycle functions.

    // This should only ever be called once.
    Init() error

    // NOTE: Resize is promised to be called once right after init.
    // Every time the terminal is resized afterwards, this will be called.
    Resize(r, c int, rows, cols int)

    // ev will never be a resize event.
    // This will always be passed to the above resize function.
    ForwardEvent(ev tcell.Event) error
    Update() error

    // Draw is called every cycle.
    // Returns true only when drawing actually occurred!
    Draw(s tcell.Screen) bool

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
    err = root.Init()
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
                cols, rows := ev.Size()
                root.Resize(0, 0, rows, cols)
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

        if root.Draw(s) {
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
    R int
    C int

    Rows int
    Cols int

    RedrawNeeded bool
}

func (cd *CompDefault) ForwardEvent(ev tcell.Event) error { 
    return nil
}

func (cd *CompDefault) Resize(r, c int, rows, cols int) { 
    cd.R = r
    cd.C = c

    cd.Rows = rows
    cd.Cols = cols

    cd.RedrawNeeded = true
}

func (cd *CompDefault) Init() error {
    return nil
}

func (cd *CompDefault) Update() error {
    return nil
}

func (cd *CompDefault) Draw(s tcell.Screen) bool {
    return false
}

func (cd *CompDefault) Cleanup() {

}

// Example Plain Text Component.
type CompPlainText struct {
    CompDefault
    label string
}

func NewCompPlainText(l string) *CompPlainText {
    return &CompPlainText{label: l} 
}

func (cpt *CompPlainText) Draw(s tcell.Screen) bool {
    if !cpt.RedrawNeeded {
        return false
    }

    st := tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite)

    for r := 0; r < cpt.Rows; r++ {
        for c := 0; c <  cpt.Cols; c++ {
            s.SetContent(c + cpt.C, r + cpt.R, ' ', nil, st)
        }
    }

    runeNum := 0
    for _, ru := range cpt.label {
        r := runeNum / cpt.Cols
        if r >= cpt.Rows {
            break
        }

        c := runeNum % cpt.Cols
        
        s.SetContent(c + cpt.C, r + cpt.R, ru, nil, st)
        runeNum++
    }

    cpt.RedrawNeeded = false
    return true
}





