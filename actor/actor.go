package actor

import (
	"errors"
	"github.com/ian-kent/go-log/log"

	"go-actor-system/entities"
	"go-actor-system/tracker"
	"sync"
)

const queue_len = 10

type Actor struct {
	id int
	wg *sync.WaitGroup
	tasks chan entities.Task
	tracker *tracker.Tracker
}

func (a *Actor) AddTask(task entities.Task) error {
	if len(a.tasks) >= queue_len {
		return errors.New("filled queue")
	}

	a.tasks <- task
	return nil
}

func (a *Actor) start() {
	defer a.wg.Done()
	a.wg.Add(1)
	log.Debug("starting actor :%d", a.id)
	for task := range a.tasks{
		task.Execute()
		a.tracker.GetTrackerChan() <- tracker.CreateCounterTrack(tracker.Task, tracker.Completed)
	}
	log.Debug("stopped actor :%d", a.id)
}

func (a *Actor) Stop() {
	//log.Debug("stopping actor :%d", a.id)
	close (a.tasks)
}

func CreateActor(wg *sync.WaitGroup, id int, tracker *tracker.Tracker) *Actor {
	actor := &Actor{
		id: id,
		wg: wg,
		tasks: make(chan entities.Task, queue_len),
		tracker: tracker,
	}
	go actor.start()
	return actor
}


