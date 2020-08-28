package splaytree

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func intComparator(key1, key2 interface{}) int {
	if a, b := key1.(int), key2.(int); a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func floatComparator(key1, key2 interface{}) int {
	if a, b := key1.(float64), key2.(float64); a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func TestCreate(t *testing.T) {
	tree := Create(intComparator)
	if tree == nil {
		t.Errorf("Tree is nil")
	}
	if tree.Root != nil {
		t.Errorf("Root is not nil")
	}
	if tree.Count != 0 {
		t.Errorf("Count is not zero")
	}
	if fmt.Sprintf("%p", intComparator) != fmt.Sprintf("%p", tree.KeyComparator) {
		t.Errorf("Key comparator is invalid")
	}
}

func TestAdd_1(t *testing.T) {
	tree := Create(intComparator)
	tree.Add(1, "value")
	if tree.Root == nil {
		t.Errorf("Root is nil")
	}
	if tree.Count != 1 {
		t.Errorf("Invalid count")
	}
	if tree.Root.Key != 1 {
		t.Errorf("Keys do not match")
	}
	if tree.Root.Value != "value" {
		t.Errorf("Values do not match")
	}
}

func TestAdd_2(t *testing.T) {
	tree := Create(intComparator)
	tree.Add(1, "1")
	tree.Add(5, "5")
	tree.Add(4, "4")

	if tree.Count != 3 {
		t.Errorf("Invalid count")
	}
	if tree.Root.Key != 4 {
		t.Errorf("Invalid root key")
	}
	if tree.Root.Left != nil && tree.Root.Left.Key != 1 {
		t.Errorf("Invalid left key")
	}
	if tree.Root.Right != nil && tree.Root.Right.Key != 5 {
		t.Errorf("Invalid right key")
	}
}

func TestAdd_3(t *testing.T) {
	tree := Create(floatComparator)
	tree.Add(0.5, "0.5")
	tree.Add(1.5, "1.5")
	tree.Add(0.7, "0.7")
	tree.Add(2.5, "2.5")
	tree.Add(7.1, "7.1")
	tree.Add(2.2, "2.2")

	cmp := func(node *Node, expectedKey float64, expectedValue string) bool {
		return node != nil &&
			math.Abs(node.Key.(float64)-expectedKey) < 0.001 &&
			node.Value.(string) == expectedValue
	}

	if tree.Count != 6 {
		t.Errorf("Invalid count")
	}
	if !cmp(tree.Root, 2.2, "2.2") {
		t.Errorf("Invalid root")
	}
	if !cmp(tree.Root.Left, 1.5, "1.5") {
		t.Errorf("Invalid left")
	}
	if !cmp(tree.Root.Left.Left, 0.7, "0.7") {
		t.Errorf("Invalid left left")
	}
	if !cmp(tree.Root.Left.Left.Left, 0.5, "0.5") {
		t.Errorf("Invalid left left left")
	}
	if !cmp(tree.Root.Right, 7.1, "7.1") {
		t.Errorf("Invalid right")
	}
	if !cmp(tree.Root.Right.Left, 2.5, "2.5") {
		t.Errorf("Invalid right left")
	}
}

func TestAdd_Traverse(t *testing.T) {
	const numElements = 10000
	tree := Create(intComparator)

	// Fill
	for i := 0; i < numElements; i++ {
		v := rand.Int()
		tree.Add(v, fmt.Sprintf("random value: %v", v))
	}
	if tree.Count != numElements {
		t.Errorf("Invalid count")
	}

	// Traverse
	keys := make([]int, 0, tree.Count)
	tree.TraverseInorder(func(node *Node) {
		keys = append(keys, node.Key.(int))
	})
	if len(keys) != numElements {
		t.Errorf("Number of keys does not match")
	}
	if !sort.SliceIsSorted(keys, func(a, b int) bool { return a < b }) {
		t.Errorf("Keys are not sorted")
	}

	// Remove
	for _, key := range keys {
		tree.Remove(key)
	}
	if tree.Count != 0 || tree.Root != nil {
		t.Errorf("Tree is not empty")
	}
}

func TestRemove_1(t *testing.T) {
	tree := Create(intComparator)
	tree.Add(5, "5")
	tree.Add(15, "15")
	tree.Add(7, "7")
	tree.Add(25, "25")
	tree.Add(71, "71")
	tree.Add(22, "22")
	tree.Remove(5)

	cmp := func(node *Node, expectedKey int) bool {
		return node != nil && node.Key.(int) == expectedKey
	}

	if tree.Count != 5 {
		t.Errorf("Invalid count")
	}
	if !cmp(tree.Root, 22) {
		t.Errorf("Invalid root")
	}
	if !cmp(tree.Root.Left, 7) {
		t.Errorf("Invalid left")
	}
	if !cmp(tree.Root.Left.Right, 15) {
		t.Errorf("Invalid left right")
	}
	if !cmp(tree.Root.Right, 71) {
		t.Errorf("Invalid right")
	}
	if !cmp(tree.Root.Right.Left, 25) {
		t.Errorf("Invalid right left")
	}

	tree.Remove(7)
	tree.Remove(71)
	if !cmp(tree.Root, 25) {
		t.Errorf("Invalid root")
	}
	if !cmp(tree.Root.Left, 22) {
		t.Errorf("Invalid left")
	}
	if !cmp(tree.Root.Left.Left, 15) {
		t.Errorf("Invalid left left")
	}

	tree.Remove(25)
	tree.Remove(22)
	tree.Remove(15)
	if tree.Count != 0 {
		t.Errorf("Count is not 0")
	}
	if tree.Root != nil {
		t.Errorf("Root is not nil")
	}
}

func TestAddTree_Traverse(t *testing.T) {
	stringCmp := func(s1, s2 interface{}) int {
		return len(s1.(string)) - len(s2.(string))
	}

	tree := Create(stringCmp)
	tree.Add("a", 1)
	tree.Add("aa", 2)
	tree.Add("aaa", 3)
	tree.Add("aaaa", 4)

	treeCopy := Create(stringCmp)
	treeCopy.AddTree(tree)

	if treeCopy.Count != 4 {
		t.Errorf("Invalid count")
	}

	keysInorder := make([]string, 0, tree.Count)
	keysPreorder := make([]string, 0, tree.Count)
	keysPostorder := make([]string, 0, tree.Count)

	treeCopy.TraverseInorder(func(node *Node) {
		keysInorder = append(keysInorder, node.Key.(string))
	})
	treeCopy.TraversePreorder(func(node *Node) {
		keysPreorder = append(keysPreorder, node.Key.(string))
	})
	treeCopy.TraversePostorder(func(node *Node) {
		keysPostorder = append(keysPostorder, node.Key.(string))
	})

	if !reflect.DeepEqual([]string{"a", "aa", "aaa", "aaaa"}, keysInorder) {
		t.Errorf("Invalid inorder")
	}
	if !reflect.DeepEqual([]string{"aaaa", "a", "aa", "aaa"}, keysPreorder) {
		t.Errorf("Invalid preorder")
	}
	if !reflect.DeepEqual([]string{"a", "aa", "aaa", "aaaa"}, keysPostorder) {
		t.Errorf("Invalid postorder")
	}
}
