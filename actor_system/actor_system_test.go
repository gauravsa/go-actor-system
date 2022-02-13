package actor_system

import (
	"github.com/ian-kent/go-log/log"
	"go-actor-system/config"
	"go-actor-system/task"
	"sync"
	"testing"
	"time"
)

func TestIOSimulationSystem(t *testing.T) {
	ioSimSystem := CreateActorSystem("io_sim", &config.ActorSystemConfig{
		Minactor: 10,
		Maxactor: 100,
		AutoScale: config.AutoScale{
			UpscaleQueueSize:   100,
			DownscaleQueueSize: 10,
		},
	})

	for i := 0; i < 1000; i += 1 {
		ioSimSystem.SubmitTask(task.CreateNumberPrinterTask(i))
		<-time.After(2 * time.Millisecond)
	}
	shutdown([]*ActorSystem{ioSimSystem})

}

func shutdown(systems []*ActorSystem) {

	wg := &sync.WaitGroup{}
	wg.Add(len(systems))
	for _, system := range systems {
		go system.Shutdown(wg)
	}
	wg.Wait()
	log.Debug("shutting down")

}
