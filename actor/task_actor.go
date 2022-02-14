package actor

import (
	"errors"
	"github.com/ian-kent/go-log/log"

	"go-actor-system/entities"
	"go-actor-system/tracker"
	"sync"
)

const taskQueueSize = 10

type TaskActor struct {
	id int
	closeSig chan bool
	wg *sync.WaitGroup
	tasks chan entities.Task
	tracker *tracker.Tracker
}

func (a *TaskActor) QueueSize() int {
	return len(a.tasks)
}

func (a *TaskActor) AddTask(task entities.Task) error {
	if len(a.tasks) >= taskQueueSize {
		return errors.New("filled queue")
	}
	a.tasks <- task
	return nil
}

func (a *TaskActor) Start() {
	defer a.wg.Done()
	a.wg.Add(1)
	log.Debug("starting actor :%d", a.id)
	for task := range a.tasks{
		task.Execute()
		a.tracker.GetTrackerChan() <- tracker.CreateCounterTrack(tracker.Task, tracker.Completed)
	}
	log.Debug("stopped actor :%d", a.id)
	a.closeSig <- true
}

func (a *TaskActor) Stop() {
	close (a.tasks)
	<- a.closeSig
}

func CreateActor(wg *sync.WaitGroup, id int, tracker *tracker.Tracker) entities.Actor {
	actor := &TaskActor{
		id: id,
		wg: wg,
		closeSig: make(chan bool),
		tasks: make(chan entities.Task, taskQueueSize),
		tracker: tracker,
	}
	go actor.Start()
	return actor
}


