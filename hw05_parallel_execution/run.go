package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNotEnoughGoroutines = errors.New("number of goroutines must be more than 0")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		wg           = sync.WaitGroup{}
		tasksCh      = make(chan Task)
		actualErrors uint64
		quit         = make(chan struct{})
	)

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	if n <= 0 {
		return ErrNotEnoughGoroutines
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				err := task()
				if err != nil {
					atomic.AddUint64(&actualErrors, 1)
					if value := atomic.LoadUint64(&actualErrors); value == uint64(m) {
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
