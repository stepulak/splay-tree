package splaytree

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"testing"
)

func intComparator(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func floatComparator(a, b float64) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func TestCreate(t *testing.T) {
	tree := Create[int, int](intComparator)
	if tree == nil {
		t.Error("Tree is nil")
	}
	if tree.Root != nil {
		t.Error("Root is not nil")
	}
	if tree.Count != 0 {
		t.Error("Count is not zero")
	}
	if fmt.Sprintf("%p", intComparator) != fmt.Sprintf("%p", tree.KeyComparator) {
		t.Error("Key comparator is invalid")
	}
}

func TestAdd_1(t *testing.T) {
	tree := Create[int, string](intComparator)
	tree.Add(1, "value")
	if tree.Root == nil {
		t.Error("Root is nil")
	}
	if tree.Count != 1 {
		t.Error("Invalid count")
	}
	if tree.Root.Key != 1 {
		t.Error("Keys do not match")
	}
	if tree.Root.Value != "value" {
		t.Error("Values do not match")
	}
}

func TestAdd_2(t *testing.T) {
	tree := Create[int, string](intComparator)
	tree.Add(1, "1")
	tree.Add(5, "5")
	tree.Add(4, "4")

	if tree.Count != 3 {
		t.Error("Invalid count")
	}
	if tree.Root.Key != 4 {
		t.Error("Invalid root key")
	}
	if tree.Root.Left != nil && tree.Root.Left.Key != 1 {
		t.Error("Invalid left key")
	}
	if tree.Root.Right != nil && tree.Root.Right.Key != 5 {
		t.Error("Invalid right key")
	}
}

func TestAdd_3(t *testing.T) {
	tree := Create[float64, string](floatComparator)
	tree.Add(0.5, "0.5")
	tree.Add(1.5, "1.5")
	tree.Add(0.7, "0.7")
	tree.Add(2.5, "2.5")
	tree.Add(7.1, "7.1")
	tree.Add(2.2, "2.2")

	cmp := func(node *Node[float64, string], expectedKey float64, expectedValue string) bool {
		return node != nil &&
			math.Abs(node.Key-expectedKey) < 0.001 &&
			node.Value == expectedValue
	}

	if tree.Count != 6 {
		t.Error("Invalid count")
	}
	if !cmp(tree.Root, 2.2, "2.2") {
		t.Error("Invalid root")
	}
	if !cmp(tree.Root.Left, 1.5, "1.5") {
		t.Error("Invalid left")
	}
	if !cmp(tree.Root.Left.Left, 0.7, "0.7") {
		t.Error("Invalid left left")
	}
	if !cmp(tree.Root.Left.Left.Left, 0.5, "0.5") {
		t.Error("Invalid left left left")
	}
	if !cmp(tree.Root.Right, 7.1, "7.1") {
		t.Error("Invalid right")
	}
	if !cmp(tree.Root.Right.Left, 2.5, "2.5") {
		t.Error("Invalid right left")
	}
}

func TestAdd_Traverse(t *testing.T) {
	const numElements = 10000
	tree := Create[int, string](intComparator)

	// Fill
	for i := 0; i < numElements; i++ {
		v := rand.Int()
		tree.Add(v, fmt.Sprintf("random value: %v", v))
	}
	if tree.Count != numElements {
		t.Error("Invalid count")
	}

	// Traverse
	keys := make([]int, 0, tree.Count)
	tree.TraverseInorder(func(node *Node[int, string]) {
		keys = append(keys, node.Key)
	})
	if len(keys) != numElements {
		t.Error("Number of keys does not match")
	}
	if !sort.SliceIsSorted(keys, func(a, b int) bool { return a < b }) {
		t.Error("Keys are not sorted")
	}

	// Remove
	for _, key := range keys {
		tree.Remove(key)
	}
	if err := tree.Remove(123); err == nil || err != TreeError(NotFoundError) {
		t.Error("Expected NotFoundError")
	}
	if tree.Count != 0 || tree.Root != nil {
		t.Error("Tree is not empty")
	}
}

func TestRemove_1(t *testing.T) {
	tree := Create[int, string](intComparator)
	tree.Add(5, "5")
	tree.Add(15, "15")
	tree.Add(7, "7")
	tree.Add(25, "25")
	tree.Add(71, "71")
	tree.Add(22, "22")
	tree.Remove(5)

	cmp := func(node *Node[int, string], expectedKey int) bool {
		return node != nil && node.Key == expectedKey
	}

	if tree.Count != 5 {
		t.Error("Invalid count")
	}
	if !cmp(tree.Root, 22) {
		t.Error("Invalid root")
	}
	if !cmp(tree.Root.Left, 7) {
		t.Error("Invalid left")
	}
	if !cmp(tree.Root.Left.Right, 15) {
		t.Error("Invalid left right")
	}
	if !cmp(tree.Root.Right, 71) {
		t.Error("Invalid right")
	}
	if !cmp(tree.Root.Right.Left, 25) {
		t.Error("Invalid right left")
	}

	tree.Remove(7)
	tree.Remove(71)
	if !cmp(tree.Root, 25) {
		t.Error("Invalid root")
	}
	if !cmp(tree.Root.Left, 22) {
		t.Error("Invalid left")
	}
	if !cmp(tree.Root.Left.Left, 15) {
		t.Error("Invalid left left")
	}

	tree.Remove(25)
	tree.Remove(22)
	tree.Remove(15)
	if tree.Count != 0 {
		t.Error("Count is not 0")
	}
	if tree.Root != nil {
		t.Error("Root is not nil")
	}
}

func TestFind_RemoveNode(t *testing.T) {
	tree := Create[int, interface{}](intComparator)

	tree.Add(1, nil)
	tree.Add(2, nil)
	tree.Add(3, nil)
	tree.Add(-1, nil)
	tree.Add(-2, nil)
	tree.Add(-3, nil)

	for _, key := range []int{3, -2, -1, 1, -3, 2} {
		node := tree.Find(key)
		if node == nil {
			t.Error("Found node is nil")
		}
		tree.RemoveNode(node)
	}

	if tree.Count != 0 {
		t.Error("Tree is not empty")
	}
}

func TestAddTree_Traverse(t *testing.T) {
	stringCmp := func(s1, s2 string) int {
		return len(s1) - len(s2)
	}

	tree := Create[string, int](stringCmp)
	tree.Add("a", 1)
	tree.Add("aa", 2)
	tree.Add("aaa", 3)
	tree.Add("aaaa", 4)

	treeCopy := Create[string, int](stringCmp)
	treeCopy.AddTree(tree)

	if treeCopy.Count != 4 {
		t.Error("Invalid count")
	}

	keysInorder := make([]string, 0, tree.Count)
	keysPreorder := make([]string, 0, tree.Count)
	keysPostorder := make([]string, 0, tree.Count)

	treeCopy.TraverseInorder(func(node *Node[string, int]) {
		keysInorder = append(keysInorder, node.Key)
	})
	treeCopy.TraversePreorder(func(node *Node[string, int]) {
		keysPreorder = append(keysPreorder, node.Key)
	})
	treeCopy.TraversePostorder(func(node *Node[string, int]) {
		keysPostorder = append(keysPostorder, node.Key)
	})

	if !reflect.DeepEqual([]string{"a", "aa", "aaa", "aaaa"}, keysInorder) {
		t.Error("Invalid inorder")
	}
	if !reflect.DeepEqual([]string{"aaaa", "a", "aa", "aaa"}, keysPreorder) {
		t.Error("Invalid preorder")
	}
	if !reflect.DeepEqual([]string{"a", "aa", "aaa", "aaaa"}, keysPostorder) {
		t.Error("Invalid postorder")
	}
}

func TestString(t *testing.T) {
	type Value struct {
		myValue string
	}

	tree := Create[int, Value](intComparator)
	tree.Add(1, Value{"x"})
	tree.Add(2, Value{"y"})
	tree.Add(3, Value{"z"})

	str := tree.String()
	rx := regexp.MustCompile(`Key:\s(\d)+;\sVal:\s\{(\w)+\};`)

	expectedKeys := []int{1, 2, 3}
	expectedValues := []string{"x", "y", "z"}
	matches := rx.FindAllStringSubmatch(str, -1)

	if len(matches) != 3 {
		t.Error("Invalid number of matches")
	}

	for i, match := range matches {
		if len(match) != 3 {
			t.Error("Invalid match")
		}
		if key, err := strconv.Atoi(match[1]); err != nil || expectedKeys[i] != key {
			t.Errorf("Invalid key: %v expected: %v", key, expectedKeys[i])
		}
		if expectedValues[i] != match[2] {
			t.Errorf("Invalid value: %v expected: %v", match[2], expectedValues[i])
		}
	}
}

func TestUncomparableKey(t *testing.T) {
	type Key map[string]string

	tree := Create[Key, int](func(a, b Key) int {
		aJson, _ := json.Marshal(a)
		bJson, _ := json.Marshal(b)
		return len(aJson) - len(bJson)
	})

	tree.Add(Key{"Hello": "world"}, 1)
	tree.Add(Key{"Go": "lang"}, 2)
	tree.Add(Key{"Praise": "God"}, 33)

	if tree.Count != 3 {
		t.Error("Invalid count")
	}
	if node := tree.Find(Key{"Hello": "world"}); node == nil || node.Value != 1 {
		t.Error("Unable to find/invalid {'Hello': 'world'}")
	}
	if tree.Remove(Key{"Go": "lang"}) != nil {
		t.Error("Unexpected error during node removal")
	}
	if node := tree.Find(Key{"Go": "lang"}); node != nil {
		t.Error("Able to find removed {'Go': 'lang'}")
	}
}
