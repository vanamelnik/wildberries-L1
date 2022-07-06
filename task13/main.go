package main

// Задание 13
// Поменять местами два числа без создания временной переменной.

import "fmt"

func main() {
	{
		a := 12.7
		b := -0.5

		// читерство))
		a, b = b, a
		fmt.Println(a, b)

		// с помощью сложения и вычитания:
		a = a + b
		b = a - b
		a = a - b
		fmt.Println(a, b)

		// с помощью умножения и деления:
		if a != 0 && b != 0 {
			a = a * b
			b = a / b
			a = a / b
			fmt.Println(a, b)
		}
	}

	{
		a := 8
		b := -3

		// с помощью XOR:
		a = a ^ b
		b = b ^ a
		a = a ^ b
		fmt.Println(a, b)
	}
}
