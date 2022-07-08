package main

import (
	"fmt"
	"sort"
	"strings"
)

// Задание 12
// Имеется последовательность строк - (cat, cat, dog, cat, tree) создать для нее собственное множество.

// StringSet представляет собой реализацию множества строк в Го. Выгоднее всего хранить его в map, т.к. время доступа
// элементу структуры близко к константному.
type StringSet map[string]struct{}

// NewStringSet создаёт новое множество из переданного слайса.
func NewStringSet(v []string) StringSet {
	set := make(StringSet)
	for _, s := range v {
		set[s] = struct{}{}
	}
	return set
}

// Has возвращает true, если elem является членом множетсва.
func (s *StringSet) Has(elem string) bool {
	_, ok := (*s)[elem]
	return ok
}

// Add добавляет во множество новый элемент.
func (s *StringSet) Add(elem string) {
	(*s)[elem] = struct{}{}
}

// Delete удаляет элемент из множетсва.
func (s *StringSet) Delete(elem string) {
	delete(*s, elem)
}

// Clear очищает множество.
func (s *StringSet) Clear() {
	*s = make(StringSet)
}

// String реализует интерейс Stringer.
func (s *StringSet) String() string {
	elements := make([]string, 0, len(*s))
	for element := range *s {
		elements = append(elements, element)
	}
	// поскольку порядок в map неопределён, из соображений идемпотентности необходимо отсортировать список.
	sort.Strings(elements)
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func main() {
	v := []string{"cat", "cat", "dog", "cat", "tree"}
	set := NewStringSet(v)
	fmt.Println(set.String())

	set.Add("gopher")
	set.Add("gopher") // ничего не происходит - суслик уже есть!
	fmt.Println(set.String())

	fmt.Println(set.Has("tree"))
	set.Delete("tree")
	set.Delete("tree") // ничего не происходит
	fmt.Println(set.Has("tree"))

	fmt.Println(set.String())
	set.Clear()
	fmt.Println(set.String())
}
