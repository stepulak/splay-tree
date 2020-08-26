package splaytree

import "testing"

var intComparator = func(key1, key2 interface{}) int {
	if a, b := key1.(int), key2.(int); a < b {
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
	if !funcEqual(intComparator, tree.KeyComparator) {
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
		t.Errorf("Count is not 1")
	}
	if tree.Root.Key != 1 {
		t.Errorf("Keys do not match")
	}
	if tree.Root.Value != "value" {
		t.Errorf("Values do not match")
	}
}
