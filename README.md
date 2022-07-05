# Wildberries задание L1

### TODO:
task06


### Устные вопросы:

1. Какой самый эффективный способ конкатенации строк?
   _Наиболее эффективно с памятью работает **strings.Builder**:_

   ```go
    var sb strings.Builder
    sb.WriteString("Wild")
    sb.WriteString("berries")
    fmt.Println(sb.String()) // "Wildberries"
   ```

2. Что такое интерфейсы, как они применяются в Go?

3. Чем отличаются RWMutex от Mutex?

4. Чем отличаются буферизированные и не буферизированные каналы?

5. Какой размер у структуры struct{}{}?
   _0 байт_

6. Есть ли в Go перегрузка методов или операторов?
   _Перегрузки операторов нет. Из-за этого возникают неудобства, например, с пакетом **math/big**. Явной перегрузки методов тоже нет, можно создавать структуры, наследующие другие структуры и создавать свои методы к ним:_

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

   _В рандомной_

8. В чем разница make и new?

9. Сколько существует способов задать переменную типа slice или map?

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
