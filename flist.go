package main

import "fmt"

// FNode представляет узел односвязного списка
type FNode struct {
	data string
	next *FNode
}

// FList представляет структуру односвязного списка
type FList struct {
	size int
	head *FNode
}

// FCREATE создает пустой список
func FCREATE() FList {
	return FList{head: nil, size: 0}
}

// FDESTROY очищает список (память очистит GC)
func FDESTROY(list *FList) {
	list.head = nil
	list.size = 0
}

// FPUSH_HEAD добавляет элемент в начало списка
func FPUSH_HEAD(list *FList, data string) {
	node := &FNode{data: data, next: list.head}
	list.head = node
	list.size++
	fmt.Println("Элемент успешно добавлен в начало")
}

// FPUSH_TAIL добавляет элемент в конец списка
func FPUSH_TAIL(list *FList, data string) {
	node := &FNode{data: data, next: nil}
	if list.head == nil {
		list.head = node
	} else {
		current := list.head
		for current.next != nil {
			current = current.next
		}
		current.next = node
	}
	list.size++
	fmt.Println("Элемент успешно добавлен в конец")
}

// FPUSH_BEFORE добавляет элемент до указанного значения
func FPUSH_BEFORE(list *FList, data string, value string) {
	if list.head == nil {
		fmt.Println("Список пуст")
		return
	}
	if list.head.data == value {
		FPUSH_HEAD(list, data)
		return
	}
	prev := list.head
	for prev.next != nil && prev.next.data != value {
		prev = prev.next
	}
	if prev.next == nil {
		fmt.Println("Узел не найден")
		return
	}
	node := &FNode{data: data, next: prev.next}
	prev.next = node
	list.size++
	fmt.Printf("Элемент успешно добавлен до значения %s\n", value)
}

// FPUSH_AFTER добавляет элемент после указанного значения
func FPUSH_AFTER(list *FList, data string, value string) {
	current := list.head
	for current != nil && current.data != value {
		current = current.next
	}
	if current == nil {
		fmt.Println("Узел не найден")
		return
	}
	node := &FNode{data: data, next: current.next}
	current.next = node
	list.size++
	fmt.Printf("Элемент успешно добавлен после %s\n", value)
}

// FDEL_HEAD удаляет элемент из начала списка
func FDEL_HEAD(list *FList) {
	if list.head == nil {
		fmt.Println("Список пуст")
		return
	}
	list.head = list.head.next
	list.size--
	fmt.Println("Элемент успешно удален")
}

// FDEL_TAIL удаляет элемент из конца списка
func FDEL_TAIL(list *FList) {
	if list.head == nil {
		fmt.Println("Список пуст")
		return
	}
	if list.head.next == nil {
		FDEL_HEAD(list)
		return
	}
	current := list.head
	for current.next.next != nil {
		current = current.next
	}
	current.next = nil
	list.size--
	fmt.Println("Элемент успешно удален")
}

// FDEL_BEFORE удаляет элемент перед указанным значением
func FDEL_BEFORE(list *FList, value string) {
	if list.head == nil {
		fmt.Println("Список пуст")
		return
	}
	if list.head.data == value {
		fmt.Println("Ошибка: удалить невозможно!")
		return
	}
	if list.head.next != nil && list.head.next.data == value {
		FDEL_HEAD(list)
		return
	}
	prev := list.head
	for prev.next != nil && prev.next.next != nil && prev.next.next.data != value {
		prev = prev.next
	}
	if prev.next == nil || prev.next.next == nil {
		fmt.Println("Ошибка: значение не найдено")
		return
	}
	prev.next = prev.next.next
	list.size--
	fmt.Println("Элемент успешно удален")
}

// FDEL_AFTER удаляет элемент после указанного значения
func FDEL_AFTER(list *FList, value string) {
	current := list.head
	for current != nil && current.data != value {
		current = current.next
	}
	if current == nil || current.next == nil {
		fmt.Println("Ошибка: значение не найдено")
		return
	}
	current.next = current.next.next
	list.size--
	fmt.Println("Элемент успешно удален")
}

// FDEL_VALUE удаляет конкретный элемент по значению
func FDEL_VALUE(list *FList, value string) {
	if list.head == nil {
		fmt.Println("Список пуст")
		return
	}
	if list.head.data == value {
		FDEL_HEAD(list)
		return
	}
	prev := list.head
	for prev.next != nil && prev.next.data != value {
		prev = prev.next
	}
	if prev.next == nil {
		fmt.Println("Ошибка: значение не найдено")
		return
	}
	prev.next = prev.next.next
	list.size--
	fmt.Println("Элемент успешно удален")
}

// FLEN возвращает длину списка
func FLEN(list FList) int {
	return list.size
}

// FGET_CLASSIC выводит список в обычном порядке
func FGET_CLASSIC(list FList) {
	if list.head == nil {
		fmt.Println("Список пуст")
		return
	}
	node := list.head
	for node != nil {
		fmt.Print(node.data, " ")
		node = node.next
	}
	fmt.Println()
}

// Рекурсивный вспомогательный метод для реверсивного вывода
func fgetReverseNode(node *FNode) {
	if node == nil {
		return
	}
	fgetReverseNode(node.next)
	fmt.Print(node.data, " ")
}

// FGET_REVERSE выводит список в обратном порядке рекурсивно
func FGET_REVERSE(list FList) {
	if list.head == nil {
		fmt.Println("Список пуст")
		return
	}
	fgetReverseNode(list.head)
	fmt.Println()
}

// FFIND ищет индекс элемента по значению
func FFIND(list FList, value string) int {
	index := 0
	current := list.head
	for current != nil {
		if current.data == value {
			return index
		}
		current = current.next
		index++
	}
	return -1
}
