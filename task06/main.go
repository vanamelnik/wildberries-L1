package main

// Задание 6
// Реализовать все возможные способы остановки выполнения горутины.

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Ниже описаны 7 способов завершения горутины.
//
//	1. Самостоятельно. Горутина завершается сама, выполнив возложенные на неё задачи. По окончании
//		она уведомляет о своём завершении, декрементируя пераданную WaitGroup (так же и в остальных примерах).
//
//	2. Контекст. Сигналом для выхода является истечение переданного контекста.
//
//	3. Закрытие канала. Сигналом к завершению работы всех горутин блока является закрытие переданного канала.
//		Можно также посылать в канал сообщения, закрывая воркеры по одиночке.
//
//	4. WaitGroup. Воркеры в блоке ждут, когда обнулится переданный sync.WaitGroup.
//
//	5. Сообщение в канале. Горутина завершается по "стоп-слову", переданному в канал.
//		Недостаток этого метода в том, что если несколько воркеров завязаны на этот канал, то количество переданных
//		сообщений со "стоп-словом" должно соответствовать количеству воркеров. Возможно, это будет полезно, когда
//		нужно закрыть часть однотипных воркеров.
//
//	6. Condition. Горутины ждут выполнения sync.Condition. Получив сигнал, они завершаются. Отправитель сигнала может
//		завершить часть воркеров, передав несколько сообщений cond.Signal, либо завершить всех слушателей, вызвав cond.Broadcast.
//
//	7. Отслеживание переменной в памяти. Горутина может следить за состоянием некоей переменной (в нашем случае сигналом
//		к завершению работы является достижение счётчика определённого порогового значения). Для безопасного отслеживания объектов
//		в памяти необходимо пользоваться синхронизирующими методами (мьютексами либо, как в приведённом примере atomic).
//		Использование этого способа противоречит принципу Р. Пайка "Don't communicate by sharing memory; share memory by communicating".
//
//
//	В приведённом  примере запускается по нескольку однотипных горутин, которые выполняют работу (увеличивают глобальный счётчик)
//	и ждут сигнала к завершению.

// numWorkersInGroup количество воркеров в блоке
const numWorkersInGroup = 5

// counter - глобальный счётчик.
var counter uint64

// selfStop - самостоятельное завершение работы горутины.
func selfStop(i int, wg *sync.WaitGroup) {
	log.Printf("\tselfStop%d started", i)

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
	atomic.AddUint64(&counter, uint64(rand.Intn(10)))

	log.Printf("\tselfStop%d stopped", i)
	wg.Done()
}

// contextStop - завершение, когда истечёт переданный контекст.
func contextStop(i int, ctx context.Context, wg *sync.WaitGroup) {
	log.Printf("\tcontextStop%d started", i)
	for {
		select {
		default:
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			atomic.AddUint64(&counter, uint64(rand.Intn(10)))
		case <-ctx.Done(): // конекст завершён
			log.Printf("\tcontextStop%d stopped", i)
			wg.Done()
			return
		}
	}
}

// chanCloseStop - завершение по закрытию канала.
func chanCloseStop(i int, stopCh <-chan struct{}, wg *sync.WaitGroup) {
	log.Printf("\tchanCloseStop%d started", i)
	for {
		select {
		default:
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
			atomic.AddUint64(&counter, uint64(rand.Intn(10)))
		case <-stopCh: // сообщение в канал либо закрытие канала
			log.Printf("\tchanCloseStop%d stopped", i)
			wg.Done()
			return
		}
	}
}

// waitgroupStop - завершение с помощью waitGroup.
func waitgroupStop(i int, wgStop, wg *sync.WaitGroup) {
	log.Printf("\twaitgroupStop%d started", i)

	atomic.AddUint64(&counter, uint64(rand.Intn(100)))
	wgStop.Wait()

	log.Printf("\twaitgroupStop%d stopped", i)
	wg.Done()
}

// msgChanStop - завершение по получении специального сообщения в канал (не подходит для групп, либо каждой нужен свой канал).
func msgChanStop(textCh <-chan string, wg *sync.WaitGroup) {
	log.Printf("\tmsgChanStop started")

	for msg := range textCh {
		log.Printf("\t\tmsgChanStop: received message: %q", msg)
		if msg == "Stop!!!" {
			log.Printf("\tmsgChanStop stopped")
			wg.Done()
			return
		}
		atomic.AddUint64(&counter, uint64(rand.Intn(100)))
	}
}

// condStop - завершение по общему сигналу в sync.Cond.
func condStop(i int, cond *sync.Cond, wg *sync.WaitGroup) {
	log.Printf("\tcondStop%d started", i)

	atomic.AddUint64(&counter, uint64(rand.Intn(150)))

	cond.L.Lock()
	defer cond.L.Unlock()
	cond.Wait()
	log.Printf("\tcondStop%d stopped", i)

	wg.Done()
}

// onCounterStop - данных в памяти (когда глобальный счётчик достигает определенного значения).
// Использование цикла с проверкой в каждом воркере - неоправданно дорогая операция. Лучше использовать один отслеживающий
// воркер, который будет слать сигнал в sync.Condition (как в предыдущем примере) либо в канал.
func onCounterStop(i int, wg *sync.WaitGroup) {
	log.Printf("\tonCounterStop%d started", i)
	for {
		if atomic.LoadUint64(&counter) > 4000 {
			log.Printf("\tonCounterStop%d stopped", i)
			wg.Done()
			return
		}
	}
}

func main() {
	wg := &sync.WaitGroup{} //общая WaitGroup
	cond := &sync.Cond{L: &sync.Mutex{}}
	textCh := make(chan string, 1)
	stopCh := make(chan struct{})
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	// запускаем группы воркеров
	log.Println("main: Children! Go for a walk!")

	for i := 1; i <= numWorkersInGroup; i++ {
		wg.Add(1)
		go selfStop(i, wg)
	}

	for i := 1; i <= numWorkersInGroup; i++ {
		wg.Add(1)
		go contextStop(i, ctxTimeout, wg)
	}

	for i := 1; i <= numWorkersInGroup; i++ {
		wg.Add(1)
		go chanCloseStop(i, stopCh, wg)
	}

	wgStop := &sync.WaitGroup{}
	wgStop.Add(1)
	for i := 1; i <= numWorkersInGroup; i++ {
		wg.Add(1)
		go waitgroupStop(i, wgStop, wg)
	}

	go msgChanStop(textCh, wg)

	for i := 1; i <= numWorkersInGroup; i++ {
		wg.Add(1)
		go condStop(i, cond, wg)
	}

	for i := 1; i <= numWorkersInGroup; i++ {
		wg.Add(1)
		go onCounterStop(i, wg)
	}

	textCh <- "Add to counter!" // сообщение воркеру msgChanStop

	log.Println("main: sleeping 4 seconds...")
	time.Sleep(time.Second * 4)
	log.Println("main: children! go home!")

	// воркеры selfStop, contextStop и onCounterStop завершатся сами

	// завершение воркеров chanCloseStop
	close(stopCh)
	// завершение waitgroupStop
	wgStop.Done()
	// завершение msgChanStop
	textCh <- "Stop!!!"
	// завершение condStop
	cond.L.Lock()
	cond.Broadcast()
	cond.L.Unlock()

	// ждем, когда воркеры "отпишутся" о завершении.
	log.Println("main: waiting for children...")
	wg.Wait()
	log.Printf("main: all goroutines are stopped; counter=%d", counter)
}
