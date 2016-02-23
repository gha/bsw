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

        if cTubes.All {
            //Reload the tube stats - will detect new tubes and drop removed tubes
            cTubes.UseAll()
        }

        for _, tube := range cTubes.Conns {
            stats, _ := tube.Stats()
            jobStats := stats["current-jobs-ready"] + " / " + stats["current-jobs-delayed"] + " / " + stats["current-jobs-buried"]
            workerStats := stats["current-waiting"] + " / " + stats["current-watching"] + " / " + stats["current-using"]
            line := fmt.Sprintf("%-35s %-22s %-22s", tube.Name, jobStats, workerStats)
            fmt.Fprintln(v, line)
        }
}

func PrintMenu(v *gocui.View) {
    tubeSelector := "Use Tube (Tab)"
    if !cTubes.All {
        tubeSelector = "Use All (Tab)"
    }

    line := fmt.Sprintf("%s | %s", "Exit (Ctrl C)", tubeSelector)
    fmt.Fprintln(v, line)
}

func MoveTubeCursor(g *gocui.Gui, mx, my int) error {
    tv, err := g.View("tubes")
    if err != nil {
        return err
    }

    cx, cy := tv.Cursor()
    ox, oy := tv.Origin()
    ny := cy + my

    //Check the cursor isn't trying to move above the first tube
    if ny < 1 && oy == 0 {
        return nil
    }

    //Check the cursor isn't trying to move past the last tube
    if ny + oy > len(cTubes.Conns) && ny > cy {
        return nil
    }

    if err = tv.SetCursor(cx, ny); err != nil {
        //If we've moved to an invalid point, update the origin
        if err = tv.SetOrigin(ox, oy + my); err != nil {
            return err
        }
    }

    return nil
}

func RefreshCursor(g *gocui.Gui) error {
    tv, err := g.View("tubes")
    if err != nil {
        return err
    }

    _, cy := tv.Cursor()
    _, oy := tv.Origin()

    if cy + oy > len(cTubes.Conns) {
        debugLog("Resetting cursor: cy: ", cy, " oy: ", oy, " t: ", len(cTubes.Conns))

        //Temporary fix for the cursor dropping off the bottom of the list
        if err = tv.SetCursor(0, 1); err != nil {
            return err
        }

        if err = tv.SetOrigin(0, 0); err != nil {
            return err
        }
    }

    return nil
}
