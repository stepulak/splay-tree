package splaytree

import (
	"fmt"
	"strings"
)

// Comparator is a comparing function that:
// returns -1 if a < b
// returns 1 if a > b
// returns 0 if a == 0
// FYI: `TKey` is `any` to support complex uncomparable types by default
type Comparator[TKey any] func(a TKey, b TKey) int

// Tree structure represents the Splay tree
// Root: pointer to root node
// Count: number of elements
// KeyComparator: function for comparing keys, see @Comparator
// FYI: `TKey` is `any` to support complex uncomparable types by default
type Tree[TKey any, TValue any] struct {
	Root          *Node[TKey, TValue]
	Count         int
	KeyComparator Comparator[TKey]
}

// Node structure that represents the key-value pair in Splay tree
// Key: key for node comparison
// Value: value stored in node
// Parent: pointer to parent node
// Left: pointer to left child node
// Right: pointer to right child node
type Node[TKey any, TValue any] struct {
	Key    TKey
	Value  TValue
	Parent *Node[TKey, TValue]
	Left   *Node[TKey, TValue]
	Right  *Node[TKey, TValue]
}

// NodeVisitor is a function for visiting nodes
type NodeVisitor[TKey any, TValue any] func(n *Node[TKey, TValue])

// TreeError type
type TreeError string

// TreeError's possible states
const (
	NotFoundError = "Not found"
	// Add more if needed
)

// Create new Splay tree with given key comparator
func Create[TKey any, TValue any](keyComparator Comparator[TKey]) *Tree[TKey, TValue] {
	return &Tree[TKey, TValue]{KeyComparator: keyComparator}
}

// Find node with given key in tree
// If node is not found, return nil
func (t *Tree[TKey, TValue]) Find(key TKey) *Node[TKey, TValue] {
	if t.Root == nil {
		return nil
	}
	node, _ := t.findNodeRec(key, t.Root, nil)
	return node
}

// Add new key-value pair into tree
// If node with given key already exists, then just update it's value with new one
// Return pointer to newly added node
func (t *Tree[TKey, TValue]) Add(key TKey, value TValue) *Node[TKey, TValue] {
	// Add root if does not exists
	if t.Root == nil {
		t.Root = &Node[TKey, TValue]{key, value, nil, nil, nil}
		t.Count = 1
		return t.Root
	}
	node, parent := t.findNodeRec(key, t.Root, nil)
	if node != nil {
		// Node with same key found, just replace the value
		node.Value = value
		return node
	}
	return t.insertNode(&Node[TKey, TValue]{Key: key, Value: value, Parent: parent})
}

// AddTree - add all elements from given tree into this tree
func (t *Tree[TKey, TValue]) AddTree(tree *Tree[TKey, TValue]) {
	tree.TraverseInorder(func(n *Node[TKey, TValue]) {
		t.Add(n.Key, n.Value)
	})
}

// Remove node with given key from this tree
// If node with given key is not found, then return NotFoundError
func (t *Tree[TKey, TValue]) Remove(key TKey) error {
	node := t.Find(key)
	if node == nil {
		return TreeError(NotFoundError)
	}
	t.RemoveNode(node)
	return nil
}

// RemoveNode - remove given node from this tree
// Node is not validated or checked if it belongs to this tree
func (t *Tree[TKey, TValue]) RemoveNode(node *Node[TKey, TValue]) {
	t.splay(node)
	t.joinSubtrees(node.Left, node.Right)
	t.Count--
}

// TraverseInorder - inorder traverse through tree using visitor
// See https://en.wikipedia.org/wiki/Tree_traversal#Implementations
func (t *Tree[TKey, TValue]) TraverseInorder(visitor NodeVisitor[TKey, TValue]) {
	if t.Root != nil {
		t.Root.traverseInorder(visitor)
	}
}

// TraversePreorder - preorder traverse through tree using visitor
// See https://en.wikipedia.org/wiki/Tree_traversal#Implementations
func (t *Tree[TKey, TValue]) TraversePreorder(visitor NodeVisitor[TKey, TValue]) {
	if t.Root != nil {
		t.Root.traversePreorder(visitor)
	}
}

// TraversePostorder - postorder traverse through tree using visitor
// See https://en.wikipedia.org/wiki/Tree_traversal#Implementations
func (t *Tree[TKey, TValue]) TraversePostorder(visitor NodeVisitor[TKey, TValue]) {
	if t.Root != nil {
		t.Root.traversePostorder(visitor)
	}
}

// Print this tree into string
// Format: { [NodeInfo]* }
func (t *Tree[TKey, TValue]) String() string {
	sb := strings.Builder{}
	sb.WriteString("{ ")
	t.TraverseInorder(func(n *Node[TKey, TValue]) {
		sb.WriteString(n.String())
		sb.WriteString(" ")
	})
	sb.WriteString("}")
	return sb.String()
}

// Print TreeError status
func (e TreeError) Error() string {
	return fmt.Sprintf("Node error: %v", string(e))
}

// Print node's detailed information into string
func (n *Node[TKey, TValue]) String() string {
	return fmt.Sprintf(
		"[Key: %v; Val: %v; Ptr: %p; Par: %p; L: %p; R: %p]",
		n.Key, n.Value, n, n.Parent, n.Left, n.Right)
}

func (n *Node[TKey, TValue]) isLeftChild() bool {
	return n.Parent != nil && n.Parent.Left == n
}

func (n *Node[TKey, TValue]) mostRightChild() *Node[TKey, TValue] {
	if n.Right == nil {
		return n
	}
	return n.Right.mostRightChild()
}

func (n *Node[TKey, TValue]) traverseInorder(visitor NodeVisitor[TKey, TValue]) {
	if n.Left != nil {
		n.Left.traverseInorder(visitor)
	}
	visitor(n)
	if n.Right != nil {
		n.Right.traverseInorder(visitor)
	}
}

func (n *Node[TKey, TValue]) traversePreorder(visitor NodeVisitor[TKey, TValue]) {
	visitor(n)
	if n.Left != nil {
		n.Left.traverseInorder(visitor)
	}
	if n.Right != nil {
		n.Right.traverseInorder(visitor)
	}
}

func (n *Node[TKey, TValue]) traversePostorder(visitor NodeVisitor[TKey, TValue]) {
	if n.Left != nil {
		n.Left.traverseInorder(visitor)
	}
	if n.Right != nil {
		n.Right.traverseInorder(visitor)
	}
	visitor(n)
}

func (t *Tree[TKey, TValue]) findNodeRec(key TKey, node, parent *Node[TKey, TValue]) (*Node[TKey, TValue], *Node[TKey, TValue]) {
	if cmp := t.KeyComparator(key, node.Key); cmp == 0 {
		return node, node.Parent
	} else if cmp < 0 && node.Left != nil {
		return t.findNodeRec(key, node.Left, node)
	} else if cmp >= 0 && node.Right != nil {
		return t.findNodeRec(key, node.Right, node)
	}
	return nil, node
}

func (t *Tree[TKey, TValue]) swapGrandparent(node, parent *Node[TKey, TValue]) {
	grandPar := parent.Parent
	node.Parent = grandPar
	if grandPar != nil {
		if parent.isLeftChild() {
			grandPar.Left = node
		} else {
			grandPar.Right = node
		}
	} else {
		t.Root = node
	}
	parent.Parent = node
}

func (t *Tree[TKey, TValue]) rightRotation(node *Node[TKey, TValue]) {
	parent, right := node.Parent, node.Right
	node.Right = parent
	parent.Left = right
	if right != nil {
		right.Parent = parent
	}
	t.swapGrandparent(node, parent)
}

func (t *Tree[TKey, TValue]) leftRotation(node *Node[TKey, TValue]) {
	parent, left := node.Parent, node.Left
	node.Left = parent
	parent.Right = left
	if left != nil {
		left.Parent = parent
	}
	t.swapGrandparent(node, parent)
}

func (t *Tree[TKey, TValue]) splay(node *Node[TKey, TValue]) {
	for node.Parent != nil {
		parent := node.Parent
		if parent.Parent == nil {
			// Zig step
			if node.isLeftChild() {
				t.rightRotation(node)
			} else {
				t.leftRotation(node)
			}
		} else if node.isLeftChild() && parent.isLeftChild() {
			// Zig-zig step
			t.rightRotation(parent)
			t.rightRotation(node)
		} else if !node.isLeftChild() && !parent.isLeftChild() {
			// Zig-zig step
			t.leftRotation(parent)
			t.leftRotation(node)
		} else if node.isLeftChild() && !parent.isLeftChild() {
			// Zig-zag step
			t.rightRotation(node)
			t.leftRotation(node)
		} else {
			t.leftRotation(node)
			t.rightRotation(node)
		}
	}
}

func (t *Tree[TKey, TValue]) joinSubtrees(left, right *Node[TKey, TValue]) {
	if left != nil {
		left.Parent = nil
		if left.Right != nil {
			mostRight := left.Right.mostRightChild()
			t.splay(mostRight)
			t.Root = mostRight
		} else {
			t.Root = left
		}
		if right != nil {
			t.Root.Right = right
			right.Parent = t.Root
		}
	} else if right != nil {
		t.Root = right
		t.Root.Parent = nil
	} else {
		// Last element?
		t.Root = nil
	}
}

func (t *Tree[TKey, TValue]) insertNode(node *Node[TKey, TValue]) *Node[TKey, TValue] {
	if t.KeyComparator(node.Key, node.Parent.Key) < 0 {
		node.Parent.Left = node
	} else {
		node.Parent.Right = node
	}
	t.splay(node)
	t.Count++
	return node
}
