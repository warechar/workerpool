package main

import (
	"reflect"
	"sync"
	"time"
	"workerpool/deque"
)

type AdaptorType interface{}

type HandlerFunc struct {
	F     func()
	Delay time.Time
	T     int
}

func (hf HandlerFunc) Get() any {
	return hf.Delay
}

func (hf HandlerFunc) Compare(handlerFunc HandlerFunc) bool {
	return hf.Delay.Sub(handlerFunc.Delay) < 0
}

type workerTimer struct {
	workerTimerQueue  chan HandlerFunc
	waitingDelayQueue *deque.DequeTimer[HandlerFunc]
	waitTimerChan     chan struct{}
	stopWaitChan      chan struct{}
}

type WorkerPool struct {
	taskQueue    chan AdaptorType
	workerQueue  chan AdaptorType
	maxWorkers   int
	waitingQueue *deque.Deque[AdaptorType]
	stopCh       chan struct{}

	workerTimer
}

func New(maxWorkers int) *WorkerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	pool := &WorkerPool{
		taskQueue:   make(chan AdaptorType),
		workerQueue: make(chan AdaptorType),
		stopCh:      make(chan struct{}),

		maxWorkers:   maxWorkers,
		waitingQueue: deque.New[AdaptorType](),

		workerTimer: workerTimer{
			workerTimerQueue:  make(chan HandlerFunc),
			waitingDelayQueue: deque.NewTimer[HandlerFunc](),
			waitTimerChan:     make(chan struct{}),
			stopWaitChan:      make(chan struct{}),
		},
	}

	go pool.distributor()
	go pool.distributorDelay()

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

	if t, ok := task.(HandlerFunc); ok {

		pool.waitingDelayQueue.Push(t, pool.waitingDelayQueue.GetHead(), pool.waitingDelayQueue.GetTail()-1)
		//pool.waitTimerChan <- struct{}{}
	} else {
		pool.taskQueue <- task
	}
}

func (pool *WorkerPool) distributorDelay() {
	defer func() {
		//close(pool.waitTimerChan)
		close(pool.stopWaitChan)
	}()

	var workerCount int
	var wg sync.WaitGroup

Loop:
	for {
		if pool.waitingDelayQueue.Len() == 0 {
			_, ok := <-pool.waitTimerChan
			if !ok {
				break Loop
			}
			continue
		}

		task := pool.waitingDelayQueue.Front()

		if isNil(task) {
			break Loop
		}

		select {
		case pool.workerTimerQueue <- pool.waitingDelayQueue.Pop():
		default:
			if workerCount < pool.maxWorkers {
				wg.Add(1)
				go pool.workerDelay(task, &wg)
				workerCount++
			} else {
				pool.waitingDelayQueue.Push(task, pool.waitingDelayQueue.GetHead(), pool.waitingDelayQueue.GetTail()-1)

			}
		}
	}

	for pool.waitingDelayQueue.Len() != 0 {
		pool.workerTimerQueue <- pool.waitingDelayQueue.Front()
		pool.waitingDelayQueue.Pop()
	}

	for workerCount > 0 {
		pool.workerTimerQueue <- HandlerFunc{}
		workerCount--
	}

	wg.Wait()
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
		close(pool.waitTimerChan)
	})
	<-pool.stopCh
	<-pool.stopWaitChan
}

func (pool *WorkerPool) waitingForQueue() bool {
	select {
	// There are new tasks, rest assured queue
	case task, ok := <-pool.taskQueue:
		if !ok {
			return false
		}

		pool.waitingQueue.Push(task)
	// cannot pop directly to prevent data loss due to workerQueue blocking after pop because workerQueue would block
	case pool.workerQueue <- pool.waitingQueue.Front():
		pool.waitingQueue.Pop()
	}

	return true
}

func (pool *WorkerPool) workerDelay(task HandlerFunc, wg *sync.WaitGroup) {
	for !isNil(task) {
		for {
			now := time.Now()
			if task.Delay.Before(now) {
				task.F()
				break
			}
		}

		task = <-pool.workerTimerQueue
	}

	wg.Done()
}

// worker work to func or etc...
func (pool *WorkerPool) worker(task AdaptorType, wg *sync.WaitGroup) {
	for !isNil(task) {
		if fc, ok := task.(func()); ok {
			fc()
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
