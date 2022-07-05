package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Задание 7
// Реализовать конкурентную запись данных в map.

// Конкурентная запись в map возможна с использованием
//	- sync.Mutex
//	- sync.RWMutex
//	- sync.Map
//
// Преимущества RWMutex перед Mutex очевидны при большом количестве операций чтения, которые в
// первом случае не блокируют друг друга. В этом задании мы сравним RWMutex и sync.Map.

// threadSafeMapper обявляет интерфейс для тестирования разных подходов к созданию потокобезопасых map.
type threadSafeMapper interface {
	// Store сохраняет данные в map.
	Store(key, value any)
	// Load загружает значение из map.
	Load(key any) (any, bool)
	// Delete удаляет значение из map.
	Delete(key any)
}

var _ threadSafeMapper = (*rwMutexMap)(nil)

// rwMutexMap - объединение map с sync.RWMutex
type rwMutexMap struct {
	mu *sync.RWMutex
	m  map[any]any
}

// Store реализует интерфейс threadSafeMapper.
func (rwm *rwMutexMap) Store(key, value any) {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()
	rwm.m[key] = value
}

// Load реализует интерфейс threadSafeMapper.
func (rwm *rwMutexMap) Load(key any) (any, bool) {
	rwm.mu.RLock()
	defer rwm.mu.RUnlock()
	val, ok := rwm.m[key]
	return val, ok
}

// Delete реализует интерфейс threadSafeMapper.
func (rwm *rwMutexMap) Delete(key any) {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()
	delete(rwm.m, key)
}

// sync.Map уже сама по себе имплементирует наш интерфейс))
var _ threadSafeMapper = (*sync.Map)(nil)

// testOneStorerManyReaders тестирует одиночную запись и множественное чтение.
func testOneStorerManyReaders(sm threadSafeMapper) time.Duration {
	const sizeOfTestData = 1000000 // объем данных для считывания
	const numReadWorkers = 200     // количество воркеров-читателей
	const sizeOfMap = 1000000      // размер мапы
	const storerDelay = 1000       // задержка для более равномерного заполнения мапы

	ch := make(chan int, 100)
	wg := &sync.WaitGroup{}

	// storer - воркер, заполняющий мапу. Ключами являются последовательные интегеры,
	// значениями - случайные числа.
	storer := func() {
		for i := 0; i < sizeOfMap; i++ {
			val := rand.Int()
			sm.Store(i, val)
			// для того, чтобы запись производилась не только в начале работы
			// реализована задержка
			time.Sleep(time.Microsecond * time.Duration(rand.Intn(storerDelay)))
		}
	}

	// reader - воркер-читатель из мапы. Просто читает.
	reader := func() {
		for key := range ch {
			sm.Load(key)
		}
		wg.Done()
	}
	start := time.Now()
	// сначала запускаем заполняющий воркер
	go storer()
	// затем воркеры-читатели
	for i := 0; i < numReadWorkers; i++ {
		wg.Add(1)
		go reader()
	}
	// загружаем воркеров-читателей работой, посылая им случайные ключи
	for i := 0; i < sizeOfTestData; i++ {
		ch <- rand.Intn(sizeOfMap)
	}
	// закрытие канала - сигнал читателям закончить работу
	close(ch)
	// ждем завершения
	wg.Wait()
	return time.Since(start)
}

func main() {
	rwm := &rwMutexMap{
		mu: &sync.RWMutex{},
		m:  make(map[any]any),
	}
	sm := &sync.Map{}
	// Тесты показывают, что при заранее заполненной мапе выигрывает связка map+RWMutex.
	// Но если запись происходит и во время работы воркеров-читателей, то обе конструкции
	// выдают примерно одинаковый результат - sync.Map в этом случае работает почти в два раза
	// быстрее, нежели с заранее заполненной мапой.
	fmt.Print("Running testOneStorerManyReaders for rwMutexMap...")
	fmt.Printf(" done in %v\n", testOneStorerManyReaders(rwm))

	fmt.Print("Running testOneStorerManyReaders for sync.Map...")
	fmt.Printf(" done in %v\n", testOneStorerManyReaders(sm))
}
