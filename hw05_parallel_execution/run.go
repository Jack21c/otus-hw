package hw05parallelexecution

import (
	"errors"
	"fmt"
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
	fmt.Println("запустились")
	errChanMain := make(chan error)
	errsFromRoutines := make(map[int]chan error, n)
	var allowed atomic.Bool
	allowed.Store(true)
	taskChan := toChan(tasks)
	fmt.Println("записали ВСЕ задачи")

	wg := new(sync.WaitGroup)
	for i := 0; i < n; i++ {
		errsFromRoutines[i] = make(chan error)
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("создали очередную горутину")
			routine(taskChan, errsFromRoutines[i], &allowed)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		accumulateChans(errsFromRoutines, errChanMain, wg)
	}()

	var err error
	counter := 0
	for err = range errChanMain {
		if err != nil {
			if !allowed.Load() {
				continue
			}
			if errors.Is(err, ErrTasksFinished) {
				allowed.Store(false)
				continue
			}
			counter++
			fmt.Printf("Counter = %d\n", counter)
			if m > 0 && counter >= m { // If m <= 0 then errors ignored
				allowed.Store(false)
			}
		}
	}

	wg.Wait()
	if m > 0 && counter >= m {
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
			fmt.Printf("очередной запуск %d\n", localCounter)
			errChan <- f()
		} else {
			fmt.Println("задачи закончились")
			errChan <- ErrTasksFinished
			close(errChan)
			return
		}
	}
	fmt.Println("закончили горутину")
	close(errChan)
}

func accumulateChans(
	errsFromRoutines map[int]chan error,
	errChan chan error,
	wg *sync.WaitGroup,
) {
	stop := make(chan struct{})
	for i, ch := range errsFromRoutines {
		wg.Add(1)
		go func(to, from chan error) {
			defer wg.Done()
			for err := range from {
				to <- err
			}
			fmt.Printf("закончили горутину %d\n", i)
			stop <- struct{}{}
		}(errChan, ch)
	}
	for i := 0; i < len(errsFromRoutines); i++ {
		<-stop
	}
	fmt.Println("закончили горутину errThread")
	close(errChan)
}
