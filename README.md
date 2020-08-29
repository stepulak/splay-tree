# Splay Tree

Self balancing binary tree with logarithmic amortized time of CRUD operations.

Install this package via: `go get github.com/stepulak/splay-tree`

Usage example:

```go
package main

import (
	"fmt"

	splaytree "github.com/stepulak/splay-tree"
)

func main() {
	tree := splaytree.Create(func(a, b interface{}) int {
		if a, b := a.(int), b.(int); a < b {
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
	tree.TraverseInorder(func(n *splaytree.Node) {
		fmt.Println(n)
	})
}
```

For more usage info look at source code `splaytree.go` and `splaytree_test.go`.