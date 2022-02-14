package actor_system

import (
	"github.com/ian-kent/go-log/log"
	"go-actor-system/actor"
	"go-actor-system/entities"
	"go-actor-system/tracker"
	"sync"
)

type ActorSystem struct {
	name           string
	assigner 	   entities.Actor
	wg             *sync.WaitGroup
	tracker        *tracker.Tracker
}

func (system *ActorSystem) SubmitTask(task entities.Task) error {
	return system.assigner.AddTask(task)
}

func (system *ActorSystem) Run() {
	log.Debug("actor system %s started \n", system.name)
	go system.assigner.Start()
}

func (system *ActorSystem) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()
	system.assigner.Stop()
	system.wg.Wait()
	system.tracker.Shutdown()
	log.Debug("actor system: %s shutdown completed ", system.name)
}

// CreateActorSystem invokes actors and returns close_sig chan to close
func CreateActorSystem(name string, config *actor.Config) *ActorSystem{
	wg := &sync.WaitGroup{}
	systracker := tracker.CreateTracker(name)
	pool := actor.CreateTaskActorPool(wg)
	system := &ActorSystem{
		name:        name,
		wg:             wg,
		tracker:        systracker,
        assigner: actor.CreateAssignerActor(pool, systracker, config),
	}

	go system.Run()

	return system
}



