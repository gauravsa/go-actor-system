package tracker

import (
	"github.com/ian-kent/go-log/log"
	"sync"
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
	sysname string
	tracker chan Track
	metrics map[TrackScope]map[TrackMetric]int
	mutex *sync.RWMutex
}

func (m *Tracker) collectMetric() {
	for track := range m.tracker {
		if m.metrics[track.Scope] == nil {
			m.metrics[track.Scope] = map[TrackMetric]int{}
		}
		m.mutex.Lock()
		m.metrics[track.Scope][track.Metric] += track.Val
		m.mutex.Unlock()
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
	m.mutex.RLock()
	log.Debug("system: %s metrics :%+v", m.sysname, m.metrics)
	m.mutex.RUnlock()
}

func CreateTracker(sysname string) *Tracker {
	tracker := &Tracker{
		close_sig: make(chan bool),
		sysname: sysname,
		tracker: make(chan Track, queue_size),
		metrics: map[TrackScope]map[TrackMetric]int{},
		mutex: &sync.RWMutex{},
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