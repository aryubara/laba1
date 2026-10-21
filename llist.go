package main

import "fmt"

// LNode представляет узел двусвязного списка
type LNode struct {
	data string
	prev *LNode
	next *LNode
}

// LList представляет структуру двусвязного списка
type LList struct {
	head *LNode
	tail *LNode
	size int
}

// LCREATE создает пустой двусвязный список
func LCREATE() LList {
	return LList{head: nil, tail: nil, size: 0}
}

// LDESTROY очищает двусвязный список
func LDESTROY(list *LList) {
	list.head = nil
	list.tail = nil
	list.size = 0
}

// LPUSH_HEAD добавляет элемент в начало
func LPUSH_HEAD(list *LList, data string) {
	node := &LNode{data: data, prev: nil, next: list.head}
	if list.head != nil {
		list.head.prev = node
	} else {
		list.tail = node
	}
	list.head = node
	list.size++
	fmt.Println("Элемент успешно добавлен в начало")
}

// LPUSH_TAIL добавляет элемент в конец
func LPUSH_TAIL(list *LList, data string) {
	node := &LNode{data: data, prev: list.tail, next: nil}
	if list.tail != nil {
		list.tail.next = node
	} else {
		list.head = node
	}
	list.tail = node
	list.size++
	fmt.Println("Элемент успешно добавлен в конец")
}

// LPUSH_BEFORE добавляет элемент до указанного значения
func LPUSH_BEFORE(list *LList, data string, value string) {
	current := list.head
	for current != nil && current.data != value {
		current = current.next
	}
	if current == nil || current.data != value {
		fmt.Println("Ошибка: значение не найдено")
		return
	}
	node := &LNode{data: data, prev: current.prev, next: current}
	if current.prev != nil {
		current.prev.next = node
	} else {
		list.head = node
	}
	current.prev = node
	list.size++
	fmt.Printf("Элемент успешно добавлен до %s\n", value)
}

// LPUSH_AFTER добавляет элемент после указанного значения
func LPUSH_AFTER(list *LList, data string, value string) {
	current := list.head
	for current != nil && current.data != value {
		current = current.next
	}
	if current == nil || current.data != value {
		fmt.Println("Ошибка: значение не найдено")
		return
	}
	node := &LNode{data: data, prev: current, next: current.next}
	if current.next != nil {
		current.next.prev = node
	} else {
		list.tail = node
	}
	current.next = node
	list.size++
	fmt.Printf("Элемент успешно добавлен после %s\n", value)
}

// LDEL_HEAD удаляет элемент с начала списка
func LDEL_HEAD(list *LList) {
	if list.head == nil {
		fmt.Println("Список пуст")
		return
	}
	list.head = list.head.next
	if list.head != nil {
		list.head.prev = nil
	} else {
		list.tail = nil
	}
	list.size--
	fmt.Println("Элемент успешно удален")
}

// LDEL_TAIL удаляет элемент с конца списка
func LDEL_TAIL(list *LList) {
	if list.tail == nil {
		fmt.Println("Список пуст")
		return
	}
	list.tail = list.tail.prev
	if list.tail != nil {
		list.tail.next = nil
	} else {
		list.head = nil
	}
	list.size--
	fmt.Println("Элемент успешно удален")
}

// LDEL_BEFORE удаляет элемент перед указанным значением
func LDEL_BEFORE(list *LList, value string) {
	current := list.head
	for current != nil && current.data != value {
		current = current.next
	}
	if current == nil || current.prev == nil {
		fmt.Println("Ошибка: Нечего удалять!")
		return
	}
	deleted := current.prev
	current.prev = deleted.prev
	if deleted.prev != nil {
		deleted.prev.next = current
	} else {
		list.head = current
	}
	list.size--
	fmt.Println("Элемент успешно удален")
}

// LDEL_AFTER удаляет элемент после указанного значения
func LDEL_AFTER(list *LList, value string) {
	current := list.head
	for current != nil && current.data != value {
		current = current.next
	}
	if current == nil || current.next == nil {
		fmt.Println("Ошибка: Нечего удалять!")
		return
	}
	deleted := current.next
	current.next = deleted.next
	if deleted.next != nil {
		deleted.next.prev = current
	} else {
		list.tail = current
	}
	list.size--
	fmt.Println("Элемент успешно удален")
}

// LDEL_VALUE удаляет узел по точному значению
func LDEL_VALUE(list *LList, value string) {
	current := list.head
	for current != nil && current.data != value {
		current = current.next
	}
	if current == nil {
		fmt.Println("Ошибка: значение не найдено!")
		return
	}
	if current.prev != nil {
		current.prev.next = current.next
	} else {
		list.head = current.next
	}
	if current.next != nil {
		current.next.prev = current.prev
	} else {
		list.tail = current.prev
	}
	list.size--
	fmt.Println("Элемент успешно удален")
}

// LGET выводит двусвязный список от начала к концу
func LGET(list LList) {
	if list.head == nil {
		fmt.Println("Список пуст")
		return
	}
	current := list.head
	for current != nil {
		fmt.Print(current.data, " ")
		current = current.next
	}
	fmt.Println()
}

// LGET_BACK выводит двусвязный список в обратном направлении (с хвоста)
func LGET_BACK(list LList) {
	if list.head == nil {
		fmt.Println("Список пуст")
		return
	}
	current := list.tail
	for current != nil {
		fmt.Print(current.data, " ")
		current = current.prev
	}
	fmt.Println()
}

// LFIND ищет индекс элемента по его значению
func LFIND(list LList, value string) int {
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
