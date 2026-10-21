package main

import "fmt"

// TNode представляет узел дерева
type TNode struct {
	Key   int
	Left  *TNode
	Right *TNode
}

// Tree представляет структуру-обертку для дерева
type Tree struct {
	Root *TNode
	Size int
}

// TCREATE создает пустое дерево
func TCREATE() Tree {
	return Tree{Root: nil, Size: 0}
}

// TDESTROY уничтожает дерево (в Go достаточно обнулить корень, остальное соберет GC)
func TDESTROY(tree *Tree) {
	tree.Root = nil
	tree.Size = 0
}

// TINSERT добавляет элемент по уровням (сверху-вниз, слева-направо)
func TINSERT(tree *Tree, value int) {
	newNode := &TNode{Key: value, Left: nil, Right: nil}

	if tree.Root == nil {
		tree.Root = newNode
		tree.Size = 1
		return
	}

	// Срез выполняет роль очереди (queue)
	q := []*TNode{tree.Root}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr.Left == nil {
			curr.Left = newNode
			tree.Size++
			return
		} else {
			q = append(q, curr.Left)
		}

		if curr.Right == nil {
			curr.Right = newNode
			tree.Size++
			return
		} else {
			q = append(q, curr.Right)
		}
	}
}

// TFIND ищет число в дереве обходом по уровням
func TFIND(tree Tree, value int) bool {
	if tree.Root == nil {
		return false
	}

	q := []*TNode{tree.Root}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr.Key == value {
			return true
		}
		if curr.Left != nil {
			q = append(q, curr.Left)
		}
		if curr.Right != nil {
			q = append(q, curr.Right)
		}
	}
	return false
}

// TDEL удаляет последний добавленный элемент дерева
func TDEL(tree *Tree, value int) bool {
	if tree.Root == nil {
		fmt.Println("Нет значений")
		return false
	}

	if tree.Left == nil && tree.Root.Right == nil {
		tree.Root = nil
		tree.Size = 0
		fmt.Println("Значение удалено")
		return true
	}

	q := []*TNode{tree.Root}
	var curr *TNode
	var lastNodeParent *TNode
	var lastNode *TNode

	for len(q) > 0 {
		curr = q[0]
		q = q[1:]

		if curr.Left != nil {
			q = append(q, curr.Left)
			lastNodeParent = curr
			lastNode = curr.Left
		}
		if curr.Right != nil {
			q = append(q, curr.Right)
			lastNodeParent = curr
			lastNode = curr.Right
		}
	}

	if lastNodeParent != nil {
		if lastNodeParent.Right == lastNode {
			lastNodeParent.Right = nil
		} else {
			lastNodeParent.Left = nil
		}
		tree.Size--
		fmt.Println("Значение удалено")
		return true
	}

	return false
}

// TGET_PREORDER Прямой обход: Корень -> Левый -> Правый
func TGET_PREORDER(node *TNode) {
	if node == nil {
		return
	}
	fmt.Print(node.Key, " ")
	TGET_PREORDER(node.Left)
	TGET_PREORDER(node.Right)
}

// TGET_INORDER Симметричный обход: Левый -> Корень -> Правый
func TGET_INORDER(node *TNode) {
	if node == nil {
		return
	}
	TGET_INORDER(node.Left)
	fmt.Print(node.Key, " ")
	TGET_INORDER(node.Right)
}

// TGET_POSTORDER Обратный обход: Левый -> Правый -> Корень
func TGET_POSTORDER(node *TNode) {
	if node == nil {
		return
	}
	TGET_POSTORDER(node.Left)
	TGET_POSTORDER(node.Right)
	fmt.Print(node.Key, " ")
}

// TIS_COMPLETE Проверка на почти полное дерево
func TIS_COMPLETE(tree Tree) bool {
	if tree.Root == nil || tree.Root.Key == -999 {
		return true
	}

	q := []*TNode{tree.Root}
	reachedEnd := false

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr.Left != nil && curr.Left.Key != -999 {
			if reachedEnd {
				return false
			}
			q = append(q, curr.Left)
		} else {
			reachedEnd = true
		}

		if curr.Right != nil && curr.Right.Key != -999 {
			if reachedEnd {
				return false
			}
			q = append(q, curr.Right)
		} else {
			reachedEnd = true
		}
	}
	return true
}

// TIS_FULL Проверка на строго полное дерево
func TIS_FULL(tree Tree) bool {
	if tree.Root == nil || tree.Root.Key == -999 {
		return true
	}

	q := []*TNode{tree.Root}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		hasLeft := curr.Left != nil && curr.Left.Key != -999
		hasRight := curr.Right != nil && curr.Right.Key != -999

		if hasLeft != hasRight {
			return false
		}

		if hasLeft {
			q = append(q, curr.Left)
		}
		if hasRight {
			q = append(q, curr.Right)
		}
	}
	return true
}

// TGET_LEVELORDER Обход по уровням
func TGET_LEVELORDER(tree Tree) {
	if tree.Root == nil {
		return
	}
	q := []*TNode{tree.Root}
	for len(q) > 0 {
		curr := q[0]
		q = q[1:]
		if curr.Key != -999 {
			fmt.Print(curr.Key, " ")
		}
		if curr.Left != nil {
			q = append(q, curr.Left)
		}
		if curr.Right != nil {
			q = append(q, curr.Right)
		}
	}
}
