package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Reads tasks from the given tasks channel and executes them.
func taskHandler(wg *sync.WaitGroup, tasksChannel chan Task, errorsChannel chan int) {
	for {
		// Get the task for executing
		task, taskChanNotEmpty := <-tasksChannel

		// If tasks channel is empty, reduce waiting group counter and return from the function
		if !taskChanNotEmpty {
			wg.Done()
			return
		}

		// Execute the task
		err := task()
		// If task returned error
		if err != nil {
			// Read one item from the errors channel
			_, errChanNotEmpty := <-errorsChannel

			// If errors channel is empty, reduce waiting group counter and return from the function
			if !errChanNotEmpty {
				wg.Done()
				return
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// I decided to avoid mutexes for locking tasks slice and errors counter,
// by using buffered channels and checking if item can be read.
func Run(tasks []Task, workersCount, maxErrorsCount int) error {
	tasksCount := len(tasks)

	// Max errors count must be 0 or greater, because I use it to create a channel
	if maxErrorsCount < 0 {
		maxErrorsCount = 0
	}

	// Reduce workers count, if tasks count is less then workers count
	if tasksCount < workersCount {
		workersCount = tasksCount
	}

	// Put all tasks to a buffered channel
	tasksChannel := make(chan Task, tasksCount)
	for _, task := range tasks {
		tasksChannel <- task
	}
	// Close tasks channel beforehand, to avoid goroutine locking, if channel is empty
	close(tasksChannel)

	// Create channel with length = max errors count
	errorsChannel := make(chan int, maxErrorsCount)
	for i := 0; i < maxErrorsCount; i++ {
		errorsChannel <- i
	}
	// Close errors channel beforehand, to avoid goroutine locking, if channel is empty
	close(errorsChannel)

	var wg sync.WaitGroup

	// Launch n workers
	for i := 0; i < workersCount; i++ {
		go taskHandler(&wg, tasksChannel, errorsChannel)
		wg.Add(1)
	}

	// Wait until all goroutines are finished
	wg.Wait()

	// Check if errors count exceeded the limit
	_, errChanNotEmpty := <-errorsChannel

	// If errors channel is empty, return ErrErrorsLimitExceeded
	if !errChanNotEmpty {
		return ErrErrorsLimitExceeded
	}

	return nil
}
