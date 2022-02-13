package actor_system

import (
	"errors"
	"github.com/ian-kent/go-log/log"
	"go-actor-system/actor"
	"go-actor-system/config"
	"go-actor-system/entities"
	"go-actor-system/tracker"

	"sync"
	"time"
)

const queue_size = 10e2

type ActorSystem struct {
	close_sig      chan bool
	tasks          chan entities.Task
	actors         []*actor.Actor
	assigner_index int
	sysname        string
	wg             *sync.WaitGroup
	actoractions   *sync.Mutex
	tracker        *tracker.Tracker
	lastActorId    int
	*config.ActorSystemConfig
}

func (system *ActorSystem) SubmitTask(t entities.Task) error {
	if len(system.tasks) > queue_size {
		system.tracker.GetTrackerChan() <- tracker.CreateCounterTrack(tracker.Task, tracker.Rejected)
		return errors.New("queue size is full")
	}
	system.tasks <- t
	system.tracker.GetTrackerChan() <- tracker.CreateCounterTrack(tracker.Task, tracker.Submitted)
	return nil
}

func (system *ActorSystem) Run() {
	log.Debug("actor system %s started \n", system.sysname)
	go system.taskAssigner()
	go system.loadBalance()
}


func(system *ActorSystem) loadBalance() {
	for {
		if len(system.tasks) > system.UpscaleQueueSize && len(system.actors) < system.Maxactor {
			system.increaseActors(1)

		} else if len(system.tasks) < system.DownscaleQueueSize && len(system.actors) > system.Minactor{
			system.decreaseActors(1)
		}

		time.Sleep(10*time.Millisecond)
	}
}

func (system *ActorSystem) taskAssigner() {
	for task := range system.tasks {
		for {
			system.actoractions.Lock()
			system.assigner_index = system.assigner_index % len(system.actors)
			actor := system.actors[system.assigner_index]
			system.assigner_index += 1
			system.actoractions.Unlock()
			err := actor.AddTask(task)
			if err == nil {
				break
			}
		}
	}
	system.close_sig <- true
}

func(system *ActorSystem) increaseActors(delta int) {
	log.Debug("increasing actors in %s by %d", system.sysname, delta)
	system.actoractions.Lock()
	actors := createAndStartActors(system.lastActorId, delta, system.tracker, system.wg)
	system.lastActorId += delta
	system.tracker.GetTrackerChan() <- tracker.CreateTrack(tracker.ActorSystem, tracker.ActiveActor, delta )
	system.actors = append(system.actors, actors...)
	system.actoractions.Unlock()
}

func(system *ActorSystem) decreaseActors(delta int) {
	log.Debug("decreasing actors in %s by %d", system.sysname, delta)
	system.actoractions.Lock()
	if delta > len(system.actors) {
		delta = len(system.actors)
	}

	system.tracker.GetTrackerChan() <- tracker.CreateTrack(tracker.ActorSystem, tracker.ActiveActor, -delta)

	actors := system.actors[:delta]
	system.actors = system.actors[delta:]

	 for _, actor := range actors {
	 	actor.Stop()
	 }
	system.actoractions.Unlock()
}

func (system *ActorSystem) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()
	close(system.tasks)
	<- system.close_sig
	system.actoractions.Lock()
	for _, actor := range system.actors {
		actor.Stop()
	}
	system.tracker.GetTrackerChan() <- tracker.CreateTrack(tracker.ActorSystem, tracker.ActiveActor, -len(system.actors))
	system.actors = system.actors[len(system.actors)-1:]

	system.actoractions.Unlock()
	system.wg.Wait()
	system.tracker.Shutdown()
	log.Debug("actor system: %s shutdown completed ", system.sysname)
}

// CreateActorSystem invokes actors and returns close_sig chan to close
func CreateActorSystem(name string, config *config.ActorSystemConfig) *ActorSystem{
	wg := &sync.WaitGroup{}
	systracker := tracker.CreateTracker("systracker", name)
	actors := createAndStartActors(0, config.Minactor, systracker, wg)

	systracker.GetTrackerChan() <- tracker.CreateTrack(tracker.ActorSystem, tracker.ActiveActor, config.Minactor)
	system := &ActorSystem{
		close_sig:      make(chan bool),
		tasks:          make(chan entities.Task, 10e3),
		actors:         actors,
		assigner_index: 0,
		sysname:        name,
		wg:             wg,
		actoractions:   &sync.Mutex{},
		tracker:        systracker,
		lastActorId: len(actors),
		ActorSystemConfig: config,
	}

	go system.Run()

	return system
}

func createAndStartActors(start, numactor int, tracker *tracker.Tracker, wg *sync.WaitGroup) []*actor.Actor {
	actors := make([]*actor.Actor, numactor)

	for i := 0; i < numactor; i += 1 {
		actors[i] = actor.CreateActor(wg, i+start,tracker)
	}
	return actors
}


