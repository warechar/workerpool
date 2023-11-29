package main

import (
	"fmt"
)

func main() {
	wp := New(1)

	requests := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	for _, r := range requests {
		r := r
		wp.Submit(func() {
			fmt.Println("Handling request:", r)
		})
	}

	wp.StopWait()
}
