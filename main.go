package main

import (
	"flag"
	"fmt"

	"github.com/gizak/termui"
	"github.com/kr/beanstalk"
)

var (
	conn        *beanstalk.Conn
	cTubes      Tubes
	tubeTable   *termui.Table
	host        = flag.String("host", "127.0.0.1:11300", "Beanstalk host address")
	refreshRate = flag.Int("refresh", 1, "Refresh rate of the tube list (seconds)")
)

func main() {
	flag.Parse()

	err := termui.Init()
	if err != nil {
		panic(err)
	}

	defer termui.Close()

	registerEventHandlers()
	tubeTable = generateTable()

	conn, err = beanstalk.Dial("tcp", *host)
	if err != nil {
		panic(err)
	}

	updateTubes()

	termui.Loop()
}

func registerEventHandlers() {
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle(fmt.Sprintf("/timer/%ds", *refreshRate), func(e termui.Event) {
		updateTubes()
	})
}

func generateTable() *termui.Table {
	tubeTable = termui.NewTable()
	tubeTable.Rows = [][]string{}
	tubeTable.FgColor = termui.ColorWhite
	tubeTable.X = 0
	tubeTable.Y = 0
	tubeTable.Width = termui.TermWidth()
	tubeTable.Height = termui.TermHeight()
	tubeTable.Border = true
	tubeTable.Separator = false

	termui.Render(tubeTable)

	return tubeTable
}

func updateTubes() {
	if err := cTubes.UseAll(); err != nil {
		panic(err)
	}

	// Update tubes
	rows := [][]string{
		[]string{"Tube", "ready/delayed/buried", "waiting/watching/using"},
	}

	for _, tube := range cTubes.Conns {
		stats, err := tube.Stats()
		if err != nil {
			panic("Error loading stats for " + tube.Name + ": " + err.Error())
		}

		jobStats := stats["current-jobs-ready"] + " / " + stats["current-jobs-delayed"] + " / " + stats["current-jobs-buried"]
		workerStats := stats["current-waiting"] + " / " + stats["current-watching"] + " / " + stats["current-using"]

		row := []string{tube.Name, jobStats, workerStats}
		rows = append(rows, row)
	}

	tubeTable.Rows = rows
	termui.Render(tubeTable)
}
