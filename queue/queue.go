package queue

type IQueue interface {
	Schedule() error
	Process()
	Close()
}

type Queue struct {
	Queue IQueue
}

func New(address string, client string) *Queue {
	queue := &Queue{}
	switch client {
	default:
		queue = NewAsynqClient(address)
	}
	return queue
}
