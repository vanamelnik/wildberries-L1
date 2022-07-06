package main

// Задание 18
// Реализовать структуру-счетчик, которая будет инкрементироваться в конкурентной среде.
// По завершению программа должна выводить итоговое значение счетчика.

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// syncCounter является thread-safe счетчиком.
type syncCounter struct {
	counter uint64
}

// Inc инкрементирует счётчик
func (c *syncCounter) Inc() {
	atomic.AddUint64(&c.counter, 1)
}

// Count возвращает текущее выражение счётчика.
func (c *syncCounter) Count() uint64 {
	return atomic.LoadUint64(&c.counter)
}

// Reset устанавливает счётчик в 0.
func (c *syncCounter) Reset() {
	atomic.StoreUint64(&c.counter, 0)
}

func main() {
	const numOfWorkers = 500
	const countsPerWorker = 1000000

	wg := &sync.WaitGroup{}
	c := &syncCounter{}

	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < countsPerWorker; i++ {
				c.Inc()
			}
		}()
	}

	wg.Wait()
	fmt.Println(c.Count()) // должен вывести значение numOfWorkers * countsperWorker
}
