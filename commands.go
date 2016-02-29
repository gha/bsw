package main

import (
    "strings"
    "errors"
    "strconv"
    "github.com/jroimartin/gocui"
)

type Cmd struct {
    CmdString string
    Handler   CmdHandler
    Args      []string
}

type CmdHandler func(*gocui.View, []string) error

var jobStates = map[string]bool{
    "ready":   true,
    "delayed": true,
    "buried":  true,
}

func ParseCmd(c string) (cmd Cmd, err error) {
    parts := strings.Split(c, " ")

    switch parts[0] {
    case "help":
        if len(parts) != 1 {
            return cmd, errors.New("Invalid command. Usage 'help'")
        }

        cmd.CmdString = c
        cmd.Handler   = Help
        cmd.Args      = []string{}
    case "clear":
        if len(parts) != 2 || !jobStates[parts[1]] {
            return cmd, errors.New("Invalid command. Usage 'clear <ready/delayed/buried>'")
        }

        cmd.CmdString = c
        cmd.Handler   = ClearTube
        cmd.Args      = []string{parts[1]}
    case "next":
        if len(parts) != 2 || !jobStates[parts[1]] {
            return cmd, errors.New("Invalid command. Usage 'next <ready/delayed/buried>'")
        }

        cmd.CmdString = c
        cmd.Handler   = NextJob
        cmd.Args      = []string{parts[1]}
    default:
        return cmd, errors.New("Invalid command. Type help for a list of commands")
    }

    return cmd, err
}

func (c *Cmd) Run(v *gocui.View) error {
    return c.Handler(v, c.Args)
}

func MoveTubeCursor(g *gocui.Gui, mx, my int) error {
    if cmdMode {
        return nil
    }

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
    cTubes.SelectedIdx = ny-1
    cTubes.Selected = cTubes.Names[cTubes.SelectedIdx]
    debugLog("Set tube to: ", cTubes.Selected)

    return nil
}

func ChangePage(g *gocui.Gui, d int) error {
    if cmdMode {
        return nil
    }

    debugLog("Changing page ", cTubes.Page, " by ", d)
    if cTubes.Page < cTubes.Pages && d > 0 {
        cTubes.Page ++
    } else if cTubes.Page > 1 && d < 0 {
        cTubes.Page --
    }

    if err := reloadTubes(g); err != nil {
        return err
    }

    if err := reloadMenu(g); err != nil {
        return err
    }

    return nil
}

func Help(v *gocui.View, _ []string) error {
    debugLog("Showing command list")

    PrintLine(v, "")
    PrintLine(v, "help - Displays a list of commands")
    PrintLine(v, "next <ready/delayed/buried> - Gets the next jobs of the given state")
    PrintLine(v, "clear <ready/delayed/buried> - Clears all jobs of the given state")

    return nil
}

func ClearTube(_ *gocui.View, a []string) error {
    debugLog("Clearing ", a[0], " queue on tube ", cTubes.Selected)

    return nil
}

func NextJob(v *gocui.View, a []string) error {
    debugLog("Getting next ", a[0], " job on tube ", cTubes.Selected)

    var id uint64
    var body []byte
    var err error

    switch a[0] {
    case "ready":
        id, body, err = cTubes.Conns[cTubes.SelectedIdx].PeekReady()
    case "delayed":
        id, body, err = cTubes.Conns[cTubes.SelectedIdx].PeekDelayed()
    case "buried":
        id, body, err = cTubes.Conns[cTubes.SelectedIdx].PeekBuried()
    }
    if err != nil {
        return err
    }

    PrintString(v, strconv.FormatUint(id, 10))
    PrintLine(v, "")
    PrintLine(v, string(body))

    return nil
}
