package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadInt считывает целое число из консоли с валидацией
func ReadInt(prompt string) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка ввода. Попробуйте еще раз.")
			continue
		}
		input = strings.TrimSpace(input)
		value, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Ошибка: нужно целое число.")
			continue
		}
		return value
	}
}

// ReadString считывает строку из консоли
func ReadString(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
