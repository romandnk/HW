package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNotEnoughGoroutine  = errors.New("number of goroutines must be more than 0")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		wg           = sync.WaitGroup{}
		tasksCh      = make(chan Task)
		actualErrors int64
		quit         = make(chan struct{})
	)

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	if n <= 0 {
		return ErrNotEnoughGoroutine
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				err := task()
				if err != nil {
					atomic.AddInt64(&actualErrors, 1)
					if actualErrors == int64(m) {
						close(quit)
						return
					}
				}
			}
		}()
	}

	for _, task := range tasks {
		select {
		case <-quit:
			close(tasksCh)
			wg.Wait()
			return ErrErrorsLimitExceeded
		default:
			tasksCh <- task
		}
	}

	close(tasksCh)
	wg.Wait()

	return nil
}
