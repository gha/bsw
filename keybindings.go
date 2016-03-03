package main

import (
    "strings"
    "fmt"
    "github.com/jroimartin/gocui"
)

func setKeyBindings(g *gocui.Gui) error {
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, moveCursorUp); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, moveCursorDown); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, nextPage); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, prevPage); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, toggleCmdMode); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, runCmd); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, exitCmdMode); err != nil {
        return err
    }

    return nil
}

func moveCursorUp(g *gocui.Gui, v *gocui.View) error {
    return MoveTubeCursor(g, 0, -1)
}

func moveCursorDown(g *gocui.Gui, v *gocui.View) error {
    return MoveTubeCursor(g, 0, 1)
}

func nextPage(g *gocui.Gui, v *gocui.View) error {
    return ChangePage(g, 1)
}

func prevPage(g *gocui.Gui, v *gocui.View) error {
    return ChangePage(g, -1)
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

        PrintLine(tv, fmt.Sprintf("Running commands on %s:\n", cTubes.Selected))
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

func runCmd(g *gocui.Gui, v *gocui.View) error {
    if !cmdMode {
        return nil
    }

    v, err := g.View("menu")
    if err != nil {
        return err
    }

    vb := v.ViewBuffer()
    cmdString := strings.TrimSpace(strings.TrimPrefix(vb, fmt.Sprintf(cmdPrefix, cTubes.Selected)))

    if cmdString == "" {
        return nil
    }

    debugLog("Received cmd: ", cmdString)

    tv, err := g.View("tubes")
    if err != nil {
        return err
    }

    PrintString(tv, cmdString + ": ")

    cmd, err := ParseCmd(cmdString)
    if err != nil {
        PrintLine(tv, err.Error())
    } else {
        if err := cmd.Run(tv); err != nil {
            PrintLine(tv, err.Error())
        }
    }

    PrintLine(tv, "------------")

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
