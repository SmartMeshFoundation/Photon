package graph

type Graph struct {
	nodes []*Node
}

func NewGraph() *Graph {
	return &Graph{
		make([]*Node, 0),
	}
}

func (g *Graph) Size() int {
	return len(g.nodes)
}

func (g *Graph) Contains(n *Node) bool {
	for _, node := range g.nodes {
		if node == n {
			return true
		}
	}

	return false
}

func (g *Graph) AddNodes(nodes ...*Node) {
	for _, node := range nodes {
		g.nodes = append(g.nodes, node)
	}
}

func (g *Graph) RemoveNodes(nodes ...*Node) {
	for _, node := range nodes {
		g.nodes = removeNode(node, g.nodes)
	}
}

func (g *Graph) GetNodes() chan *Node {
	nodes := make(chan *Node)

	go func() {
		defer close(nodes)
		for _, node := range g.nodes {
			nodes <- node
		}
	}()

	return nodes
}

func (g *Graph) GetEdges() chan *Edge {
	edges := make(chan *Edge)

	go func() {
		defer close(edges)
		for _, node := range g.nodes {
			for edge := range node.GetEdges() {
				edges <- edge
			}
		}
	}()

	return edges
}

// Get the edges that cross from g to o
func (g *Graph) GetCut(o *Graph) chan *Edge {
	edges := make(chan *Edge)

	go func() {
		defer close(edges)
		for edge := range g.GetEdges() {
			if o.Contains(edge.Tail) {
				edges <- edge
			}
		}
	}()

	return edges
}

func (g *Graph) Search(id int) *Node {
	for _, node := range g.nodes {
		if node.Id == id {
			return node
		}
	}

	return nil
}
