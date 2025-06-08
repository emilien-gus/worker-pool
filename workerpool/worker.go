package workerpool

import "fmt"

// worker represents a single worker that processes tasks from a shared task channel.
type worker struct {
	id    int           // Unique ID of the worker
	tasks <-chan task   // Shared channel to receive tasks from the pool
	stop  chan struct{} // Channel used to signal the worker to stop
}

// newWorker creates a new worker with a given ID and shared task channel.
func newWorker(id int, tasks <-chan task) *worker {
	return &worker{
		id:    id,
		tasks: tasks,
		stop:  make(chan struct{}),
	}
}

// Start launches the worker loop that listens for tasks or a stop signal.
// If the task channel is closed, the worker exits gracefully.
func (w *worker) Start() {
	fmt.Printf("Worker %d started\n", w.id)
	for {
		select {
		case task, ok := <-w.tasks:
			if !ok {
				// The task channel was closed; stop the worker.
				fmt.Printf("Worker %d: tasks channel closed: worker stopped\n", w.id)
				return
			}
			// Process the received task.
			fmt.Printf("Worker %d processing task: ", w.id)
			task.process()
		case <-w.stop:
			// Received stop signal; exit loop.
			fmt.Printf("Worker %d stopped\n", w.id)
			return
		}
	}
}

// Stop sends a signal to the worker to stop.
func (w *worker) Stop() {
	close(w.stop)
}
