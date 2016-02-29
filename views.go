package main

import (
    "fmt"
    "math"
    "strings"
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

        //Calculate the size for paging
        _, vy := v.Size()
        cTubes.Pages = int(math.Ceil(float64(len(cTubes.Conns)) / float64(vy - 1)))
        offset := vy * (cTubes.Page - 1)
        limit := vy * cTubes.Page
        if limit > len(cTubes.Conns) {
            limit = len(cTubes.Conns)
        }
        displayed := cTubes.Conns[offset:limit]

        for _, tube := range displayed {
            stats, _ := tube.Stats()
            jobStats := stats["current-jobs-ready"] + " / " + stats["current-jobs-delayed"] + " / " + stats["current-jobs-buried"]
            workerStats := stats["current-waiting"] + " / " + stats["current-watching"] + " / " + stats["current-using"]
            line := fmt.Sprintf("%-35s %-22s %-22s", tube.Name, jobStats, workerStats)
            fmt.Fprintln(v, line)
        }
}

func PrintString(v *gocui.View, s string) {
    fmt.Fprintf(v, s)
}

func PrintLine(v *gocui.View, line string) {
    fmt.Fprintf(v, fmt.Sprintf("%s\n", line))
}

func PrintMenu(v *gocui.View) {
    v.Editable = true

    menuItems := []interface{}{
        "Exit (Ctrl C)",
        "Toggle Cmd Mode (Ctrl T)",
    }

    if cTubes.Page < cTubes.Pages {
        menuItems = append(menuItems, "Next Page (Ctrl N)")
    }

    if cTubes.Page > 1 {
        menuItems = append(menuItems, "Prev Page (Ctrl P)")
    }

    if !cmdMode {
        verbs := []string{}
        for _, _ = range menuItems {
            verbs = append(verbs, "%s")
        }
        line := fmt.Sprintf(strings.Join(verbs, " | "), menuItems...)
        fmt.Fprintln(v, line)
    } else {
        prefix := fmt.Sprintf(cmdPrefix, cTubes.Selected)
        fmt.Fprintln(v, prefix)
    }
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
