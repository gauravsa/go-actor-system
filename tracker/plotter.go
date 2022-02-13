package tracker

import (
	"github.com/ian-kent/go-log/log"
	"os"
	"strconv"
)

var printmetrics []TrackMetric = []TrackMetric{ActiveActor, Submitted, Completed}

type Line struct {
	Time int64
	Metrics map[TrackMetric]int

}

type Plotter struct {
	lines []Line
}

func GetPlotter() *Plotter {
	return &Plotter{lines: []Line{}}
}

func (p *Plotter) AddMetric(line Line) {
	p.lines = append(p.lines, line)
}

func (p *Plotter) Save() {
	fi, err := os.Create("../plot.csv")
	if err != nil {
		log.Error("unable to write to file: %s", err)
	}
	defer fi.Close()


	fi.Write([]byte("Time, " + ActiveActor+ ", " + Submitted + ", " +Completed +"\n"))

	for _, line := range p.lines {
		s := strconv.Itoa(int(line.Time)) + ", "
		for _, metric := range printmetrics {
			s += strconv.Itoa(line.Metrics[metric]) + ", "
		}
		s = s[:len(s)-2]
		s += "\n"
		fi.Write([]byte(s))
	}

}