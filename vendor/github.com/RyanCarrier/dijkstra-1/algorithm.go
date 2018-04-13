package dijkstra

func (g *Graph) ShortestPath(src, dest int) (x int) {
	g.visit(src)
	v := g.vertices[src]
	h := make(Candidates, len(v.Arcs))
	// initialize the heap with out edges from src
	for id, y := range v.Arcs {
		v := g.vertices[id]
		// update the vertices being pointed to with the distance.
		v.Dist = y
		g.vertices[id] = v
		h.Push(v)
	}
	for src != dest {
		if h.IsEmpty() {
			return 1000000
		}
		v = h.Pop()
		src = v.ID
		if g.visited[src] {
			continue
		}
		g.visit(src)
		for w, d := range v.Arcs {
			if g.visited[w] {
				continue
			}
			c := g.vertices[w]
			distance := d + v.Dist
			if distance < c.Dist {
				c.Dist = distance
				g.vertices[w] = c
			}
			h.Push(c)
		}
	}
	v = g.vertices[dest]
	return v.Dist
}
