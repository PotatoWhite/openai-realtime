package utils

import "sync"

// RunGoroutine go Routine 실행하는 함수
func RunGoroutine(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fn()
	}()
}
