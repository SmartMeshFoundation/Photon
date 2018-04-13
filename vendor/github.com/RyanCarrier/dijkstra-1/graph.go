package dijkstra

import (
	"bufio"
	// "container/heap"
	"log"
	"os"
	"strconv"
	"strings"
)

type Edge struct {
	Tail   int
	Head   int
	Weight int
}

type Vertex struct {
	ID   int
	Dist int
	Arcs map[int]int // arcs[vertex id] = weight
}

type Candidates []Vertex

func (h Candidates) Len() int           { return len(h) }
func (h Candidates) Less(i, j int) bool { return h[i].Dist < h[j].Dist }
func (h Candidates) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Candidates) IsEmpty() bool      { return len(h) == 0 }

func (h *Candidates) Push(v Vertex) {
	var changed bool
	old := *h
	updated := *h
	// insert Vertex at the correct position (keyed by distance)
	for i, w := range old {
		if v.ID == w.ID {
			if changed {
				if i+1 < len(updated) {
					updated = append(updated[:i], updated[i+1:]...)
				} else {
					updated = updated[:i]
				}
			} else if v.Dist < w.Dist {
				updated[i] = v
			}
			changed = true
		} else if v.Dist < w.Dist {
			changed = true
			updated = append(old[:i], v)
			updated = append(updated, w)
			updated = append(updated, old[i+1:]...)
		}
	}
	if !changed {
		updated = append(old, v)
	}
	*h = updated
}

func (h *Candidates) Pop() (v Vertex) {
	old := *h
	v = old[0]
	*h = old[1:]
	return
}

type Graph struct {
	visited  map[int]bool
	vertices map[int]Vertex
}

func NewGraph(vs map[int]Vertex) *Graph {
	g := new(Graph)
	g.visited = make(map[int]bool)
	g.vertices = make(map[int]Vertex)
	for i, v := range vs {
		v.Dist = 1000000
		g.vertices[i] = v
	}
	return g
}

func NewGraphFromFile(fn string) *Graph {
	v := make(map[int]Vertex)
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(string(line), "\t")
		t := words[0]
		tail, err := strconv.Atoi(t)
		// a Vertex is an id (tail) and a map of arcs
		arcs := make(map[int]int)
		for i := 1; i < len(words); i++ {
			item := strings.Split(words[i], ",")
			if len(item) < 2 {
				break
			}
			h, w := item[0], item[1]
			head, err := strconv.Atoi(h)
			if err != nil {
				log.Print(err)
			}
			weight, err := strconv.Atoi(w)
			if err != nil {
				log.Print(err)
			}
			arcs[head] = weight
		}
		v[tail] = Vertex{ID: tail, Arcs: arcs, Dist: 0}
		if err != nil {
			log.Print(err)
		}
		err = nil
	}
	return NewGraph(v)
}

func (g *Graph) Len() int    { return len(g.vertices) }
func (g *Graph) visit(v int) { g.visited[v] = true }
