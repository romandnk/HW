package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		wg           = sync.WaitGroup{}
		tasksCh      = make(chan Task, len(tasks))
		errorCh      = make(chan struct{})
		quit         = make(chan struct{})
		actualErrors = 0
	)
	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case task := <-tasksCh:
					err := task()
					if err != nil {
						errorCh <- struct{}{}
					}
				case <-quit:
					return
				}
			}
		}()
	}
	for range errorCh {
		actualErrors++
		if actualErrors == m {
			close(quit)
			wg.Wait()
			return ErrErrorsLimitExceeded
		}
	}
	wg.Wait()
	return nil
}
