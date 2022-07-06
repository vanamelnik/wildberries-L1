package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// Задание 15
// К каким негативным последствиям может привести данный фрагмент кода, и как это исправить? Приведите корректный пример реализации.
// var justString string
// func someFunc() {
//   v := createHugeString(1 << 10)
//   justString = v[:100]
// }
// func main() {
//   someFunc()
// }

var justString string

func someFunc() {
	v := createHugeString(1 << 10)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v)).Data // получаем реальный указатель
	fmt.Println(hdr)

	justString = v[:100]
	// Когда осуществляется обрезка строки, по факту новая память не выделяется, а "обрезок" (даже если он сохранён в новую переменную) указывает на тот же
	// адрес в памяти, что и изначальная строка. Поскольку строки в Go неизменяемы, данным ничего не угрожает.
	// Но если мы присваиваем обрезанную строку глобальной переменной и выходим из функции, происходит утечка памяти, т.к. глобальная переменная указывает
	// на тот же адрес, что и локальная.
	hdr = (*reflect.SliceHeader)(unsafe.Pointer(&justString)).Data // видим, что указатели совпадают!
	fmt.Println(hdr)
}

// корректный вариант реализации
func someFuncCorrect1() {
	v := createHugeString(1 << 10)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v)).Data
	fmt.Println(hdr)

	// Чтобы избежать утечки памяти, склонируем строку. Функция Clone выполняет копирование урезаной строки в новую память,
	// и область памяти, занятая ненужной большой строкой, при выход из функции может быть освобождена сборщиком мусора.
	justString = strings.Clone(v[:100])
	hdr = (*reflect.SliceHeader)(unsafe.Pointer(&justString)).Data
	fmt.Println(hdr)
}

// Кроме того, когда мы "разрезаем" строку подобным образом, есть опасность, что мы "рассечём" руну (если в строке встречаются руны юникода по 2 или 3 байта).
// Например:
// 		s := "古池や　蛙飛び込む　水の音"
// 		fmt.Println(s[:8])
// Результат будет "古池��"
// Лучше резать строку, представляя её массивом рун. В этом случае новая строка копируется в новую память методами языка.
func someFuncCorrect2() {
	v := createHugeString(1 << 10)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v)).Data
	fmt.Println(hdr)

	justString = string([]rune(v[:100]))
	hdr = (*reflect.SliceHeader)(unsafe.Pointer(&justString)).Data
	fmt.Println(hdr)
}

func createHugeString(n int) string {
	s := strings.Repeat("a", n)
	return s
}

func main() {
	someFunc()
	fmt.Println()
	someFuncCorrect1()
	fmt.Println()
	someFuncCorrect2()
}
