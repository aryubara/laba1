package main

import "fmt"

// Array представляет структуру динамического массива строк
type Array struct {
	data []string
}

// MCREATE создает новый массив с заданной начальной емкостью
func MCREATE(initCapacity int) Array {
	if initCapacity <= 0 {
		fmt.Println("Начальная вместимость должна быть > 0. По умолчанию - 4")
		initCapacity = 4
	}
	return Array{
		data: make([]string, 0, initCapacity),
	}
}

// MDESTROY очищает динамический массив
func MDESTROY(arr *Array) {
	arr.data = nil
}

// MPUSH_INDEX добавляет элемент по индексу с проверкой корректности
func MPUSH_INDEX(arr *Array, value string, index int) {
	for index < 0 || index > len(arr.data) {
		fmt.Printf("Некорректный индекс. Допустимо от 0 до %d\n", len(arr.data))
		index = ReadInt("Повторите ввод индекса: ")
	}

	// Делаем вставку в слайс эффективным Go-способом
	arr.data = append(arr.data, "")
	copy(arr.data[index+1:], arr.data[index:])
	arr.data[index] = value
}

// MPUSH добавляет элемент в конец массива
func MPUSH(arr *Array, value string) {
	arr.data = append(arr.data, value)
}

// MRECEIVE возвращает элемент по индексу
func MRECEIVE(arr *Array, index int) string {
	if len(arr.data) == 0 {
		fmt.Println("Массив пуст, действие невозможно")
		return ""
	}
	for index < 0 || index >= len(arr.data) {
		fmt.Printf("Некорректный индекс. Допустимо от 0 до %d\n", len(arr.data)-1)
		index = ReadInt("Повторите ввод индекса: ")
	}
	return arr.data[index]
}

// MDEL удаляет элемент по индексу и возвращает его значение
func MDEL(arr *Array, index int) string {
	if len(arr.data) == 0 {
		fmt.Println("Массив пуст, действие невозможно")
		return ""
	}
	for index < 0 || index >= len(arr.data) {
		fmt.Printf("Некорректный индекс. Допустимо от 0 до %d\n", len(arr.data)-1)
		index = ReadInt("Повторите ввод индекса: ")
	}

	deletedVal := arr.data[index]
	// Удаление из слайса со сдвигом элементов
	arr.data = append(arr.data[:index], arr.data[index+1:]...)
	return deletedVal
}

// MREPLACE заменяет элемент по индексу
func MREPLACE(arr *Array, value string, index int) {
	if len(arr.data) == 0 {
		fmt.Println("Массив пуст, действие невозможно")
		return
	}
	for index < 0 || index >= len(arr.data) {
		fmt.Printf("Некорректный индекс. Допустимо от 0 до %d\n", len(arr.data)-1)
		index = ReadInt("Повторите ввод индекса: ")
	}
	arr.data[index] = value
}

// MLEN возвращает текущую длину массива
func MLEN(arr Array) int {
	return len(arr.data)
}

// MGET выводит весь массив в консоль
func MGET(arr Array) {
	if len(arr.data) == 0 {
		fmt.Println("Массив пуст")
		return
	}
	for _, val := range arr.data {
		fmt.Print(val, " ")
	}
	fmt.Println()
}
