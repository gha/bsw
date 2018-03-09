//Copyright (c) 2014, David Persson. All rights reserved.
//https://github.com/davidpersson/bsa
package main

import (
	"github.com/kr/beanstalk"
)

type Tubes struct {
	Names []string
	Conns []beanstalk.Tube
}

func (t *Tubes) UseAll() error {
	t.Reset()

	allTubes, err := conn.ListTubes()
	if err != nil {
		return err
	}

	for _, tube := range allTubes {
		t.Names = append(t.Names, tube)
		t.Conns = append(t.Conns, beanstalk.Tube{conn, tube})
	}

	return nil
}

func (t *Tubes) Reset() {
	t.Names = t.Names[:0]
	t.Conns = t.Conns[:0]
}
