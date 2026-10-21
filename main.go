package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Названия типов данных
const (
	TypeArray  = "ARRAY"
	TypeStack  = "STACK"
	TypeQueue  = "QUEUE"
	TypeTree   = "TREE"
	TypeFList  = "FLIST"
	TypeDList  = "DLIST"
	TypeUnknown = "UNKNOWN"
)

// Args хранит аргументы командной строки
type Args struct {
	file  string
	query string
}

// Query представляет распарсенную команду
type Query struct {
	command string
	name    string
	args    []string
}

// Entry представляет одну запись структуры данных в БД
type Entry struct {
	dataType string
	name     string
	values   []string
}

// DB отвечает за работу с файлом базы данных
type DB struct {
	filename string
	entries  []Entry
}

func ParseArgs(osArgs []string) (Args, bool) {
	var args Args
	for i := 1; i < len(osArgs); i++ {
		arg := osArgs[i]
		if arg == "--file" && i+1 < len(osArgs) {
			i++
			args.file = osArgs[i]
		} else if arg == "--query" && i+1 < len(osArgs) {
			i++
			args.query = osArgs[i]
		} else {
			fmt.Println("Неизвестный аргумент:", arg)
			return args, false
		}
	}
	if args.file == "" || args.query == "" {
		fmt.Println("Использование: go run . --file <file> --query '<cmd>'")
		return args, false
	}
	return args, true
}

func ParseQuery(str string) Query {
	fields := strings.Fields(str)
	if len(fields) < 2 {
		return Query{}
	}
	return Query{
		command: fields[0],
		name:    fields[1],
		args:    fields[2:],
	}
}

// Load считывает базу данных из файла
func (db *DB) Load() bool {
	file, err := os.Open(db.filename)
	if err != nil {
		// Если файла нет, это нормально (создадим новый)
		return true
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		headerFields := strings.Fields(line)
		if len(headerFields) < 3 {
			return false
		}

		entry := Entry{
			dataType: headerFields[0],
			name:     headerFields[1],
		}

		size, err := strconv.Atoi(headerFields[2])
		if err != nil {
			return false
		}

		for i := 0; i < size; i++ {
			if !scanner.Scan() {
				break
			}
			entry.values = append(entry.values, scanner.Text())
		}
		db.entries = append(db.entries, entry)
	}
	return true
}

// Save записывает текущее состояние БД в файл
func (db *DB) Save() bool {
	file, err := os.Create(db.filename)
	if err != nil {
		return false
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, entry := range db.entries {
		fmt.Fprintf(writer, "%s %s %d\n", entry.dataType, entry.name, len(entry.values))
		for _, value := range entry.values {
			fmt.Fprintln(writer, value)
		}
	}
	writer.Flush()
	return true
}

func (db *DB) GetOrCreate(dataType, name string) *Entry {
	for i := range db.entries {
		if db.entries[i].name == name {
			if db.entries[i].dataType != dataType {
				fmt.Println("Тип не совпадает: ожидается", db.entries[i].dataType)
				return nil
			}
			return &db.entries[i]
		}
	}
	db.entries = append(db.entries, Entry{dataType: dataType, name: name})
	return &db.entries[len(db.entries)-1]
}

func TypeByCommand(cmd string) string {
	if len(cmd) == 0 {
		return TypeUnknown
	}
	switch cmd[0] {
	case 'M':
		return TypeArray
	case 'S':
		return TypeStack
	case 'Q':
		return TypeQueue
	case 'T':
		return TypeTree
	case 'F':
		return TypeFList
	case 'L':
		return TypeDList
	default:
		return TypeUnknown
	}
}

// --- Функции конвертации (Маршалинг) ---

func ToArray(vec []string) Array {
	arr := MCREATE(len(vec))
	for _, str := range vec {
		MPUSH(&arr, str)
	}
	return arr
}

func FromArray(arr *Array, entry *Entry) {
	entry.values = make([]string, len(arr.data))
	copy(entry.values, arr.data)
	MDESTROY(arr)
}

func ToStack(vec []string) Stack {
	stack := SCREATE(4)
	for _, item := range vec {
		SPUSH(&stack, item)
	}
	return stack
}

func FromStack(stack *Stack, entry *Entry) {
	entry.values = make([]string, len(stack.data))
	copy(entry.values, stack.data)
	SDESTROY(stack)
}

func ToQueue(vec []string) Queue {
	queue := QCREATE(4)
	for _, item := range vec {
		QPUSH(&queue, item)
	}
	return queue
}

func FromQueue(queue *Queue, entry *Entry) {
	entry.values = make([]string, len(queue.data))
	copy(entry.values, queue.data)
	QDESTROY(queue)
}

func ToTree(vec []string) Tree {
	tree := TCREATE()
	for _, str := range vec {
		if num, err := strconv.Atoi(str); err == nil {
			TINSERT(&tree, num)
		} else if str == "null" {
			TINSERT(&tree, -999)
		}
	}
	return tree
}

func FromTree(tree *Tree, entry *Entry) {
	entry.values = nil
	if tree.Root == nil {
		TDESTROY(tree)
		return
	}
	q := []*TNode{tree.Root}
	for len(q) > 0 {
		curr := q[0]
		q = q[1:]
		entry.values = append(entry.values, strconv.Itoa(curr.Key))
		if curr.Left != nil {
			q = append(q, curr.Left)
		}
		if curr.Right != nil {
			q = append(q, curr.Right)
		}
	}
	TDESTROY(tree)
}

func ToFList(vec []string) FList {
	list := FCREATE()
	for _, str := range vec {
		FPUSH_TAIL(&list, str)
	}
	return list
}

func FromFList(list *FList, entry *Entry) {
	entry.values = nil
	for node := list.head; node != nil; node = node.next {
		entry.values = append(entry.values, node.data)
	}
	FDESTROY(list)
}

func ToDList(vec []string) LList {
	list := LCREATE()
	for _, str := range vec {
		LPUSH_TAIL(&list, str)
	}
	return list
}

func FromDList(list *LList, entry *Entry) {
	entry.values = nil
	for node := list.head; node != nil; node = node.next {
		entry.values = append(entry.values, node.data)
	}
	LDESTROY(list)
}

// --- Обработчик команд СУБД ---

type Command struct {
	argc    int
	handler func(entry *Entry, query Query) bool
}

func ExecuteCommand(entry *Entry, query Query) bool {
	cmdName := query.command

	switch cmdName {
	// ====== ARRAY ======
	case "MPUSH":
		arr := ToArray(entry.values)
		MPUSH(&arr, query.args[0])
		FromArray(&arr, entry)
		return true

	case "MPUSH_INDEX":
		idx, err := strconv.Atoi(query.args[1])
		if err != nil {
			fmt.Println("Не число:", query.args[1])
			return false
		}
		arr := ToArray(entry.values)
		MPUSH_INDEX(&arr, query.args[0], idx)
		FromArray(&arr, entry)
		return true

	case "MDEL":
		idx, err := strconv.Atoi(query.args[0])
		if err != nil {
			fmt.Println("Не число:", query.args[0])
			return false
		}
		arr := ToArray(entry.values)
		fmt.Println(MDEL(&arr, idx))
		FromArray(&arr, entry)
		return true

	case "MREPLACE":
		idx, err := strconv.Atoi(query.args[1])
		if err != nil {
			fmt.Println("Не число:", query.args[1])
			return false
		}
		arr := ToArray(entry.values)
		MREPLACE(&arr, query.args[0], idx)
		FromArray(&arr, entry)
		return true

	case "MRECEIVE":
		idx, err := strconv.Atoi(query.args[0])
		if err != nil {
			fmt.Println("Не число:", query.args[0])
			return false
		}
		arr := ToArray(entry.values)
		fmt.Println(MRECEIVE(&arr, idx))
		FromArray(&arr, entry)
		return true

	case "MLEN":
		arr := ToArray(entry.values)
		fmt.Println(MLEN(arr))
		FromArray(&arr, entry)
		return true

	case "MGET":
		arr := ToArray(entry.values)
		MGET(arr)
		FromArray(&arr, entry)
		return true

	// ====== STACK ======
	case "SPUSH":
		stack := ToStack(entry.values)
		SPUSH(&stack, query.args[0])
		fmt.Println(query.args[0])
		FromStack(&stack, entry)
		return true

	case "SPOP":
		stack := ToStack(entry.values)
		fmt.Println(SPOP(&stack))
		FromStack(&stack, entry)
		return true

	case "SGET":
		stack := ToStack(entry.values)
		SGET(stack)
		FromStack(&stack, entry)
		return true

	case "SGET_REVERSE":
		stack := ToStack(entry.values)
		SGET_REVERSE(stack)
		FromStack(&stack, entry)
		return true

	// ====== QUEUE ======
	case "QPUSH":
		queue := ToQueue(entry.values)
		QPUSH(&queue, query.args[0])
		fmt.Println(query.args[0])
		FromQueue(&queue, entry)
		return true

	case "QPOP":
		queue := ToQueue(entry.values)
		fmt.Println(QPOP(&queue))
		FromQueue(&queue, entry)
		return true

	case "QGET":
		queue := ToQueue(entry.values)
		QGET(queue)
		FromQueue(&queue, entry)
		return true

	case "QGET_BACK":
		queue := ToQueue(entry.values)
		QGET_BACK(queue)
		FromQueue(&queue, entry)
		return true

	// ====== TREE ======
	case "TINSERT":
		num, err := strconv.Atoi(query.args[0])
		if err != nil {
			fmt.Println("Не число:", query.args[0])
			return false
		}
		tree := ToTree(entry.values)
		if tree.Size > 0 && TFIND(tree, num) {
			fmt.Println("Значение уже есть в дереве:", num)
			return false
		}
		TINSERT(&tree, num)
		FromTree(&tree, entry)
		return true

	case "TDEL":
		num, err := strconv.Atoi(query.args[0])
		if err != nil {
			fmt.Println("Не число:", query.args[0])
			return false
		}
		tree := ToTree(entry.values)
		TDEL(&tree, num)
		FromTree(&tree, entry)
		return true

	case "TFIND":
		num, err := strconv.Atoi(query.args[0])
		if err != nil {
			fmt.Println("Не число:", query.args[0])
			return false
		}
		tree := ToTree(entry.values)
		if tree.Root == nil {
			fmt.Println("Дерево пустое")
			return true
		}
		if TFIND(tree, num) {
			fmt.Println("TRUE")
		} else {
			fmt.Println("FALSE")
		}
		FromTree(&tree, entry)
		return true

	case "TGET_PREORDER":
		tree := ToTree(entry.values)
		fmt.Printf("Дерево полное: %t\n", TIS_COMPLETE(tree))
		if tree.Size > 0 {
			TGET_PREORDER(tree.Root)
		} else {
			fmt.Print("Нет значений")
		}
		fmt.Println()
		FromTree(&tree, entry)
		return true

	case "TGET_INORDER":
		tree := ToTree(entry.values)
		fmt.Printf("Дерево полное: %t\n", TIS_COMPLETE(tree))
		if tree.Size > 0 {
			TGET_INORDER(tree.Root)
		} else {
			fmt.Print("Нет значений")
		}
		fmt.Println()
		FromTree(&tree, entry)
		return true

	case "TGET_POSTORDER":
		tree := ToTree(entry.values)
		fmt.Printf("Дерево полное: %t\n", TIS_COMPLETE(tree))
		if tree.Size > 0 {
			TGET_POSTORDER(tree.Root)
		} else {
			fmt.Print("Нет значений")
		}
		fmt.Println()
		FromTree(&tree, entry)
		return true

	case "TGET_LEVELORDER":
		tree := ToTree(entry.values)
		fmt.Printf("Дерево полное: %t\n", TIS_COMPLETE(tree))
		if tree.Size > 0 {
			TGET_LEVELORDER(tree)
		} else {
			fmt.Print("Нет значений")
		}
		fmt.Println()
    FromTree(&tree, entry)
return true
// ====== FLIST ======
case "FPUSH_HEAD":
list := ToFList(entry.values)
FPUSH_HEAD(&list, query.args[0])
FromFList(&list, entry)
return true
case "FPUSH_TAIL":
list := ToFList(entry.values)
FPUSH_TAIL(&list, query.args[0])
FromFList(&list, entry)
return true
case "FPUSH_BEFORE":
list := ToFList(entry.values)
FPUSH_BEFORE(&list, query.args[0], query.args[1])
FromFList(&list, entry)
return true
case "FPUSH_AFTER":
list := ToFList(entry.values)
FPUSH_AFTER(&list, query.args[0], query.args[1])
FromFList(&list, entry)
return true
case "FDEL_HEAD":
list := ToFList(entry.values)
FDEL_HEAD(&list)
FromFList(&list, entry)
return true
case "FDEL_TAIL":
list := ToFList(entry.values)
FDEL_TAIL(&list)
FromFList(&list, entry)
return true
case "FDEL_BEFORE":
list := ToFList(entry.values)
FDEL_BEFORE(&list, query.args[0])
FromFList(&list, entry)
return true
case "FDEL_AFTER":
list := ToFList(entry.values)
FDEL_AFTER(&list, query.args[0])
FromFList(&list, entry)
return true
case "FDEL_VALUE":
list := ToFList(entry.values)
FDEL_VALUE(&list, query.args[0])
FromFList(&list, entry)
return true
case "FFIND":
list := ToFList(entry.values)
idx := FFIND(list, query.args[0])
if idx < 0 {
fmt.Println("Нет такого значения")
} else {
fmt.Println(idx)
}
FromFList(&list, entry)
return true
case "FLEN":
list := ToFList(entry.values)
fmt.Println(FLEN(list))
FromFList(&list, entry)
return true
case "FGET":
list := ToFList(entry.values)
FGET_CLASSIC(list)
FromFList(&list, entry)
return true
case "FGET_REVERSE":
list := ToFList(entry.values)
FGET_REVERSE(list)
FromFList(&list, entry)
return true
// ====== DLIST ======
case "LPUSH_HEAD":
list := ToDList(entry.values)
LPUSH_HEAD(&list, query.args[0])
FromDList(&list, entry)
return true
case "LPUSH_TAIL":
list := ToDList(entry.values)
LPUSH_TAIL(&list, query.args[0])
FromDList(&list, entry)
return true
case "LPUSH_BEFORE":
list := ToDList(entry.values)
LPUSH_BEFORE(&list, query.args[0], query.args[1])
FromDList(&list, entry)
return true
case "LPUSH_AFTER":
list := ToDList(entry.values)
LPUSH_AFTER(&list, query.args[0], query.args[1])
FromDList(&list, entry)
return true
case "LDEL_HEAD":
list := ToDList(entry.values)
LDEL_HEAD(&list)
FromDList(&list, entry)
return true
case "LDEL_TAIL":
list := ToDList(entry.values)
LDEL_TAIL(&list)
FromDList(&list, entry)
return true
case "LDEL_BEFORE":
list := ToDList(entry.values)
LDEL_BEFORE(&list, query.args[0])
FromDList(&list, entry)
return true
case "LDEL_AFTER":
list := ToDList(entry.values)
LDEL_AFTER(&list, query.args[0])
FromDList(&list, entry)
return true
case "LDEL_VALUE":
list := ToDList(entry.values)
LDEL_VALUE(&list, query.args[0])
FromDList(&list, entry)
return true
case "LFIND":
list := ToDList(entry.values)
idx := LFIND(list, query.args[0])
if idx < 0 {
fmt.Println("Нет такого значения")
} else {
fmt.Println(idx)
}
FromDList(&list, entry)
return true
case "LGET":
list := ToDList(entry.values)
LGET(list)
FromDList(&list, entry)
return true
case "LGET_BACK":
list := ToDList(entry.values)
LGET_BACK(list)
FromDList(&list, entry)
return true
}
fmt.Println("Неизвестная команда:", cmdName)
return false
}
func main() {
args, ok := ParseArgs(os.Args)
if !ok {
os.Exit(1)
}
query := ParseQuery(args.query)
if query.command == "" || query.name == "" {
fmt.Println("Неверный запрос")
os.Exit(1)
}
dataType := TypeByCommand(query.command)
if dataType == TypeUnknown {
fmt.Println("Не удалось определить тип для команды:", query.command)
os.Exit(1)
}
db := DB{filename: args.file}
if !db.Load() {
fmt.Println("Ошибка загрузки")
os.Exit(1)
}
entry := db.GetOrCreate(dataType, query.name)
if entry == nil {
os.Exit(1)
}
if !ExecuteCommand(entry, query) {
os.Exit(1)
}
if !db.Save() {
fmt.Println("Ошибка сохранения")
os.Exit(1)
}
}

