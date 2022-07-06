package main

// Задание 24
// Разработать программу нахождения расстояния между двумя точками, которые представлены
// в виде структуры Point с инкапсулированными параметрами x,y и конструктором.

import (
	"fmt"
	"math"
)

// Point - точка на плоскости.
type Point struct {
	x, y float64
}

// NewPoint создаёт новую точку.
func NewPoint(x, y float64) Point {
	return Point{x, y}
}

// Distance возвращает расстояние от данной точки до родительской.
func (p Point) Distance(t Point) float64 {
	square := func(n float64) float64 { return n * n }
	// Пифагоровы штаны...
	return math.Sqrt(square(t.x-p.x) + square(t.y-p.y))
}

func main() {
	p1 := NewPoint(-2, -1)
	p2 := NewPoint(1, 3)
	fmt.Printf("The distance between p1 and p2 is %.2f\n", p1.Distance(p2)) // для проверки считаем расстояние и в ту,
	fmt.Printf("The distance between p2 and p1 is %.2f\n", p2.Distance(p1)) // и в другую стороны...
}
