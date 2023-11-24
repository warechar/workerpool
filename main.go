package main

import (
	"fmt"
	"time"
)

func main() {

	//q := deque.New[int]()
	//
	//ch := make(chan int)
	//
	//go func() {
	//
	//	for {
	//		fmt.Println(q)
	//		select {
	//		case ch <- q.Pop():
	//			//if q.Front() != 0 {
	//			//	q.Pop()
	//			//}
	//		}
	//	}
	//}()
	//
	//go func() {
	//Lo:
	//	for {
	//		select {
	//		case i := <-ch:
	//			if i == 0 {
	//				break Lo
	//			}
	//			fmt.Println("输出", i)
	//		default:
	//
	//		}
	//	}
	//}()
	//
	//for i := 1; i <= 10; i++ {
	//	q.Push(i)
	//}
	//
	//time.Sleep(5 * time.Second)
	//
	//os.Exit(1)

	pool := New(1)

	for i := 1; i <= 5; i++ {
		pool.Submit(i)
	}

	time.Sleep(4 * time.Second)
	//pool.StopWait()
	fmt.Printf("结束main")
}
