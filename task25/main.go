package main

// Задание 25
// Реализовать собственную функцию sleep.

import (
	"fmt"
	"time"
)

// sleepLoop запускает цикл, работающий до момента, когда текущее время совпадает с заданным
// работает, но грузит процессор!
func sleepLoop(interval time.Duration) {
	wakeUp := time.Now().Add(interval)
	for time.Now().Before(wakeUp) {
	}
}

// sleepTimeAfter использует функцию time.After. Процессор не грузит.
func sleepTimeAfter(interval time.Duration) {
	<-time.After(interval)
}

func main() {
	interval := time.Second * 3
	fmt.Printf("sleepTimeAfter: Sleeping for %v...", interval)
	sleepTimeAfter(interval)
	fmt.Println(" OK")

	fmt.Printf("sleepLoop: Sleeping for %v...", interval)
	sleepLoop(interval)
	fmt.Println(" OK")
}
