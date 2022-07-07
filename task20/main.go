package main

// Задание 20
// Разработать программу, которая переворачивает слова в строке.
// Пример: «snow dog sun — sun dog snow».

import (
	"fmt"
	"strings"
)

// reverseWords меняет порядок слов в строке на противоположный.
func reverseWords(s string) string {
	words := strings.Fields(s)
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 { // действуем с двух сторон))
		words[i], words[j] = words[j], words[i]
	}
	return strings.Join(words, " ")
}

func main() {
	s := "катацумури соро соро ноборэ фудзи но яма"
	fmt.Printf("%q =====> %q\n", s, reverseWords(s))
}
