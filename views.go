package main

import (
    "fmt"
    "github.com/jroimartin/gocui"
)

func PrintTubeList(v *gocui.View) {
        //Reload the tube stats
        cTubes.UseAll()

        line := fmt.Sprintf("%-40s %-30s %-30s", "Tube", "ready/delayed/buried", "waiting/watching/using")
        fmt.Fprintln(v, line)

        for _, tube := range cTubes.Conns {
            stats, _ := tube.Stats()
            jobStats := stats["current-jobs-ready"] + "/" + stats["current-jobs-delayed"] + "/" + stats["current-jobs-buried"]
            workerStats := stats["current-waiting"] + "/" + stats["current-watching"] + "/" + stats["current-using"]
            line := fmt.Sprintf("%-40s %-30s %-30s", tube.Name, jobStats, workerStats)
            fmt.Fprintln(v, line)
        }
}

func PrintMenu(v *gocui.View) {
    line := fmt.Sprintf("%s", "Exit (Ctrl C)")
    fmt.Fprintln(v, line)
}
