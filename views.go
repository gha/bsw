package main

import (
    "fmt"
    "github.com/jroimartin/gocui"
)

func PrintTubeList(v *gocui.View) {
        //Reload the tube stats
        cTubes.UseAll()

        line := fmt.Sprintf("%-35s %-22s %-22s", "Tube", "ready/delayed/buried", "waiting/watching/using")
        fmt.Fprintln(v, line)

        v.Highlight = true
        v.Wrap      = true

        for _, tube := range cTubes.Conns {
            stats, _ := tube.Stats()
            jobStats := stats["current-jobs-ready"] + " / " + stats["current-jobs-delayed"] + " / " + stats["current-jobs-buried"]
            workerStats := stats["current-waiting"] + " / " + stats["current-watching"] + " / " + stats["current-using"]
            line := fmt.Sprintf("%-35s %-22s %-22s", tube.Name, jobStats, workerStats)
            fmt.Fprintln(v, line)
        }
}

func PrintMenu(v *gocui.View) {
    line := fmt.Sprintf("%s | %s", "Exit (Ctrl C)", "Use Tube (Enter)")
    fmt.Fprintln(v, line)
}

func MoveTubeCursor(g *gocui.Gui, mx, my int) error {
    tv, err := g.View("tubes")
    if err != nil {
        return err
    }

    maxX, maxY := tv.Size()
    //Set the max height to the number of tubes so we cant scroll past the last tube
    maxY        = len(cTubes.Conns)
    //Get the current cursor position
    cx, cy     := tv.Cursor()

    //If the current cursor exceeds the bounds of the view, move it back
    //This usually happens if the bottom tube is highlighed and a tube drops off the list
    if cx > maxX || cy > maxY {
        return tv.SetCursor(0, maxY)
    }

    //Update the cursor with the new position
    cx, cy      = cx + mx, cy + my

    //If the new cursor exceeds the bounds of the view dont move it
    if cx < 0 || cx > maxX || cy < 1 || cy > maxY {
        return nil
    }

    return tv.SetCursor(cx, cy)
}
