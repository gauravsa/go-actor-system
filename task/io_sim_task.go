package task

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
	duration := time.Duration(rand.Intn(100)) * time.Millisecond
	<- time.After(duration)

}
