package splaytree

import (
	"fmt"
)

type Node struct {
	Key interface{}
	Value interface{}
	Parent *Node
	Left *Node
	Right *Node
}

type Tree struct {
	Root *Node
	Counter int
	KeyComparator func (key1 interface{}, key2 interface{}) int
}

