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
		tasksCh      = make(chan Task)
		quit         = make(chan struct{})
		mu           = sync.Mutex{}
		actualErrors = 0
	)

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-tasksCh:
					if !ok {
						return
					}
					err := task()
					if err != nil {
						mu.Lock()
						actualErrors++
						if actualErrors == m {
							close(quit)
						}
						mu.Unlock()
					}
				case <-quit:
					return
				}
			}
		}()
	}

	for _, task := range tasks {
		select {
		case <-quit:
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
