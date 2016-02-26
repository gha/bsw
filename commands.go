package main

import (
    "strings"
    "github.com/jroimartin/gocui"
)

func moveCursorUp(g *gocui.Gui, v *gocui.View) error {
    return MoveTubeCursor(g, 0, -1)
}

func moveCursorDown(g *gocui.Gui, v *gocui.View) error {
    return MoveTubeCursor(g, 0, 1)
}

func toggleUseTube(g *gocui.Gui, v *gocui.View) error {
    if cTubes.All {
        v, err := g.View("tubes")
        if err != nil {
            return err
        }

        _, cy := v.Cursor()

        tubes := []string{
            cTubes.Names[cy-1],
        }

        cTubes.Use(tubes)
    } else {
        cTubes.UseAll()
    }

    if err := reloadTubes(g); err != nil {
        return err
    }

    if err := reloadMenu(g); err != nil {
        return err
    }

    return nil
}

func toggleCmdMode(g *gocui.Gui, v *gocui.View) error {
    var nv string

    if !cmdMode {
        cmdMode = true
        g.Cursor = true
        nv = "menu"
    } else {
        cmdMode = false
        g.Cursor = false
        nv = "tubes"
    }

    if err := reloadMenu(g); err != nil {
        return err
    }

    return g.SetCurrentView(nv)
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
    cmd := strings.TrimSpace(strings.TrimPrefix(vb, cmdPrefix))
    debugLog("Sent cmd: ", cmd)

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
