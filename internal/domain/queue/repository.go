package queue

type QueueRepository interface {
	Enqueue(any) error
	Dequeue() (any, error)
}
