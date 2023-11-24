package main

import (
	"fmt"
	"strconv"
	"sync"
	"workerpool/deque"
)

type WorkerPool struct {
	taskQueue    chan int
	workerQueue  chan int
	maxWorkers   int
	waiting      int32
	waitingQueue *deque.Deque[int]
	stopCh       chan struct{}
}

func New(maxWorkers int) *WorkerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	pool := &WorkerPool{
		taskQueue:    make(chan int),
		workerQueue:  make(chan int),
		stopCh:       make(chan struct{}),
		maxWorkers:   maxWorkers,
		waitingQueue: deque.New[int](),
	}

	go pool.distributor()

	return pool
}

func (pool *WorkerPool) Size() int {
	return pool.maxWorkers
}

// Submit submit task to taskQueue
func (pool *WorkerPool) Submit(task int) {
	if task == 0 {
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
			fmt.Println(pool.waitingQueue, "pool.waitingQueue.Len()"+strconv.Itoa(pool.waitingQueue.Len()))
			if !pool.waitingForQueue() {
				break Loop
			}

			continue
		}

		select {
		case task, ok := <-pool.taskQueue:
			if !ok {
				fmt.Println("没有东西运行")
				break Loop
			}

			select {
			case pool.workerQueue <- task:
			default:
				if workerCount < pool.maxWorkers {
					wg.Add(1)
					go pool.worker(task, pool.workerQueue, &wg)
					workerCount++
				} else {
					fmt.Println(task, "插入一个task到队列中")
					pool.waitingQueue.Push(task)
				}
			}
		}
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

		fmt.Println("重新插入到wait", task)
		pool.waitingQueue.Push(task)
		//case pool.workerQueue <- pool.waitingQueue.Front():
	case pool.workerQueue <- pool.waitingQueue.Pop():

		//fmt.Println("waitingForQueue推送一个pop过去" + strconv.Itoa(pool.waitingQueue.Pop()))
		fmt.Println("waitingForQueue推送一个pop过去")
	}

	//atomic.StoreInt32(&pool.waiting, int32(pool.waitingQueue.Len()))
	return true
}

func (pool *WorkerPool) worker(task int, workerQueue chan int, wg *sync.WaitGroup) {
	for task != 0 {
		fmt.Println("这里是worker输出" + strconv.Itoa(task))
		task = <-workerQueue
	}

	wg.Done()
}
