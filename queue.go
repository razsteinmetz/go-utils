package razutils

import (
	"errors"
	"golang.org/x/exp/slices"
	"sync"
)

/*
a Simple thread safe FIFO Queue implementation
All access to the queue information is done under a Mutex locking, so its thread safe to be called from multiple
goroutines.
For some functions to work the items in the queue must be comparable.  However, with using the unique feature any
type can be added to the queue.

Not copyright and no warranty is made - use at your own discretion
*/

type Queue struct {
	data        []interface{} // it must provide the comparison to work
	length      int
	totalPushed int
	mu          sync.Mutex
}

// MakeQueue - create a new Queue with a starting capacity (it's a slice based, so its just allocating initial capacity).
func MakeQueue(initSize int) Queue {
	if initSize <= 0 {
		return Queue{}
	}
	return Queue{data: make([]interface{}, 0, initSize), length: 0, totalPushed: 0}
}

// TotalIn - return the total number of items added to the queue
func (q *Queue) TotalIn() int {
	q.mu.Lock()
	x := q.totalPushed
	q.mu.Unlock()
	return x
}

// Top - return the top (i.e. the oldest) item without removing it. Error is returned if the queue is empty
func (q *Queue) Top() (interface{}, error) {
	//log.Println("Q pop: before", q.data)
	q.mu.Lock()
	if q.length == 0 {
		q.mu.Unlock()
		return nil, errors.New("queue empty")
	}
	item := q.data[0]
	q.mu.Unlock()
	return item, nil

}

// Pop - return the top (i.e. the oldest) item while removing it. Error is returned if the queue is empty
func (q *Queue) Pop() (interface{}, error) {
	//log.Println("Q pop: before", q.data)
	q.mu.Lock()
	if q.length == 0 {
		q.mu.Unlock()
		return nil, errors.New("queue empty")
	}
	item := q.data[0]
	q.data = q.data[1:]
	q.length -= 1
	q.mu.Unlock()
	return item, nil
}

// Push - Push an item into the queue
func (q *Queue) Push(dt interface{}) {
	q.mu.Lock()
	q.data = append(q.data, dt)
	q.length += 1
	q.totalPushed += 1
	q.mu.Unlock()
}

// PushUnique - Push an item into the queue only if It's not already in it
func (q *Queue) PushUnique(dt interface{}) {
	if !q.InQueue(dt) {
		q.mu.Lock()
		q.data = append(q.data, dt)
		q.length += 1
		q.totalPushed += 1
		q.mu.Unlock()
	}
}

// PushMany - Push many items into the queue. If unique is true only new items will be pushed
func (q *Queue) PushMany(dt []interface{}, unique bool) {
	if len(dt) > 0 {
		if !unique {
			q.mu.Lock()
			q.data = append(q.data, dt...)
			q.length += len(dt)
			q.mu.Unlock()
		} else {
			for _, d := range dt {
				q.PushUnique(d)
			}
		}
	}
}

// Len - return the queue length
func (q *Queue) Len() int {
	q.mu.Lock()
	res := q.length
	q.mu.Unlock()
	return res
}

// IsEmpty - check if a queue is empty
func (q *Queue) IsEmpty() bool {
	q.mu.Lock()
	res := q.length == 0
	q.mu.Unlock()
	return res
}

// InQueue - check if an item is in the queue
func (q *Queue) InQueue(s interface{}) bool {
	q.mu.Lock()
	if q.Len() == 0 {
		return false
	}
	res := slices.IndexFunc(q.data, func(c interface{}) bool { return c == s }) != -1
	q.mu.Unlock()
	return res
}
