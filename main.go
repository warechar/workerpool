package main

import "fmt"

func main() {
	wp := New(5)

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
