# Splay Tree

Self balancing binary tree with logarithmic amortized time of CRUD operations.

Install this package via: `go get github.com/stepulak/splay-tree`.

Usage example:

```go
import splaytree "github.com/stepulak/splay-tree"

tree := splaytree.Create(func(a, b interface{}) int {
    if a, b := key1.(int), key2.(int); a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
})

tree.Add(1, "value 1")
tree.Add(2, "value 2")
tree.Add(3, "value 3")

fmt.Println(tree)
fmt.Println(tree.Root)
fmt.Println(tree.Count)

node := tree.Find(1)
fmt.Println(node)

tree.Remove(1)
tree.TraverseInorder(func (n *Node) {
    fmt.Println(n)
})
```

