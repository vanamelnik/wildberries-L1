package main

import "fmt"

// Задание 8
// Дана переменная int64. Разработать программу которая устанавливает i-й бит в 1 или 0.

// setBit устанавливает i-й бит в числе num в значение set(true = 1, false = 0)
func setBit(num uint64, i int, set bool) uint64 {
	if set {
		num = num | (1 << i) // если нужно установить - то OR
	} else {
		num = num ^ (1 << i) // если сбросить - то XOR
	}
	return num
}

func main() {
	var n uint64
	n = setBit(n, 10, true) // +1024
	fmt.Println(n)
	n = setBit(n, 2, true) // +4
	fmt.Println(n)
	n = setBit(n, 10, false) // -1024
	fmt.Println(n)
	n = setBit(n, 0, true) // +1
	fmt.Println(n)         // = 5
}
