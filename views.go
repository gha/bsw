package main

import (
    "fmt"
    "github.com/jroimartin/gocui"
)

func PrintTubeList(v *gocui.View) {
        line := fmt.Sprintf("%-35s %-22s %-22s", "Tube", "ready/delayed/buried", "waiting/watching/using")
        fmt.Fprintln(v, line)

        v.Highlight = true
        v.Wrap      = true
        v.Editable  = false

        //Reload the tube stats - will detect new tubes and drop removed tubes
        cTubes.UseAll()

        for _, tube := range cTubes.Conns {
            stats, _ := tube.Stats()
            jobStats := stats["current-jobs-ready"] + " / " + stats["current-jobs-delayed"] + " / " + stats["current-jobs-buried"]
            workerStats := stats["current-waiting"] + " / " + stats["current-watching"] + " / " + stats["current-using"]
            line := fmt.Sprintf("%-35s %-22s %-22s", tube.Name, jobStats, workerStats)
            fmt.Fprintln(v, line)
        }
}

func PrintCmd(v *gocui.View, line string) {
    fmt.Fprintf(v, line)
}

func PrintMenu(v *gocui.View) {
    v.Editable = true

    if !cmdMode {
        line := fmt.Sprintf("%s | %s", "Exit (Ctrl C)", "Toggle Cmd Mode (Ctrl T)")
        fmt.Fprintln(v, line)
    } else {
        prefix := fmt.Sprintf(cmdPrefix, cTubes.Selected)
        fmt.Fprintln(v, prefix)
    }
}

func MoveTubeCursor(g *gocui.Gui, mx, my int) error {
    tv, err := g.View("tubes")
    if err != nil {
        return err
    }

    cx, cy := tv.Cursor()
    ny := cy + my

    //Check the cursor isn't trying to move above the first tube or past the last tube
    if ny < 1 || ny > len(cTubes.Conns) {
        return nil
    }

    if err = tv.SetCursor(cx, ny); err != nil {
        return err
    }

    //Set the selected tube to the currently highlighted row
    cTubes.Selected = cTubes.Names[ny-1]
    debugLog("Set tube to: ", cTubes.Selected)

    return nil
}

func RefreshCursor(g *gocui.Gui) error {
    tv, err := g.View("tubes")
    if err != nil {
        return err
    }

    _, cy := tv.Cursor()

    if cy > len(cTubes.Conns) {
        debugLog("Resetting cursor position ", cy, " to ", len(cTubes.Conns))

        //Temporary fix for the cursor dropping off the bottom of the list
        if err = tv.SetCursor(0, len(cTubes.Conns)); err != nil {
            return err
        }
    }

    return nil
}
