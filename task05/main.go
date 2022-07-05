package main

// Задание 5
// Разработать программу, которая будет последовательно отправлять значения в канал,
// а с другой стороны канала — читать. По истечению N секунд программа должна завершаться.

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// writer отправляет последовательные числа в канал. Горутина завершается при истечении контекста.
func writer(ctx context.Context, ch chan<- int) {
	for i := 0; true; i++ { // выполнять до лучших времён
		select {
		default:
			ch <- i
			time.Sleep(time.Millisecond * 10) // задержка для наглядности
		case <-ctx.Done():
			close(ch) // закрытие канала является сигналом завершения для ридера
			return
		}
	}
}

// reader читает числа из канала до посинения (зачёркнуто) до его закрытия.
func reader(ch <-chan int, wg *sync.WaitGroup) {
	for num := range ch {
		fmt.Printf("%d\t", num) // сделаем красивые колонки табуляцией
	}
	fmt.Println()
	wg.Done() // отпускаем главный поток
}

func main() {
	// создаём контекст, который завершится через 5 сек.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()            // на всякий случай, чтобы избежать утечки контекста, и чтобы линтер не ругался
	wg := new(sync.WaitGroup) // последнее звено в цепи остановок - ридер, wg нужна для него
	ch := make(chan int, 1)
	wg.Add(1)
	go writer(ctx, ch)
	go reader(ch, wg)
	wg.Wait() // ждем остановки ридера
}
