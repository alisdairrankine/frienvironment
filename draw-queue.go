package frienvironment

import "sync"

type DrawQueue struct {
	queue []func()

	sync.Mutex
}

func (q *DrawQueue) AddCommand(cmd func()) {
	q.Lock()
	defer q.Unlock()
	q.queue = append(q.queue, cmd)
}

func (q *DrawQueue) Run() {
	q.Lock()
	defer q.Unlock()
	for _, cmd := range q.queue {
		cmd()
	}
	// q.queue = []func(){}
}
