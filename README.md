# workerpool

参考
- https://github.com/gammazero/workerpool  

并且简单实现了基于FIFO队列根据时间大小按小排序，时间小的优先出列的workerDelay

### Example

```golang
    // 延迟
    wp := New(2)

	wp.Submit(HandlerFunc{
		F: func() {
			fmt.Println("echo mike")
		},
		Delay: time.Now().Add(3 * time.Second),
	})

	wp.Submit(HandlerFunc{
		F: func() {
			fmt.Println("echo nick")
		},
		Delay: time.Now().Add(30 * time.Second),
	})

	wp.Submit(HandlerFunc{
		F: func() {
			fmt.Println("echo james")
		},
		Delay: time.Now().Add(15 * time.Second),
	})

	wp.Submit(HandlerFunc{
		F: func() {
			fmt.Println("echo ware")
		},
		Delay: time.Now().Add(1 * time.Second),
	})

    wp.StopWait()
```

```golang
    wp := New(2)
	requests := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

	for _, r := range requests {
		r := r
		wp.Submit(func() {
			fmt.Println("Handling request:", r)
		})
	}

	wp.StopWait()

```