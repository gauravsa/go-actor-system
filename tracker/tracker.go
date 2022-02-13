package tracker

import (
	"github.com/ian-kent/go-log/log"
	"time"
)

const queue_size = 10e4

type TrackScope string
const (
	ActorSystem TrackScope = "actor-system"
	Task TrackScope = "task"
)

type TrackMetric string

const (
	Submitted TrackMetric = "submitted"
	Completed TrackMetric = "completed"
	Rejected  TrackMetric = "rejected"
	ActiveActor TrackMetric = "active-actors"
)

type Track struct {
	Scope TrackScope
	Metric TrackMetric
	Val int
}

type Tracker struct {
	close_sig chan bool
	trackername string
	sysname string
	tracker chan Track
	metrics map[TrackScope]map[TrackMetric]int
}

func (m *Tracker) collectMetric() {
	for track := range m.tracker {
		if m.metrics[track.Scope] == nil {
			m.metrics[track.Scope] = map[TrackMetric]int{}
		}
		m.metrics[track.Scope][track.Metric] += track.Val
	}
	m.printMetric()
	m.close_sig <- true
}

func (m *Tracker) Shutdown() {
	close(m.tracker)
	<- m.close_sig
}

func (m *Tracker) GetTrackerChan() chan Track {
	return m.tracker
}

func (m *Tracker) foreverPrintMetric() {
	for {
		time.Sleep(1 * time.Second)
		m.printMetric()
	}
}

func (m *Tracker) printMetric() {
	log.Debug("system: %s metrics :%+v", m.sysname, m.metrics)
}

func CreateTracker(trackername, sysname string) *Tracker {
	tracker := &Tracker{
		close_sig: make(chan bool),
		trackername: trackername,
		sysname: sysname,
		tracker: make(chan Track, queue_size),
		metrics: map[TrackScope]map[TrackMetric]int{},
	}
	go tracker.collectMetric()
	go tracker.foreverPrintMetric()
	return tracker
}

func CreateCounterTrack(scope TrackScope, metric TrackMetric) Track{
	return Track{
		Scope:  scope,
		Metric: metric,
		Val:    1,
	}
}

func CreateTrack(scope TrackScope, metric TrackMetric, delta int) Track{
	return Track{
		Scope:  scope,
		Metric: metric,
		Val:    delta,
	}
}