package dijkstra

import (
	"errors"
	"fmt"
)

//GetMapped gets the key assosciated with the mapped int
func (g *Graph) GetMapped(a int) (string, error) {
	if !g.usingMap || g.mapping == nil {
		return "", errors.New("Map is not being used/initialised")
	}
	for k, v := range g.mapping {
		if v == a {
			return k, nil
		}
	}
	return "", errors.New(fmt.Sprint(a, " not found in mapping"))
}

//GetMapping gets the index associated with the specified key
func (g *Graph) GetMapping(a string) (int, error) {
	if !g.usingMap || g.mapping == nil {
		return -1, errors.New("Map is not being used/initialised")
	}
	if b, ok := g.mapping[a]; ok {
		return b, nil
	}
	return -1, errors.New(fmt.Sprint(a, " not found in mapping"))
}

//AddMappedVertex adds a new Vertex with a mapped ID (or returns the index if
// ID already exists).
func (g *Graph) AddMappedVertex(ID string) int {
	if !g.usingMap || g.mapping == nil {
		g.usingMap = true
		g.mapping = map[string]int{}
		g.highestMapIndex = 0
	}
	if i, ok := g.mapping[ID]; ok {
		return i
	}
	i := g.highestMapIndex
	g.highestMapIndex++
	g.mapping[ID] = i
	return g.AddVertex(i).ID
}

//AddMappedArc adds a new Arc from Source to Destination, for when verticies are
// referenced by strings.
func (g *Graph) AddMappedArc(Source, Destination string, Distance int64) error {
	return g.AddArc(g.AddMappedVertex(Source), g.AddMappedVertex(Destination), Distance)
}

//AddArc is the default method for adding an arc from a Source Vertex to a
// Destination Vertex
func (g *Graph) AddArc(Source, Destination int, Distance int64) error {
	if len(g.Verticies) <= Source || len(g.Verticies) <= Destination {
		return errors.New("Source/Destination not found")
	}
	g.Verticies[Source].AddArc(Destination, Distance)
	return nil
}

//delete one arc from graph
func (g *Graph) DeleteArc(Source, Destination int) error {
	if len(g.Verticies) <= Source || len(g.Verticies) <= Destination {
		return errors.New("Source/Destination not found")
	}
	g.Verticies[Source].DeleteArc(Destination)
	return nil
}

func (g *Graph) GetAllNeighbors(Source int) (neighbors []int, err error) {
	if len(g.Verticies) <= Source {
		err = errors.New("Source/Destination not found")
		return
	}
	for k, _ := range g.Verticies[Source].arcs {
		neighbors = append(neighbors, k)
	}
	return neighbors, nil
}
