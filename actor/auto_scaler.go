package actor

import (
	"github.com/ian-kent/go-log/log"
	"go-actor-system/entities"
	"go-actor-system/tracker"
	"time"
)

// auto scalar is part of task assigner actor
// It scales task actor pool based on queue len size
type autoScalar struct {
	*AssignerActor
	lastActorId int
	closingSig  chan bool
	closedSig   chan bool
}

func GetAutoScaler(assigner *AssignerActor, poolStarted chan bool) *autoScalar {
	scalar := &autoScalar{
		AssignerActor: assigner,
		lastActorId:   0,
		closingSig:    make(chan bool),
		closedSig:     make(chan bool),
	}
	go scalar.run(poolStarted)
	return scalar
}

func(scalar *autoScalar) run(poolStarted chan bool) {
	log.Debug("running auto scalar with min actor")
	scalar.provisionActors(scalar.Config.MinActor)
	poolStarted <- true
	completed := false
	for !completed {
		select {
		case <- scalar.closingSig:
			completed = true
		case <-time.After(100 * time.Millisecond):
			if scalar.QueueSize() > scalar.UpscaleQueueSize && len(scalar.pool) < scalar.MaxActor {
				scalar.provisionActors(1)

			} else if scalar.QueueSize() < scalar.DownscaleQueueSize && len(scalar.pool) > scalar.MinActor {
				scalar.deprovisionActors(1)
			}
		}
	}
	scalar.deprovisionActors(len(scalar.pool))
	log.Debug("scalar exited")
	scalar.closedSig <- true
}

func (scalar *autoScalar) stop() {
	scalar.closingSig <- true
	<- scalar.closedSig
}

func(scalar *autoScalar) deprovisionActors(delta int) {
	log.Debug("de-provision actors in %s by %d", scalar.name, delta)
	scalar.poolLock.Lock()
	if delta > len(scalar.pool) {
		delta = len(scalar.pool)
	}

	scalar.tracker.GetTrackerChan() <- tracker.CreateTrack(tracker.ActorSystem, tracker.ActiveActor, -delta)

	actors := scalar.pool[:delta]
	scalar.pool = scalar.pool[delta:]

	for _, actor := range actors {
		actor.Stop()
	}
	scalar.poolLock.Unlock()
}

func(scalar *autoScalar) provisionActors(delta int) []entities.Actor {
	log.Debug("provision actors in %s by %d", scalar.name, delta)

	actors := make([]entities.Actor, delta)
	for i := 0; i < delta; i += 1 {
		actors[i] = CreateActor(scalar.wg, i+scalar.lastActorId,scalar.tracker)
	}
	scalar.lastActorId += delta
	scalar.poolLock.Lock()
	scalar.tracker.GetTrackerChan() <- tracker.CreateTrack(tracker.ActorSystem, tracker.ActiveActor, delta )
	scalar.pool = append(scalar.pool, actors...)
	scalar.poolLock.Unlock()
	return actors
}
