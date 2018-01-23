package dijkstra

import "github.com/chanxuehong/util/math"

func (g *Graph) ShortestPath(src, dest int) (x int) {
	//a clean visit memory
	g.visited = make(map[int]bool)
	g.visit(src)
	v := g.vertices[src]
	h := make(Candidates, len(v.arcs))
	// initialize the heap with out edges from src
	for id, y := range v.arcs {
		v := g.vertices[id]
		// update the vertices being pointed to with the distance.
		v.dist = y
		g.vertices[id] = v
		h.Push(v)
	}
	for src != dest {
		if h.IsEmpty() {
			return math.MaxInt
		}
		v = h.Pop()
		src = v.id
		if g.visited[src] {
			continue
		}
		g.visit(src)
		for w, d := range v.arcs {
			if g.visited[w] {
				continue
			}
			c := g.vertices[w]
			distance := d + v.dist
			if distance < c.dist {
				c.dist = distance
				g.vertices[w] = c
			}
			h.Push(c)
		}
	}
	v = g.vertices[dest]
	return v.dist
}
