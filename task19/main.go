package main

// Задание 19
// Разработать программу, которая переворачивает подаваемую на ход строку (например: «главрыба — абырвалг»). Символы могут быть unicode.

import "fmt"

// reverse переворачивает переданную строку.
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func main() {
	s := "\xF0\x9F\x98\x8B привет 世界 \xF0\x9F\x92\xA5"
	fmt.Printf("%q ----> %q\n", s, reverse(s))
}
