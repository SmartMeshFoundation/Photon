package graph

type Edge struct {
	Head   *Node
	Tail   *Node
	Weight int
}

func NewEdge(head, tail *Node, weight int) *Edge {
	return &Edge{
		head,
		tail,
		weight,
	}
}
