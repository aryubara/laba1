package main

import "fmt"

// Queue представляет структуру очереди строк
type Queue struct {
	data []string
}

// QCREATE создает новую очередь
func QCREATE(initCapacity int) Queue {
	if initCapacity <= 0 {
		fmt.Println("Начальная вместимость должна быть > 0. По умолчанию - 4")
		initCapacity = 4
	}
	return Queue{
		data: make([]string, 0, initCapacity),
	}
}

// QDESTROY очищает очередь
func QDESTROY(queue *Queue) {
	queue.data = nil
}

// QPUSH добавляет элемент в конец очереди
func QPUSH(queue *Queue, value string) {
	queue.data = append(queue.data, value)
}

// QPOP удаляет и возвращает элемент из начала очереди
func QPOP(queue *Queue) string {
	if len(queue.data) == 0 {
		fmt.Println("Ошибка: Очередь пуста")
		return ""
	}
	deleted := queue.data[0]
	queue.data = queue.data[1:] // Сдвигаем срез, удаляя первый элемент
	return deleted
}

// QGET выводит очередь от начала к концу
func QGET(queue Queue) {
	if len(queue.data) == 0 {
		fmt.Println("Очередь пуста")
		return
	}
	for _, val := range queue.data {
		fmt.Print(val, " ")
	}
	fmt.Println()
}

// QGET_BACK выводит очередь от конца к началу
func QGET_BACK(queue Queue) {
	if len(queue.data) == 0 {
		fmt.Println("Очередь пуста")
		return
	}
	for i := len(queue.data) - 1; i >= 0; i-- {
		fmt.Print(queue.data[i], " ")
	}
	fmt.Println()
}
