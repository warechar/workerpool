package main

import (
	"fmt"
)

func main() {

	//r := []int{1, 4, 7, 2, 8, 9, 3, 10, 2, 6}
	////r := []int{1, 4, 7, 2, 8, 9, 3, 10, 2}
	//seq := deque.NewSeq[HandlerFunc]()
	//for k, i := range r {
	//
	//	seq.Push(HandlerFunc{T: int32(i), Delay: time.Now().Add(time.Duration(i) * time.Second)}, 0, k-1)
	//	fmt.Println(seq)
	//	fmt.Println("===============")
	//}
	//
	//os.Exit(0)

	wp := New(1)

	requests := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	//
	//for i, r := range requests {
	//	r := r
	//	if i%2 == 0 {
	//		wp.Submit(HandlerFunc{
	//			F: func() {
	//				fmt.Println("i是", r)
	//			},
	//			Delay: time.NewTicker(1 * time.Second),
	//		})
	//	} else {
	//		wp.Submit(HandlerFunc{
	//			F: func() {
	//				fmt.Println("i是", r)
	//			},
	//			Delay: time.NewTicker(2 * time.Second),
	//		})
	//	}
	//}

	//
	for _, r := range requests {
		r := r
		wp.Submit(func() {
			fmt.Println("Handling request:", r)
		})
	}

	wp.StopWait()
}
