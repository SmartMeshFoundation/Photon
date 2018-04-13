package graph

// Remove a node from a list of nodes... the slow/painful way
func removeNode(node *Node, nodes []*Node) []*Node {
	for i, n := range nodes {
		if n == node {
			return remove(i, nodes)
		}
	}

	return nodes
}

func remove(index int, nodes []*Node) []*Node {
	nodes[index] = nodes[len(nodes)-1]
	return nodes[:len(nodes)-1]
}
