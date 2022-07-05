package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Задание 3
// Дана последовательность чисел: 2,4,6,8,10. Найти сумму их квадратов(22+32+42….)
// с использованием конкурентных вычислений.

// useAtomic использует возможности пакета sync.Atomic.
func useAtomic(numbers []int) int {
	var result int64
	// используем sync.WaitGroup, чтобы дождаться окончания работы воркеров.
	wg := new(sync.WaitGroup)
	for _, num := range numbers {
		wg.Add(1)
		// передаём число в качестве аргумента горутине, чтобы избежать data race.
		go func(n int) {
			// atomic гарантирует, что не будет гонки данных
			atomic.AddInt64(&result, int64(n*n))
			wg.Done()
		}(num)
	}
	// ждём окончания
	wg.Wait()
	return int(result)
}

// useChannels использует передачу данных и результата по каналам.
// Оставляем возможность задать размер буфера канала для экспериментов.
func useChannels(numbers []int, chSize int) int {
	wg := new(sync.WaitGroup)

	chSqr := make(chan int, chSize) // канал для передачи квадратов
	chResult := make(chan int)      // канал для результата

	for _, num := range numbers {
		wg.Add(1)
		go func(n int) {
			// передаем квадрат в канал
			chSqr <- n * n
			wg.Done()
		}(num)
	}
	// эта горутина собирает квадраты из канала и складывает их
	go func() {
		var result int
		// цикл завершится, когда канал будет закрыт
		for num := range chSqr {
			result += num
		}
		chResult <- result
	}()
	// ждём окончания
	wg.Wait()
	// закрытие канала даёт сигнал горутине послать результат
	close(chSqr)
	return <-chResult
}

func main() {
	// выполняем задание
	nums := []int{2, 4, 6, 8, 10}
	fmt.Println(useAtomic(nums), useChannels(nums, 5))

	// сравним скорость работы двух функций
	//
	// Результаты эксперимента показывают, что буфер канала ожидаемо становится узким местом,
	// т.к. горутинам приходится ждать возможности записать результат в канал.
	// Если размер буфера канала равен числу обрабатываемых элементов, то разницы в скорости
	// не обнаруживается. Но буфер канала забирает ощутимо больше памяти, поэтому подход с использованием
	// пакета atomic является более эффективным в данном случае.
	const sliceSize = 5000000
	const chanSize = 5000000
	numbers := make([]int, 0, sliceSize)
	for i := 1; i <= sliceSize; i++ {
		numbers = append(numbers, i*2)
	}
	start := time.Now()
	atomicRes := useAtomic(numbers)
	t1 := time.Since(start)
	start = time.Now()
	chanRes := useChannels(numbers, chanSize)
	t2 := time.Since(start)
	fmt.Printf("useAtomic:\tresult: %d, time: %v\n", atomicRes, t1)
	fmt.Printf("useChannels:\tresult: %d, time: %v\n", chanRes, t2)
}
