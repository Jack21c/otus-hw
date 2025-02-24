package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrTasksFinished       = errors.New("all tasks finished")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errChan := make(chan error)
	var allowed atomic.Bool
	allowed.Store(true)
	taskChan := toChan(tasks)

	wg := new(sync.WaitGroup)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			routine(taskChan, errChan, &allowed)
		}()
	}

	var err error
	countErrors := 0
	countFinished := 0
	for countFinished < n {
		err = <-errChan
		if err != nil {
			if errors.Is(err, ErrTasksFinished) {
				allowed.Store(false)
				countFinished++
				continue
			}
			if !allowed.Load() {
				continue
			}
			countErrors++
			if m > 0 && countErrors >= m { // If m <= 0 then errors ignored
				allowed.Store(false)
			}
		}
	}

	wg.Wait()
	if m > 0 && countErrors >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func toChan(s []Task) chan Task {
	ch := make(chan Task, len(s))
	for _, t := range s {
		ch <- t
	}
	close(ch)
	return ch
}

func routine(
	taskChan chan Task,
	errChan chan error,
	allowed *atomic.Bool,
) {
	localCounter := 0
	for allowed.Load() {
		localCounter++
		f, ok := <-taskChan
		if ok {
			errChan <- f()
		} else {
			break
		}
	}
	errChan <- ErrTasksFinished
}
