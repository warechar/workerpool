package main

import (
	"fmt"
	"time"
)

func main() {
	wp := New(2)

	wp.Submit(HandlerFunc{
		F: func() {
			fmt.Println("hello")
		},
		Delay: time.Now().Add(3 * time.Second),
	})

	wp.Submit(HandlerFunc{
		F: func() {
			fmt.Println("hello1")
		},
		Delay: time.Now().Add(30 * time.Second),
	})

	wp.Submit(HandlerFunc{
		F: func() {
			fmt.Println("hello2")
		},
		Delay: time.Now().Add(15 * time.Second),
	})

	wp.Submit(HandlerFunc{
		F: func() {
			fmt.Println("hi")
		},
		Delay: time.Now().Add(1 * time.Second),
	})

	requests := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	for _, r := range requests {
		r := r
		wp.Submit(func() {
			fmt.Println("Handling request:", r)
		})
	}

	wp.StopWait()
}
