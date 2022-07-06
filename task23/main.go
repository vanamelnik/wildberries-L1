package main

// Задание 23
// Удалить i-ый элемент из слайса.

import "fmt"

// используем дженерики
type slice[T any] []T

func (s *slice[T]) deleteElement(i int) {
	if i < 0 || i > len(*s) {
		panic("index out of range")
	}
	// отрезаем до удаляемого элемента и добавляем то, что после
	*s = append((*s)[:i], (*s)[i+1:]...)
}

// благодаря дженерикам наш метод работает со слайсами любого типа
func main() {
	// интегеры:
	n := slice[int]{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	n.deleteElement(5)
	fmt.Println(n)

	// строки:
	s := slice[string]{"zero", "one", "two", "three", "four", "five"}
	s.deleteElement(3)
	fmt.Println(s)
}
