package main

import (
	"fmt"
	"sync"
)

// Задание 2
// Написать программу, которая конкурентно рассчитает значение квадратов чисел взятых
// из массива (2,4,6,8,10) и выведет их квадраты в stdout.

func main() {
	numbers := []int{2, 4, 6, 8, 10}
	// используем sync.WaitGroup, чтобы дождаться окончания работы воркеров.
	wg := new(sync.WaitGroup)
	for _, num := range numbers {
		wg.Add(1)
		// передаём число в качестве аргумента горутине, чтобы избежать data race.
		go func(n int) {
			fmt.Println(n * n)
			wg.Done()
		}(num)
	}
	// ждём окончания
	wg.Wait()
}
