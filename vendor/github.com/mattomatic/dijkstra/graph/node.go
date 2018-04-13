package graph

type Node struct {
	Id      int
	Visited bool
	edges   []*Edge
}

func NewNode(id int) *Node {
	return &Node{
		Id:    id,
		edges: make([]*Edge, 0),
	}
}

func (n *Node) AddEdges(edges ...*Edge) {
	for _, edge := range edges {
		n.edges = append(n.edges, edge)
	}
}

func (n *Node) GetEdges() chan *Edge {
	edges := make(chan *Edge)

	go func() {
		defer close(edges)
		for _, edge := range n.edges {
			edges <- edge
		}
	}()

	return edges
}
