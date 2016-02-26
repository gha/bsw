package main

import (
    "strings"
    "fmt"
    "github.com/jroimartin/gocui"
)

func moveCursorUp(g *gocui.Gui, v *gocui.View) error {
    if !cmdMode {
        return MoveTubeCursor(g, 0, -1)
    }

    return nil
}

func moveCursorDown(g *gocui.Gui, v *gocui.View) error {
    if !cmdMode {
        return MoveTubeCursor(g, 0, 1)
    }

    return nil
}

func toggleCmdMode(g *gocui.Gui, v *gocui.View) error {
    var nv string

    if !cmdMode {
        cmdMode = true
        g.Cursor = true
        nv = "menu"

        //Clear the tube list for command responses
        tv, err := g.View("tubes")
        if err != nil {
            return err
        }

        tv.Clear()

        PrintCmd(tv, fmt.Sprintf("Running commands on %s:\n\n", cTubes.Selected))
    } else {
        cmdMode = false
        g.Cursor = false
        nv = "tubes"

        //Reload the tube list
        if err := reloadTubes(g); err != nil {
            return err
        }
    }

    if err := reloadMenu(g); err != nil {
        return err
    }

    return g.SetCurrentView(nv)
}

func nextPage(g *gocui.Gui, v *gocui.View) error {
    if cmdMode {
        return nil
    }

    return ChangePage(g, 1)
}

func prevPage(g *gocui.Gui, v *gocui.View) error {
    if cmdMode {
        return nil
    }

    return ChangePage(g, -1)
}

func runCmd(g *gocui.Gui, v *gocui.View) error {
    if !cmdMode {
        return nil
    }

    v, err := g.View("menu")
    if err != nil {
        return err
    }

    vb := v.ViewBuffer()
    cmd := strings.TrimSpace(strings.TrimPrefix(vb, fmt.Sprintf(cmdPrefix, cTubes.Selected)))

    if cmd == "" {
        return nil
    }

    debugLog("Received cmd: ", cmd)

    tv, err := g.View("tubes")
    if err != nil {
        return err
    }

    PrintCmd(tv, fmt.Sprintf("%s:\n", cmd))

    return reloadMenu(g)
}

func exitCmdMode(g *gocui.Gui, v *gocui.View) error {
    if cmdMode {
        return toggleCmdMode(g, v)
    }

    return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
    stop <- true

    return gocui.ErrQuit
}
