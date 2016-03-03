package main

import (
    "github.com/jroimartin/gocui"
    "strconv"
)

type Get struct{
    Args []string
}

func (c *Get) Validator(a []string) error {
    if len(a) < 2 || len(a) > 3 {
        return ErrInvalidSyntax
    }

    if a[1] == "next" {
        if len(a) != 3 || !jobStates[a[2]] {
            return ErrInvalidSyntax
        }

        return nil
    }

    if _, err := strconv.ParseUint(a[1], 10, 64); err != nil {
        debugLog("Invalid job ID: ", err.Error())
        return ErrInvalidSyntax
    }

    return nil
}

func (c *Get) Run(v *gocui.View) error {
    debugLog("Getting next ", c.Args[0], " job on tube ", cTubes.Selected)

    var jobID uint64
    var body []byte
    var err error

    switch c.Args[0] {
    case "next":
        debugLog("Getting next ", c.Args[0])
        jobID, body, err = getNext(c.Args[1])
    default:
        jobID, _ = strconv.ParseUint(c.Args[0], 10, 64)
        body, err = conn.Peek(jobID)
    }

    if err != nil {
        return err
    }

    debugLog("Got job ID ", jobID)
    PrintLine(v, strconv.FormatUint(jobID, 10))
    PrintLine(v, string(body))

    return nil
}

func (c *Get) SetArgs(a []string) {
    c.Args = a
}

func (c *Get) GetArgs() []string {
    return c.Args
}

func (c *Get) Usage() string {
    return "get <{job ID}/next> [<ready/delayed/buried>]"
}

func (c *Get) Description() string {
    return "Gets a job of the given state or ID"
}
