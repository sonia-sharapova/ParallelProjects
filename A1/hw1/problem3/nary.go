// References:
// Trees in Go: https://ieftimov.com/posts/golang-datastructures-trees/
// Creating objects: https://stackoverflow.com/questions/39397865/how-to-create-object-for-a-struct-in-golang
// Goland fields: https://www.geeksforgeeks.org/strings-fields-function-in-golang-with-examples/

package problem3

import (
	"strings"
)

// You cannot modify the name of this struct
type Node struct {
	children []*Node // We will store the children of each node here
}

// You cannot modify this data structure at all (including its types). However
// You are not required to use all of its fields.
type Tree struct {
	root  *Node
	nAry  int
	nodes map[string]*Node
}

// SearchTree takes in a tree and returns the number of internal nodes that have children that are greater
// than the arty of the tree, which is given by the second argument to the function.
// Note: Please make sure to read the assignment description for more details.
func SearchTree(edges []string, nAry int) int {

	// Make a tree object
	tree := Tree{
		nAry:  nAry,
		nodes: make(map[string]*Node),
	}

	// iterate over edges to build tree
	for _, edge := range edges {

		// get parent and children
		parts := strings.Fields(edge)
		parent := parts[0]
		child := parts[1]

		parentNode, exists := tree.nodes[parent]
		if !exists { // make node if not present
			parentNode = &Node{}
			tree.nodes[parent] = parentNode
		}

		childNode, exists := tree.nodes[child]
		if !exists {
			childNode = &Node{} // make node if not present
			tree.nodes[child] = childNode
		}

		// add child
		parentNode.children = append(parentNode.children, childNode)

		// index 0 to represent root
		if tree.root == nil {
			tree.root = parentNode
		}
	}

	// count the nodes with more children than given arity
	tooMany := 0
	for _, node := range tree.nodes {
		if len(node.children) > nAry {
			tooMany++
		}
	}

	return tooMany
}
