package main

import (
	"fmt"
	"pool/workerpool"
	"time"
)

func main() {
	// Function to be executed for each task
	printTask := func(data interface{}) {
		msg, ok := data.(string)
		if ok {
			fmt.Println(msg)
		}
	}

	// Create a worker pool with a task buffer of size 10
	pool := workerpool.NewWorkerPool(10, printTask)

	// Add 3 workers to the pool
	pool.AddWorker()
	pool.AddWorker()
	pool.AddWorker()

	// Submit 10 tasks to the pool
	for i := 0; i < 10; i++ {
		time.Sleep(50 * time.Millisecond)
		pool.Submit(fmt.Sprintf("Task #%d", i))
	}

	// Wait for some tasks to be processed
	time.Sleep(1 * time.Second)

	// Remove one worker from the pool
	pool.RemoveWorker()

	// Submit 5 more tasks
	for i := 10; i < 15; i++ {
		time.Sleep(50 * time.Millisecond)
		pool.Submit(fmt.Sprintf("Task #%d", i))
	}

	// Give workers time to process the remaining tasks
	time.Sleep(1 * time.Second)

	// stop all workers
	pool.Stop()
}
