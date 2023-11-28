package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"
	"workerpool/deque"
)

type AdaptorType interface{}

type HandlerFunc struct {
	F     func()
	Delay time.Time
	T     int32
}

func (hf HandlerFunc) Get() int32 {
	return hf.T
}

func (hf HandlerFunc) Compare(T any) bool {
	return hf.Delay.Sub(reflect.ValueOf(T).Interface().(HandlerFunc).Delay) < 0
}

type WorkerPool struct {
	taskQueue         chan AdaptorType
	workerQueue       chan AdaptorType
	maxWorkers        int
	waiting           int32
	waitingQueue      *deque.Deque[AdaptorType]
	waitingDelayQueue *deque.Deque[AdaptorType]
	stopCh            chan struct{}
}

func New(maxWorkers int) *WorkerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	pool := &WorkerPool{
		taskQueue:         make(chan AdaptorType),
		workerQueue:       make(chan AdaptorType),
		stopCh:            make(chan struct{}),
		maxWorkers:        maxWorkers,
		waitingQueue:      deque.New[AdaptorType](),
		waitingDelayQueue: deque.New[AdaptorType](),
	}

	go pool.distributor()

	return pool
}

func (pool *WorkerPool) Size() int {
	return pool.maxWorkers
}

// Submit submit task to taskQueue
func (pool *WorkerPool) Submit(task AdaptorType) {
	if isNil(task) {
		return
	}

	pool.taskQueue <- task
}

func (pool *WorkerPool) distributor() {
	defer close(pool.stopCh)

	var workerCount int
	var wg sync.WaitGroup

Loop:
	for {
		// The waitingQueue is processed first
		if pool.waitingQueue.Len() != 0 {
			if !pool.waitingForQueue() {
				break Loop
			}

			continue
		}

		if pool.waitingDelayQueue.Len() != 0 {
			if !pool.waitingForQueue1() {
				break Loop
			}

			continue
		}

		select {
		case task, ok := <-pool.taskQueue:
			if !ok {
				break Loop
			}

			select {
			case pool.workerQueue <- task:
			default:
				if workerCount < pool.maxWorkers {
					wg.Add(1)
					go pool.worker(task, &wg)
					workerCount++
				} else {
					pool.waitingQueue.Push(task)
				}
			}
		}
	}

	for pool.waitingQueue.Len() != 0 {
		pool.workerQueue <- pool.waitingQueue.Front()
		pool.waitingQueue.Pop()
	}

	// dec workerCount when task is nil, worker exit
	for workerCount > 0 {
		pool.workerQueue <- nil
		workerCount--
	}

	wg.Wait()
}

func (pool *WorkerPool) StopWait() {
	once := sync.Once{}
	once.Do(func() {
		close(pool.taskQueue)
	})
	<-pool.stopCh
}

func (pool *WorkerPool) waitingForQueue() bool {
	select {
	// There are new tasks, rest assured queue
	case task, ok := <-pool.taskQueue:
		if !ok {
			return false
		}
		if t, ok := task.(HandlerFunc); ok {
			pool.waitingDelayQueue.Push(t)
		} else {
			pool.waitingQueue.Push(task)
		}
	// cannot pop directly to prevent data loss due to workerQueue blocking after pop because workerQueue would block
	case pool.workerQueue <- pool.waitingQueue.Front():
		pool.waitingQueue.Pop()
	}

	return true
}

func (pool *WorkerPool) waitingForQueue1() bool {
	select {
	// There are new tasks, rest assured queue
	case task, ok := <-pool.taskQueue:
		if !ok {
			return false
		}
		if t, ok := task.(HandlerFunc); ok {
			pool.waitingDelayQueue.Push(t)
		} else {
			pool.waitingQueue.Push(task)
		}
	// cannot pop directly to prevent data loss due to workerQueue blocking after pop because workerQueue would block
	case pool.workerQueue <- pool.waitingDelayQueue.Front():
		fmt.Println("?????????")
		pool.waitingDelayQueue.Pop()
	}

	return true
}

// worker work to func or etc...
func (pool *WorkerPool) worker(task AdaptorType, wg *sync.WaitGroup) {
	for !isNil(task) {
		if fc, ok := task.(func()); ok {
			fc()
		}

		if fc, ok := task.(HandlerFunc); ok {
			fc.F()
		}
		task = <-pool.workerQueue
	}

	wg.Done()
}

/**
Determine whether it is nil
*/
func isNil(task AdaptorType) bool {
	switch task.(type) {
	case HandlerFunc:
		if task == nil {
			return true
		}

		v := reflect.ValueOf(task)
		st := v.Interface().(HandlerFunc)
		if st.F == nil {
			return true
		}
	case int:
		if task == 0 {
			return true
		}
	case func():
		if task == nil {
			return true
		}
	case string:
		if task == "" {
			return true
		}
	default:
		return true
	}

	return false
}
