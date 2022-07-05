package main

// Задание 4
// Реализовать постоянную запись данных в канал (главный поток). Реализовать набор
// из N воркеров, которые читают произвольные данные из канала и выводят в stdout.
// Необходима возможность выбора количества воркеров при старте.
// Программа должна завершаться по нажатию Ctrl+C. Выбрать и обосновать способ завершения
// работы всех воркеров.

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// chWriter запускает воркеры и пишет случайные числа в канал, дожидаясь отмены контекста.
// Выполняется в главном потоке.
//
// Остановка воркеров реализована через закрытие канала. Это наилучший вариант т.к.
// не нужно создавать новых сущностей для достижения результата. В качестве альтернативы
// можно было бы передать общий контекст воркерам, но в этом случае трудно добиться завершения обработки
// всех сгенерированных данных.
func chWriter(ctx context.Context, n int) {
	ch := make(chan int, 5)
	wg := &sync.WaitGroup{}
	// запускаем воркеры
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go worker(i, ch, wg)
	}
	for {
		select {
		default:
			ch <- rand.Int()
		case <-ctx.Done(): // контекст завершён!
			close(ch) // закрытие канала служит сигналом остановки одновременно всем воркерам
			wg.Wait() // дожидаемся, пока воркеры завершат работу
			return
		}
	}
}

// worker читает данные из канала и выводит их в stdout.
func worker(i int, ch chan int, wg *sync.WaitGroup) {
	log.Printf("Worker %d started", i)
	for number := range ch { // закрытие канала остановит все запущенные воркеры
		time.Sleep(time.Millisecond * 500) // задержка для наглядности
		fmt.Println(number)
	}
	log.Printf("Worker %d stopped", i)
	wg.Done()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:\n\ttask04 <n> - where n is number of workers")
		os.Exit(1)
	}
	arg := os.Args[1]
	numOfWorkers, err := strconv.Atoi(arg)
	if err != nil || numOfWorkers < 1 {
		fmt.Printf("wrong number of workers: %s\n", arg)
		os.Exit(1)
	}
	// создаём контекст с отменой для плавного завершения
	ctx, cancel := context.WithCancel(context.Background())
	// подписываемся на сигнал остановки от ОС
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	// эта горутина ждёт сигнала от ОС и завершает контекст, давая тем самым сигнал остановить работу
	go func() {
		<-sigint
		log.Println("Shutting down...")
		cancel()
	}()
	chWriter(ctx, numOfWorkers)
	log.Println("Bye!")
	os.Exit(0)
}
