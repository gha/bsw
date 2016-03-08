package main

import (
    "github.com/jroimartin/gocui"
)

type Help struct{
    Args []string
}

func (c *Help) Validator(a []string) error {
    if len(a) != 1 {
        return ErrInvalidSyntax
    }

    return nil
}

func (c *Help) Run(v *gocui.View) error {
    debugLog("Showing command list")

    PrintLine(v, "")

    for _, cmd := range validCmds {
        line := cmd.Usage() + " - " + cmd.Description()
        PrintLine(v, line)
    }

    return nil
}

func (c *Help) SetArgs(a []string) {
    c.Args = a
}

func (c *Help) GetArgs() []string {
    return c.Args
}

func (c *Help) Usage() string {
    return "help"
}

func (c *Help) Description() string {
    return "Displays a list of commands"
}
