package main

import (
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

func quit(g *gocui.Gui, v *gocui.View) error {
    stop <- true

    return gocui.ErrQuit
}
