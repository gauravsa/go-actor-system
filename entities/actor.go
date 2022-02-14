package entities


type Actor interface {
	AddTask(task Task) error
	//QueueSize() int
	Start()
	Stop()
}
