package main

import (
    "errors"
    "github.com/jroimartin/gocui"
    "strings"
)

type Cmd interface {
    Run(*gocui.View) error
    Validator([]string) error
    SetArgs([]string)
    GetArgs() []string
    Usage() string
    Description() string
}

var validCmds = map[string]Cmd{
    "help": &Help{},
    "get": &Get{},
}

var jobStates = map[string]bool{
    "ready":   true,
    "delayed": true,
    "buried":  true,
}

var (
    ErrInvalidSyntax   = errors.New("Invalid syntax")
    ErrInvalidCommand  = errors.New("Invalid command. Type help for a list of commands")
    ErrInvalidJobState = errors.New("Invalid job state. Should be one of ready/delayed/buried")
)

func ParseCmd(c string) (cmd Cmd, err error) {
    parts := strings.Split(c, " ")

    cmd, exists := validCmds[parts[0]]
    if !exists {
        return cmd, ErrInvalidCommand
    }

    if err := cmd.Validator(parts); err != nil {
        if err == ErrInvalidSyntax {
            return cmd, errors.New(err.Error() + ": Usage '" + cmd.Usage() + "'")
        }

        return cmd, err
    }

    //Add additional cmd args if specified
    if len(parts) > 1 {
        cmd.SetArgs(parts[1:])
        debugLog("Set cmd args to ", cmd.GetArgs())
    }

    return cmd, nil
}

func getNext(state string) (uint64, []byte, error) {
    switch state {
    case "ready":
        return cTubes.Conns[cTubes.SelectedIdx].PeekReady()
    case "delayed":
        return cTubes.Conns[cTubes.SelectedIdx].PeekDelayed()
    case "buried":
        return cTubes.Conns[cTubes.SelectedIdx].PeekBuried()
    }

    debugLog("Invalid state ", state)

    return 0, []byte(""), ErrInvalidJobState
}
