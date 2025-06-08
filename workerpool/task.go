package workerpool

// task represents a unit of work to be processed by a worker.
// It holds arbitrary data and a function to execute with that data.
type task struct {
	data     interface{}       // Data associated with the task
	execFunc func(interface{}) // Function to execute the task
}

// newTask creates a new task with the provided function and data.
func newTask(execFunc func(interface{}), data interface{}) task {
	return task{data: data, execFunc: execFunc}
}

// process executes the task's function using the associated data.
func (t *task) process() {
	t.execFunc(t.data)
}
