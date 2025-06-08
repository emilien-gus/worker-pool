// Package workerpool provides a simple dynamically-scalable worker pool
// that processes tasks using a shared task queue.
// Each task wraps a user-defined function and input data.
//
// Example usage:
//
//	pool := NewWorkerPool(10, func(v interface{}) { fmt.Println(v) })
//	pool.AddWorker()
//	pool.Submit("hello")
//	pool.Stop()
package workerpool

import (
	"fmt"
	"sync"
)

// WorkerPool manages a pool of workers processing tasks concurrently
type WorkerPool struct {
	tasks    chan task         // Shared task queue
	workers  map[int]*worker   // Map of active workers
	nextID   int               // Next worker ID
	taskFunc func(interface{}) // Function to process each task
	wg       sync.WaitGroup    // WaitGroup to track worker completion
	mu       sync.Mutex        // Mutex to protect shared state
}

// NewWorkerPool creates and returns a new WorkerPool.
// Doesn't create workers. To process tasks you need to add worker.
func NewWorkerPool(bufferSize int, taskFunc func(interface{})) *WorkerPool {
	return &WorkerPool{
		tasks:    make(chan task, bufferSize), // Buffered channel for tasks
		taskFunc: taskFunc,                    // Task processing function
		nextID:   1,
		workers:  make(map[int]*worker),
	}
}

// AddWorker creates a new worker, starts it, and adds it to the pool
func (wp *WorkerPool) AddWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	worker := newWorker(wp.nextID, wp.tasks) // Create new worker with unique ID
	wp.workers[wp.nextID] = worker           // Add to map of workers
	wp.nextID++
	wp.wg.Add(1)

	// Start worker in a separate goroutine
	go func() {
		defer wp.wg.Done()
		worker.Start()
	}()
}

// RemoveWorker stops and removes one worker from the pool
func (wp *WorkerPool) RemoveWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if wp.nextID == 1 {
		return
	}
	// Remove last worker from the pool
	wp.nextID--
	wp.workers[wp.nextID].Stop()  // Signal worker to stop
	delete(wp.workers, wp.nextID) // Remove from map
}

// Submit adds a task to the pool's task queue
func (wp *WorkerPool) Submit(s interface{}) {
	wp.tasks <- newTask(wp.taskFunc, s) // Wrap input in task and send to channel
}

// Stop forces to stop all workers and waits for their completion
func (wp *WorkerPool) Stop() {
	close(wp.tasks) // Close task channel so workers can exit
	wp.wg.Wait()    // Wait for all workers to finish
	fmt.Println("All workers stopped")
}
