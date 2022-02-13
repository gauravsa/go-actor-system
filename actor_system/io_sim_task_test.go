package actor_system

import (
	"go-actor-system/entities"
	"math/rand"
	"time"
)

// SimIOTask this is task to simulate IO operation.
// It would randomly take [0-10) milliseconds to finish
type SimIOTask struct {
	num int
}

func CreateNumberPrinterTask(i int) entities.Task{
	return &SimIOTask{
		num: i,
	}
}

func (t *SimIOTask) Execute() {
	x := 0
	if time.Now().Second() > 30 {
		x = 50
	}
	duration := time.Duration(x+rand.Intn(50)) * time.Millisecond
	<- time.After(duration)

}
