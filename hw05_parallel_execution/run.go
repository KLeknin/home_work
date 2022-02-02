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
	var currentTask, errorsCount int32
	mu := new(sync.Mutex)
	if n > len(tasks) {
		n = len(tasks)
	}
	wg := sync.WaitGroup{}
	wg.Add(n)
	if m <= 0 {
		m = 1
	}

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				if currentTask >= int32(len(tasks)) || errorsCount >= int32(m) {
					mu.Unlock()
					return
				}
				task := tasks[currentTask]
				currentTask++
				mu.Unlock()
				err := task()
				if err != nil {
					mu.Lock()
					errorsCount++
					mu.Unlock()
					// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
					// Вопрос к преподавателю:
					// Почему если вместо предидущих трех строк использовать атомик:
					// atomic.AddInt32(&errorsCount, 1)
					// то при тестировании возникает "гонка" (go test -v -count=1 -race -timeout=1m .)
					// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				}
			}
		}()
	}

	wg.Wait()
	if int(atomic.LoadInt32(&errorsCount)) >= m {
		return (ErrErrorsLimitExceeded)
	}
	return nil
}
