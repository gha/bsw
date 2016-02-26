package main

import (
    "log"
    "log/syslog"
    "time"
    "flag"
    "github.com/jroimartin/gocui"
    "github.com/kr/beanstalk"
)

var (
    conn *beanstalk.Conn
    cTubes Tubes
    watch = false
    stop = make(chan bool)
    host = flag.String("host", "127.0.0.1:11300", "Beanstalk host address")
    refreshRate = flag.Int("refresh", 1, "Refresh rate of the tube list (seconds)")
    debug = flag.Bool("debug", false, "Enable debug logging")
    logWriter *syslog.Writer
    cmdMode = false
)

const (
    cmdPrefix = ": "
)

func main() {
    var err error

    logWriter, err = syslog.New(syslog.LOG_INFO, "bsw")
    if err != nil {
        log.Fatal(err)
    }
    log.SetOutput(logWriter)

    flag.Parse()

    if conn, err = beanstalk.Dial("tcp", *host); err != nil {
        log.Fatal(err)
    }
    debugLog("Connected to beanstalk")

    //Use all tubes by default
    cTubes.UseAll()

    g := gocui.NewGui()
    if err := g.Init(); err != nil {
        log.Fatal(err)
    }
    defer g.Close()

    setKeyBindings(g)
    debugLog("Set keybindings")

    g.SetLayout(setLayout)
    debugLog("Set layout")
    g.Editor = gocui.EditorFunc(cmdEditor)
    debugLog("Set editor")
    g.Cursor = true
    go watchTubes(g)

    debugLog("Starting main loop")

    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Fatal(err)
    }
}

func cmdEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
    switch {
    case ch != 0 && mod == 0:
        v.EditWrite(ch)
    case key == gocui.KeySpace:
        v.EditWrite(' ')
    case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
        v.EditDelete(true)
    case key == gocui.KeyDelete:
        v.EditDelete(false)
    case key == gocui.KeyArrowLeft:
        v.MoveCursor(-1, 0, false)
    case key == gocui.KeyArrowRight:
        v.MoveCursor(1, 0, false)
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

    if err := g.SetKeybinding("", gocui.KeyCtrlT, gocui.ModNone, toggleCmdMode); err != nil {
        log.Fatal(err)
    }

    if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, runCmd); err != nil {
        log.Fatal(err)
    }

    if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, exitCmdMode); err != nil {
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

    if cmdMode {
        if err = v.SetCursor(2, 0); err != nil {
            return err
        }
    }

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
    //Refresh the cursor
    return RefreshCursor(g)
}

func watchTubes(g *gocui.Gui) {
    for {
        select {
            case <-stop:
                watch = false
                return
            case <-time.After(time.Duration(*refreshRate) * time.Second):
                watch = true
                //Refresh tube list
                g.Execute(func(g *gocui.Gui) error {
                    return reloadTubes(g)
                })

                //Reload the menu if we're not in cmd mode
                if !cmdMode {
                    _ = reloadMenu(g);
                }
        }
    }
}

func debugLog(s ...interface{}) {
    if *debug {
        log.Print(s...)
    }
}
