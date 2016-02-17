package main

import (
    "log"
    "time"
    "github.com/jroimartin/gocui"
    "github.com/kr/beanstalk"
)

var (
    conn *beanstalk.Conn
    cTubes Tubes
    watch = false
    stop = make(chan bool)
)

func main() {
    var err error
    if conn, err = beanstalk.Dial("tcp", "127.0.0.1:11300"); err != nil {
        log.Fatal(err)
    }

    //Use all tubes by default
    cTubes.UseAll()

    g := gocui.NewGui()
    if err := g.Init(); err != nil {
        log.Fatal(err)
    }
    defer g.Close()

    setKeyBindings(g)

    g.SetLayout(setLayout)
    g.Cursor = true
    go watchTubes(g)

    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Fatal(err)
    }
}

func setLayout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := g.SetView("tubes", 0, 0, maxX-1, maxY-3); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }

        PrintTubeList(v)

        //Move the cursor to the first tube
        if err = MoveTubeCursor(g, 0, 1); err != nil {
            return err
        }
    }

    if v, err := g.SetView("menu", 0, maxY-3, maxX-1, maxY-1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }

        PrintMenu(v)
    }

    return nil
}

func setKeyBindings(g *gocui.Gui) {
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        log.Fatal(err)
    }

    if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, moveCursorUp); err != nil {
        log.Fatal(err)
    }

    if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, moveCursorDown); err != nil {
        log.Fatal(err)
    }

    if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, toggleUseTube); err != nil {
        log.Fatal(err)
    }
}

func reloadMenu(g *gocui.Gui) error {
    v, err := g.View("menu")
    if err != nil {
        return err
    }

    v.Clear()
    PrintMenu(v)

    _, err = g.SetViewOnTop("menu")

    return err
}

func reloadTubes(g *gocui.Gui) error {
    v, err := g.View("tubes")
    if err != nil {
        return err
    }

    //Clear the current tube list
    v.Clear()
    //Print the new tube list
    PrintTubeList(v)
    //Check the cursor hasn't fallen off the bottom
    return MoveTubeCursor(g, 0, 0)
}

func watchTubes(g *gocui.Gui) {
    for {
        select {
            case <-stop:
                watch = false
                return
            case <-time.After(1 * time.Second):
                watch = true
                //Refresh tube list
                g.Execute(func(g *gocui.Gui) error {
                    return reloadTubes(g)
                })

                _ = reloadMenu(g);
        }
    }
}
