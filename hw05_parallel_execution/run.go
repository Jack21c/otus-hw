package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var counter atomic.Int32
	wg := new(sync.WaitGroup)
	taskChan := toChan(tasks, wg)

	if m <= 0 {
		m = len(tasks) + 1
	}

	for i := 0; i < n; i++ {
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

func toChan(s []Task, wg *sync.WaitGroup) chan Task {
	ch := make(chan Task, len(s))
	wg.Add(1)
	go func(ch chan Task, s []Task) {
		defer wg.Done()
		for _, t := range s {
			ch <- t
		}
		close(ch)
	}(ch, s)

	return ch
}

func routine(
	taskChan chan Task,
	counter *atomic.Int32,
	m int32,
) {
	for counter.Load() < m {
		f, ok := <-taskChan
		if ok {
			err := f()
			if err != nil {
				counter.Add(1)
			}
		} else {
			break
		}
	}
}
