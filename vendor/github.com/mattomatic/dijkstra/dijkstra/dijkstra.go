package dijkstra

import (
	"github.com/mattomatic/dijkstra/graph"
)

type PathMap map[*graph.Node]int

const MAXWEIGHT = 1000000 // a token maximum value

// Compute the shortest path from the source node to each other node in the graph
func Dijkstra(g *graph.Graph, s *graph.Node) PathMap {
	p := make(PathMap)
	v := graph.NewGraph()
	u := g

	precompute(s, g, p)
	initialize(v, u, s, p)

	for u.Size() > 0 {
		step(s, v, u, p)
	}

	return p
}

// Determine which nodes are unreachable from s and set their values in p
// to be MAXWEIGHT
func precompute(s *graph.Node, g *graph.Graph, p PathMap) {
	dfs(s)

	for node := range g.GetNodes() {
		if node.Visited == false {
			p[node] = MAXWEIGHT
			g.RemoveNodes(node)
		}
	}
}

func dfs(node *graph.Node) {
	if node.Visited {
		return
	}

	node.Visited = true

	for edge := range node.GetEdges() {
		dfs(edge.Tail)
	}
}

func initialize(v *graph.Graph, u *graph.Graph, s *graph.Node, p PathMap) {
	v.AddNodes(s)
	u.RemoveNodes(s)
	p[s] = 0
}

func step(s *graph.Node, v *graph.Graph, u *graph.Graph, p PathMap) {
	var minNode *graph.Node
	var minScore int = MAXWEIGHT

	for edge := range v.GetCut(u) {
		score := scoreEdge(edge, p)

		if score < minScore {
			minNode = edge.Tail
			minScore = score
		}
	}

	if minScore == MAXWEIGHT {
		panic("something bad happened")
	}

	// Choose the min edge and update its score
	v.AddNodes(minNode)
	u.RemoveNodes(minNode)
	p[minNode] = minScore
}

func scoreEdge(edge *graph.Edge, p PathMap) int {
	score, _ := p[edge.Head]
	return score + edge.Weight
}
