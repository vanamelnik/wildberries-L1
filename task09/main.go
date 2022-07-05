package main

// Задание 9
// Разработать конвейер чисел. Даны два канала: в первый пишутся числа (x) из массива,
// во второй — результат операции x*2, после чего данные из второго канала должны выводиться в stdout.

import (
	"fmt"
	"math"
	"math/rand"
)

// producer извлекает числа из переданного массива и посылает их в канал.
func producer(arr []int, ch chan<- int) {
	for _, x := range arr {
		ch <- x
	}
	close(ch) // закрытие канала - сигнал к окончанию работы следующего звена
}

// multiplier читает числа из канала chIn, умножает их на 2 и передаёт в канал chOut.
func multiplier(chIn <-chan int, chOut chan<- int) {
	for x := range chIn {
		chOut <- x * 2
	}
	close(chOut)
}

func main() {
	const arraySize = 1000
	arr := make([]int, arraySize)
	for i := range arr {
		arr[i] = rand.Intn(math.MaxInt / 2) // чтобы избежать переполнения при удвоении (и появления отрицательных чисел)
	}
	fmt.Println()
	prodToMult := make(chan int)
	multToOut := make(chan int)
	go producer(arr, prodToMult)
	go multiplier(prodToMult, multToOut)
	for x := range multToOut {
		fmt.Println(x)
	}
}
