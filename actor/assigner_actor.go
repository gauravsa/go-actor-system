package actor

import (
	"errors"
	"go-actor-system/entities"
	"go-actor-system/tracker"
)
const assignerQueueSize = 10e2

func CreateAssignerActor(pool *TaskActorPool, tracker *tracker.Tracker, config *Config) entities.Actor {
	return &AssignerActor{
		closeSig:      make(chan bool),
		tasks:         make(chan entities.Task, assignerQueueSize),
		assignerIndex: 0,
		TaskActorPool: pool,
		tracker:       tracker,
		Config:        config,
	}
}

type AssignerActor struct {
	name     string
	closeSig chan bool
	tasks    chan entities.Task
	assignerIndex int
	tracker *tracker.Tracker
	scalar  *autoScalar
	*TaskActorPool
	*Config
}

func (assigner *AssignerActor) QueueSize() int {
	return len(assigner.tasks)
}

func (assigner *AssignerActor) AddTask(task entities.Task) error{
	if len(assigner.tasks) >= assignerQueueSize {
		assigner.tracker.GetTrackerChan() <- tracker.CreateCounterTrack(tracker.Task, tracker.Rejected)
		return errors.New("task queue is full")
	}
	assigner.tasks <- task
	assigner.tracker.GetTrackerChan() <- tracker.CreateCounterTrack(tracker.Task, tracker.Submitted)
	return nil
}

func (assigner *AssignerActor) Start() {
	poolStarted := make(chan bool)
	assigner.scalar = GetAutoScaler(assigner, poolStarted)
	<- poolStarted
	for task := range assigner.tasks {
		for {
			assigner.poolLock.Lock()
			assigner.assignerIndex = assigner.assignerIndex % len(assigner.pool)
			actor := assigner.pool[assigner.assignerIndex]
			assigner.assignerIndex += 1
			assigner.poolLock.Unlock()
			err := actor.AddTask(task)
			if err == nil {
				break
			}
		}
	}
	assigner.closeSig <- true
}

func (assigner *AssignerActor) Stop() {
	close(assigner.tasks)
	<- assigner.closeSig
	assigner.scalar.stop()
}