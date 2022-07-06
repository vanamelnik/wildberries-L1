package main

// Задание 14
// Разработать программу, которая в рантайме способна определить тип переменной: int, string, bool, channel из переменной типа interface{}.

import "fmt"

// whatsThis выводит в stdout тип переданной переменной и её значение.
func whatsThis(x interface{}) {
	fmt.Printf("x (%T) = %+v\n", x, x)
}

// xOperation производит различные операции над переданной переменной в зависимости от её типа и возвращает результат.
func xOperation(x interface{}) interface{} {
	switch y := x.(type) {
	// интегер возводим в квадрат
	case int:
		return y * y
	// строку переворачиваем:
	case string:
		runes := []rune(y)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	// bool инвертируем
	case bool:
		return !y
	// в канал посылаем сообщение
	case chan string:
		go func() {
			y <- "Hello!"
			close(y)
		}()
		return nil
	// у слайса возвращаем сумму
	case []int:
		var sum int
		for _, n := range y {
			sum += n
		}
		return sum
	}
	return nil
}

func main() {
	xs := []interface{}{123456789, "Preved medved!", true, make(chan string), []int{1, -2, 3, -4, 5, -6, 7}}
	for _, x := range xs {
		whatsThis(x)
		if chX, ok := x.(chan string); ok {
			xOperation(x)
			fmt.Printf("xOperation(x) sended %q to the channel\n\n", <-chX)
			continue
		}
		fmt.Printf("xOperation(x) = %v\n\n", xOperation(x))
	}
}
