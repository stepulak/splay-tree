package splaytree

type Node struct {
	Key    interface{}
	Value  interface{}
	Parent *Node
	Left   *Node
	Right  *Node
}

func (n *Node) isLeftChild() bool {
	return n.Parent != nil && n.Parent.Left == n
}

type Tree struct {
	Root          *Node
	Counter       int
	KeyComparator func(key1 interface{}, key2 interface{}) int
}

func (t *Tree) findNodeRec(key interface{}, node *Node, parent *Node) (*Node, *Node) {
	if cmp := t.KeyComparator(key, node.Key); cmp == 0 {
		return node, node.Parent
	} else if cmp < 0 && node.Left != nil {
		return t.findNodeRec(key, node.Left, node)
	} else if cmp >= 0 && node.Right != nil {
		return t.findNodeRec(key, node.Right, node)
	}
	return nil, node
}

func (t *Tree) rightRotation(node *Node) {
	par, right := node.Parent, node.Right
	node.Right = par
	par.Left = right
	if right != nil {
		right.Parent = par
	}
	t.swapGrandparent(node, par)
}

func (t *Tree) leftRotation(node *Node) {
	par, left := node.Parent, node.Left
	node.Left = par
	par.Right = left
	if left != nil {
		left.Parent = par
	}
	t.swapGrandparent(node, par)
}

func (t *Tree) splay(node *Node) {
	for node.Parent != nil {
		par := node.Parent
		if par.Parent == nil {
			// Zig step
			if node.isLeftChild() {
				t.rightRotation(node)
			} else {
				t.leftRotation(node)
			}
		} else if node.isLeftChild() && par.isLeftChild() {
			// Zig-zig step
			t.rightRotation(par)
			t.rightRotation(node)
		} else if !node.isLeftChild() && !par.isLeftChild() {
			// Zig-zig step
			t.leftRotation(par)
			t.leftRotation(node)
		} else if node.isLeftChild() && !par.isLeftChild() {
			// Zig-zag step
			t.leftRotation(node)
			t.rightRotation(par)
		}
	}
}

func (t *Tree) insert(key, value interface{}, parent *Node) *Node {
	node := &Node{}
	if t.KeyComparator(key, parent.Key) < 0 {
		parent.Left = node
	} else {
		parent.Right = node
	}
	t.splay(node)
	t.Counter++
	return node
}

func (t *Tree) Add(key, value interface{}) *Node {
	// Add root if does not exists
	if t.Root == nil {
		t.Root = &Node{key, value, nil, nil, nil}
		t.Counter++
		return t.Root
	}
	node, par := t.findNodeRec(key, t.Root, nil)
	if node != nil {
		// Node with same key found, just replace the value
		node.Value = value
		return node
	}
	return t.insert(key, value, par)
}
