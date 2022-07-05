package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// Задание 1
// Дана структура Human (с произвольным набором полей и методов). Реализовать встраивание методов
// в структуре Action от родительской структуры Human (аналог наследования).

// Human представляет человека с его целями.
type Human struct {
	Name    string
	Age     int
	Address string

	mu *sync.RWMutex // для потокобезопасного использования map.
	// goals - цели реализованы в виде map для бысроты поиска и простоты удаления.
	goals map[string]struct{}
}

// AddGoal добавляет цель к списку целей человека.
func (h *Human) AddGoal(goal string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.goals[goal] = struct{}{}
	return
}

// GetGoals - получить все цели в виде слайса строк.
func (h *Human) GetGoals() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	result := make([]string, 0, len(h.goals))
	for goal := range h.goals {
		result = append(result, goal)
	}
	return result
}

// String реализует интерфейс Stringer.
func (h *Human) String() string {
	return fmt.Sprintf("%s %d лет от роду, живущий в %s", h.Name, h.Age, h.Address)
}

// Action репрезентует действие человека. Наследует методы структуры Human.
type Action struct {
	Human

	actionType string
	deadline   time.Time
}

// String реализует интерфейс Stringer. Использует методы "родительского объекта".
func (a Action) String() string {
	var strGoals strings.Builder
	goals := a.Human.GetGoals()
	for i, goal := range goals {
		if len(goals) > 1 {
			if i == len(goals)-1 {
				strGoals.WriteString(" и ")
			} else if i > 0 {
				strGoals.WriteString(", ")
			}
			strGoals.WriteString(goal)
		}
	}

	return fmt.Sprintf("%s, имеющий цели %s, должен %s до %v.",
		a.Human.String(), strGoals.String(), a.actionType, a.deadline.Format(time.RFC822))
}

func main() {
	deadline, err := time.Parse(time.RFC822, "01 Jan 23 04:00 UTC")
	if err != nil {
		log.Fatal(err)
	}
	// определяем человека
	polikarp := Human{
		Name:    "Поликарп Серафимович Гудбайло",
		Age:     93,
		Address: "г.Мухосранск, ул.Барака Обамы д.22 кв.13",
		mu:      &sync.RWMutex{},
		goals:   make(map[string]struct{}),
	}
	// добавляем человеку цель
	polikarp.AddGoal("стать губернатором области")

	// определяем действие
	action := Action{
		Human:      polikarp,
		actionType: "дойти до ручки",
		deadline:   deadline,
	}
	// метод AddGoal "наследуется" структурой Action,
	// поэтому цель можно задать и через реализацию этой структуры:
	action.AddGoal("достичь духовного просветления")
	// либо так:
	action.Human.AddGoal("накопить миллион йен")

	// выводим результат
	fmt.Println(action.String())

	// определим Поликарпу новое действие:
	action1 := Action{
		Human:      polikarp,
		actionType: "быть хорошим самураем",
	}
	// цели Поликарпа доступны и через это действие:
	goals := action1.GetGoals()

	// в "дочерней" структуре доступны не только экспортируемые методы, но и поля:
	fmt.Printf("Чтобы %s в %d года, %s должен %s.\n",
		action1.actionType, action.Age, action.Name, strings.Join(goals, ", "))
}
