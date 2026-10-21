package main

import "fmt"

// Stack представляет структуру стека строк
type Stack struct {
	data []string
}

// SCREATE создает новый стек с заданной начальной емкостью
func SCREATE(initCapacity int) Stack {
	if initCapacity <= 0 {
		fmt.Println("Начальная вместимость должна быть > 0. По умолчанию - 4")
		initCapacity = 4
	}
	return Stack{
		data: make([]string, 0, initCapacity),
	}
}

// SDESTROY очищает стек (в Go память очистит сборщик мусора)
func SDESTROY(stack *Stack) {
	stack.data = nil
}

// SPUSH добавляет элемент на вершину стека
func SPUSH(stack *Stack, value string) {
	stack.data = append(stack.data, value)
}

// SPOP удаляет и возвращает элемент с вершины стека
func SPOP(stack *Stack) string {
	if len(stack.data) == 0 {
		fmt.Println("Стек пуст")
		return ""
	}
	topIndex := len(stack.data) - 1
	value := stack.data[topIndex]
	stack.data = stack.data[:topIndex]
	return value
}

// SGET выводит стек от вершины к основанию
func SGET(stack Stack) {
	if len(stack.data) == 0 {
		fmt.Println("Стек пуст")
		return
	}
	for i := len(stack.data) - 1; i >= 0; i-- {
		fmt.Print(stack.data[i], " ")
	}
	fmt.Println()
}

// SGET_REVERSE выводит стек от основания к вершине
func SGET_REVERSE(stack Stack) {
	if len(stack.data) == 0 {
		fmt.Println("Стек пуст")
		return
	}
	for i := 0; i < len(stack.data); i++ {
		fmt.Print(stack.data[i], " ")
	}
	fmt.Println()
}
