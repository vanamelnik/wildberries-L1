package main

// Задание 17
// Реализовать бинарный поиск встроенными методами языка.

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

// binarySearch возвращает индекс элемента val в массиве v.
// Если элемент отсутствует, функция возвращает -1.
// Для корректного результата массив должен быть отсортированным.
// Если в массиве есть повторяющиеся элементы, может быть возвращён индекс любого из них.
func binarySearch(v []int, val int) int {
	if !sort.IntsAreSorted(v) {
		panic("slice is not sorted")
	}
	offset := 0
	for len(v) != 0 {
		mid := len(v) / 2
		switch {
		case val == v[mid]:
			return offset + mid
		case val < v[mid]:
			v = v[:mid]
		case val > v[mid]:
			v = v[mid+1:]
			offset += mid + 1
		}
	}
	return -1
}

// тестирование функции binarySearch рандомными значениями
func main() {
	const sliceSize = 10000 // размер массива
	const maxStep = 50      // максимальный шаг между двумя элементами
	rand.Seed(time.Now().UnixNano())

	firstElement := rand.Intn(sliceSize/2) - rand.Intn(sliceSize) // первый элемент

	// randStep возвращает рандомный шаг
	randStep := func() int {
		if maxStep == 1 {
			return 1
		}
		return rand.Intn(maxStep-1) + 1 // шаг не должен быть равен 0
	}

	testSlice := make([]int, sliceSize) // тестовый слайс
	// notInSlice содержит все значения, не входящие в тестовый слайс, лежащие в его границах
	// (а также по одному значению, выходящему за границы тестового слайса).
	notInSlice := []int{firstElement - randStep()}

	// заполняем тестовые слайсы
	for i, k := 0, firstElement; i < sliceSize; i, k = i+1, k+randStep() {
		testSlice[i] = k
		if i == 0 {
			continue
		}
		for n := testSlice[i-1] + 1; n < k; n++ {
			notInSlice = append(notInSlice, n)
		}
	}
	notInSlice = append(notInSlice, testSlice[sliceSize-1]+randStep())
	errCount := 0

	// выполняем поиск каждого элемента, содержащегося в слайсе
	for _, val := range testSlice {
		idx := binarySearch(testSlice, val)
		if testSlice[idx] != val { // проверка: по возвращенному индексу должно быть искомое значение
			errCount++
			fmt.Printf("idx=%d, val=%d, testSlice[%d]=%d\n", idx, val, idx, testSlice[idx])
		}
	}

	// пытаемся найти в слайсе то, чего там нет
	for _, val := range notInSlice {
		if idx := binarySearch(testSlice, val); idx != -1 {
			fmt.Printf("val=%d, idx=%d (must be -1), testSlice[%d]=%d\n", val, idx, idx, testSlice[idx])
			errCount++
		}
	}
	if errCount > 0 {
		fmt.Printf("%d errors ocuried\n", errCount)
		os.Exit(1)
	}
	fmt.Println("OK")
}
