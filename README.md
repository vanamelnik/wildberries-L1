# Wildberries задание L1

### Устные вопросы:

1. Какой самый эффективный способ конкатенации строк?

   Наиболее эффективно с памятью работает **strings.Builder**:

   ```go
    var sb strings.Builder
    sb.WriteString("Wild")
    sb.WriteString("berries")
    fmt.Println(sb.String()) // "Wildberries"
   ```

2. Что такое интерфейсы, как они применяются в Go?

   Интерфейсы в Go представляют собой абстрактное описание поведения других типов. Интерфейс определяется через набор методов. Если у типа имеются все перечисленные в данном интерфейсе методы - значит он имплементирует этот интерфейс. Применение интерфейсов позволяет сделать код более гибким и наглядным. Если какая-либо функция в качестве аргумента принимает интерфейс, то ей можно передать любой тип, реализующий этот интерфейс (например, широко распространены интерфейсы Reader и Writer, из которых можно составлять целые цепочки). Также удобно использование интерфейсов для создания шаблонов в архитектуре приложения. Часто используется пустой интерфейс `interface{}`, вместо которого может быть подставлен любой тип.

3. Чем отличаются RWMutex от Mutex?

   RWMutex применяется со структурами, безопасными для конкурентного чтения, но небезопасными для записи. RWMutex позволяет не блокировать чтение, если в данный момент не идёт запись.

4. Чем отличаются буферизированные и не буферизированные каналы?

   При отправке сообщения в небуферизированный канал отправитель блокируется и ждёт, когда получатель извлечёт сообщение. В буферизированном канале отправитель блокируется, только если буфер заполнен полностью (и ждёт, пока освободится место для одного сообщения)

5. Какой размер у структуры struct{}{}?

   0 байт

6. Есть ли в Go перегрузка методов или операторов?

   Перегрузки операторов нет. Из-за этого возникают неудобства, например, с пакетом **math/big**. Явной перегрузки методов тоже нет, можно создавать структуры, наследующие другие структуры и создавать свои методы к ним:

   ```go
   type Cat struct {
   }

   func (c Cat) NumOfLegs() int {
       return 4
   }

   func (c Cat) Voice() string {
       return "Meow!"
   }

   type Kitten struct {
       Cat
   }

   func (k Kitten) Voice() string {
       return "Wheeee!"
   }

   func main() {
       k := Kitten{
           Cat: Cat{},
       }
       fmt.Println(k.NumOfLegs()) // 4
       fmt.Println(k.Voice())     // Wheeee!
   }
   ```

7. В какой последовательности будут выведены элементы map[int]int?

   Пример:

   ```go
   m[0]=1
   m[1]=124
   m[2]=281
   ```

   Обход map принципиально рандомный. Это реализовано на уровне языка.

8. В чем разница make и new?

   make создаёт только объекты слайс, map и канал и возвращает сам объект. new выделяет память и возвращает указатель на созданный объект. С помощью new невозможно создать map и chan.

9. Сколько существует способов задать переменную типа slice или map?

   map: 2 - явно(с хотя бы одним инициализированным элементом) и пустую через make
   slice: 4 (или 5) - явно, make, new, var (либо v := []int{}, что тоже можно считать явным заданием)

10. Что выведет данная программа и почему?

```go
func update(p *int) {
  b := 2
  p = &b
}

func main() {
  var (
     a = 1
     p = &a
  )
  fmt.Println(*p)
  update(p)
  fmt.Println(*p)
}
```

    Выведет 1 1. Update не работает, т.к. внутри функции меняется не значение по указателю, а значение самого указателя, которое остаётся в локальной видимости/ Функция должна либо вернуть новый указатель, либо она может изменить значение по указателю *p = 2

11. Что выведет данная программа и почему?

```go
func main() {
wg := sync.WaitGroup{}
for i := 0; i < 5; i++ {
wg.Add(1)
go func(wg sync.WaitGroup, i int) {
fmt.Println(i)
wg.Done()
}(wg, i)
}
wg.Wait()
fmt.Println("exit")
}
```

_Deadlock - т.к. в горутины передано скопированное значение WaitGroup. Необходимо передать указатель._

12. Что выведет данная программа и почему?

```go
func main() {
n := 0
if true {
n := 1
n++
}
fmt.Println(n)
}
```

    0, т.к. внутри блока if переменная n пересоздаётся оператором :=, и она не видна за пределами этого блока.

13. Что выведет данная программа и почему?

```go
func someAction(v []int8, b int8) {
  v[0] = 100
  v = append(v, b)
}

func main() {
  var a = []int8{1, 2, 3, 4, 5}
  someAction(a, 6)
  fmt.Println(a)
}
```

    100, 2, 3, 4, 5
    Функция someAction меняет нулевой элемент изначального слайса. Поскольку capacity изначального слайса недостаточно для добавления ещё одного элемента, функция append создаёт новый слайс в новом блоке памяти с cap = 2*cap, но время жизни его заканчивается вместе с функцией someAction.

14. Что выведет данная программа и почему?

```go
func main() {
  slice := []string{"a", "a"}

  func(slice []string) {
     slice = append(slice, "a")
     slice[0] = "b"
     slice[1] = "b"
     fmt.Print(slice)
  }(slice)
  fmt.Print(slice)
}
```

    ["b", "b", "a"]["a", "a"]
    Аналогично предыдущему примеру - функция работает с переданной переменной, поэтому всё, что происходит со слайсом после изменения размера, остаётся в локальной области видимости и уходит с завершением функции.
