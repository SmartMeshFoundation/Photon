package dijkstra

import (
	"testing"

	"github.com/chanxuehong/util/math"
)

func AssertEqual(a, b interface{}, t *testing.T) {
	if a != b {
		t.Logf("%v == %v", a, b)
		t.Fail()
	}
}

func TestAssertedLengthofDataFile(t *testing.T) {
	const expected = 200
	g := NewGraphFromFile("dijkstraData.txt")
	result := g.Len()
	AssertEqual(expected, result, t)
}

func TestShortestPathIsLowestSingleArc(t *testing.T) {
	const expected = 15
	v := map[int]Vertex{
		0: {
			id: 0,
			arcs: map[int]int{
				1: 14,
				2: 15,
			},
		},
		1: {
			id: 1,
			arcs: map[int]int{
				2: 8,
			},
		},
		2: {
			id:   2,
			arcs: map[int]int{},
		},
	}
	g := NewGraph(v)
	result := g.ShortestPath(0, 2)
	AssertEqual(result, expected, t)
	if t.Failed() {
		t.Log(g)
	}
}

func TestShortestTakesAHop(t *testing.T) {
	const expected = 10
	v := map[int]Vertex{
		0: {
			id: 0,
			arcs: map[int]int{
				1: 2,
				2: 15,
			},
		},
		1: {
			id: 1,
			arcs: map[int]int{
				2: 8,
			},
		},
		2: {
			id:   2,
			arcs: map[int]int{},
		},
	}
	g := NewGraph(v)
	result := g.ShortestPath(0, 2)
	AssertEqual(result, expected, t)
	if t.Failed() {
		t.Log(g)
	}
}

func TestUnnonnectedNodeReturnsMillion(t *testing.T) {
	const expected = math.MaxInt
	v := map[int]Vertex{
		0: {
			id:   0,
			arcs: map[int]int{},
		},
		1: {
			id:   1,
			arcs: map[int]int{},
		},
	}
	g := NewGraph(v)
	result := g.ShortestPath(0, 2)
	AssertEqual(result, expected, t)
	if t.Failed() {
		t.Log(g)
	}
}

func TestGraph_AllShortestPath(t *testing.T) {
	v := map[int]Vertex{
		0: {
			id: 0,
			arcs: map[int]int{
				1: 1,
				2: 1,
			},
		},
		1: {
			id: 1,
			arcs: map[int]int{
				0: 1,
				3: 1,
			},
		},
		2: {
			id: 2,
			arcs: map[int]int{
				0: 1,
				3: 1,
			},
		},
		3: {
			id: 3,
			arcs: map[int]int{
				1: 1,
				2: 1,
			},
		},
	}
	g := NewGraph(v)
	result := g.AllShortestPath(0, 3)
	/*
		result:=[[0,1,3],[0,2,3]]
	*/
	if len(result) != 2 {
		t.Error("shoude be two shortest path")
	}
}

func TestGraph_AllShortestPath2(t *testing.T) {
	v := map[int]Vertex{
		0: {
			id: 0,
			arcs: map[int]int{
				1: 1,
				2: 1,
			},
		},
		1: {
			id: 1,
			arcs: map[int]int{
				0: 1,
				3: 1,
			},
		},
		2: {
			id: 2,
			arcs: map[int]int{
				0: 1,
				3: 1,
				4: 1,
			},
		},
		3: {
			id: 3,
			arcs: map[int]int{
				1: 1,
				2: 1,
				4: 1,
			},
		},
		4: {
			id: 4,
			arcs: map[int]int{
				2: 1,
				3: 1,
			},
		},
	}
	g := NewGraph(v)
	result := g.AllShortestPath(0, 3)
	/*
		result:=[[0,1,3],[0,2,3]]
	*/
	if len(result) != 2 {
		t.Error("shoude be two shortest path")
	}
}
