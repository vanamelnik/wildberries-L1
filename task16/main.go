package main

// Задание 16
// Реализовать быструю сортировку массива (quicksort) встроенными методами языка.

import (
	"fmt"
	"math/rand"
	"time"
)

// quickSortRecursive реализует рекурсивный алгоритм быстрой сортировки.
func quickSortRecursive(arr []int) {
	if len(arr) < 2 {
		return
	}
	// разбиваем массив на 2 подмассива
	p := partition(arr)
	if p > 1 {
		quickSortRecursive(arr[:p]) // сортируем левую часть...
	}
	if p < len(arr)-2 {
		quickSortRecursive(arr[p+1:]) // и правую.
	}
}

// quickSort реалихует алгоритм быстрой сортировки без рекурсии.
// Использует собственный стек для хранения проверяемых интервалов.
// В случае больших массивов должен быть более эффективным, чем рекурсивный алгоритм, т.к. последний расходует больше памяти на хранение
// адресов вызовов в стеке программы.
func quickSort(arr []int) {
	if len(arr) < 2 {
		return
	}
	type sortRange struct {
		left, right int
	}
	sortStack := []sortRange{{left: 0, right: len(arr)}} // стек интервалов, которые нужно отсортировать
	for len(sortStack) != 0 {
		sp := len(sortStack) - 1 // указатель на последний элемент стека
		// извлекаем из стека интервал
		sRange := sortStack[sp]
		sortStack = sortStack[:sp]

		// разбиваем массив на 2 подмассива
		p := partition(arr[sRange.left:sRange.right])
		// добавляем в стек интервалы
		if p > 1 {
			sortStack = append(sortStack, sortRange{left: sRange.left, right: sRange.left + p}) // сортируем левую часть...
		}
		if sRange.right-(sRange.left+p+1) > 1 {
			sortStack = append(sortStack, sortRange{left: sRange.left + p + 1, right: sRange.right}) // сортируем левую часть...
		}
	}
}

// partition - разбиение массива по Хоару. Возвращает индекс опорного элемента
// (любой элемент массива левее данного индекса меньше не больше любого элемента правее данного индекса).
func partition(arr []int) int {
	pivot := arr[len(arr)/2] // опорным элементом выбираем средний элемент массива
	left, right := 0, len(arr)-1
	// разбиение массива на две части: любой элемент левой части меньше опорного,
	// а любой элемент в правой части - больше опорного.
	for {
		for ; left < len(arr) && arr[left] < pivot; left++ { // промотать left до первого элемента >= опорному
		}
		for ; right >= 0 && arr[right] > pivot; right-- { // промотать right до первого элемента <= опорному
		}
		if left >= right {
			break
		}
		if arr[right] == arr[left] { // наткнулись на повтор опорного элемента?
			right-- // пропустим правый экземпляр в правую часть
			continue
		}
		arr[left], arr[right] = arr[right], arr[left] // меняем местами элементы
	}
	return left
}

func main() {
	const numOfTests = 1000
	const arrayLength = 10000
	rand.Seed(time.Now().UnixNano())
	fillFn := func() int {
		return rand.Intn(1000) - 500
	}
	// тестируем рекурсивную реализацию алгоритма
	fmt.Print("Testing quickSortRecursive... ")
	testSort(numOfTests, arrayLength, quickSortRecursive, fillFn)

	// тестируем реализацию алгоритма
	fmt.Print("Testing quickSort... ")
	testSort(numOfTests, arrayLength, quickSort, fillFn)
}

// testSort проверяет функцию sortFn заданным количеством тестов.
// Тестовые массивы заполняются при помощи функции fillFn.
func testSort(numOfTests, arrayLength int, sortFn func([]int), fillFn func() int) {
	for i := 0; i < numOfTests; i++ {
		// подготавливаем тестовый слайс
		arr := make([]int, arrayLength)
		for idx := range arr {
			arr[idx] = fillFn()
		}
		backupArr := make([]int, len(arr))
		copy(backupArr, arr)
		quickSort(arr)
		if !isSorted(arr) {
			fmt.Printf("FAIL\n%v ---> %v\n", backupArr, arr)
		}
	}
	fmt.Println("OK")
}

// isSorted возвращает true, если переданный массив отсортирован в порядке неубывания.
func isSorted(arr []int) bool {
	if len(arr) < 2 {
		return true
	}
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
}
