package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
//
//nolint:gosec
func Run(tasks []Task, n, m int) error {
	var counter atomic.Int32
	wg := new(sync.WaitGroup)
	if m <= 0 {
		m = len(tasks) + 1
	}
	taskChan := toChan(tasks, &counter, int32(m), wg)

	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			routine(taskChan, &counter, int32(m))
		}()
	}

	wg.Wait()
	if counter.Load() >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func toChan(
	s []Task,
	counter *atomic.Int32,
	m int32,
	wg *sync.WaitGroup,
) chan Task {
	ch := make(chan Task)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, t := range s {
			if counter.Load() >= m {
				break
			}
			ch <- t
		}
		close(ch)
	}()

	return ch
}

func routine(
	taskChan chan Task,
	counter *atomic.Int32,
	m int32,
) {
	for f := range taskChan {
		if counter.Load() >= m {
			break
		}
		err := f()
		if err != nil {
			counter.Add(1)
		}
	}
}
