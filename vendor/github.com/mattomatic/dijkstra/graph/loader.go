package graph

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func LoadGraph(filename string) *Graph {
	reader := getReader(filename)
	lines := getLines(reader)

	g := NewGraph()

	for line := range lines {
		fields := strings.Split(line, " ")
		headId := atoi(fields[0])
		edges := fields[1:]

		getOrCreate(g, headId)

		for _, pair := range edges {
			if len(strings.Split(pair, ",")) == 1 {
				continue
			}

			tailId := atoi(strings.Split(pair, ",")[0])
			weight := atoi(strings.Split(pair, ",")[1])

			head := getOrCreate(g, headId)
			tail := getOrCreate(g, tailId)
			edge := NewEdge(head, tail, weight)
			head.AddEdges(edge)
		}
	}

	return g
}

func atoi(s string) int {
	result, err := strconv.Atoi(s)

	if err != nil {
		panic("failed to convert string")
	}

	return result
}

func getOrCreate(g *Graph, id int) *Node {
	node := g.Search(id)

	if node == nil {
		node = NewNode(id)
		g.AddNodes(node)
	}

	return node
}

func getReader(filename string) *bufio.Reader {
	fp, err := os.Open(filename)

	if err != nil {
		panic("could not open file")
	}

	return bufio.NewReader(fp)
}

func getLines(reader *bufio.Reader) chan string {
	ch := make(chan string, 0)

	go func() {
		defer close(ch)

		for {
			line, err := reader.ReadString('\n')

			if err != nil {
				break
			}

			if strings.HasPrefix(line, "#") {
				continue
			}

			ch <- strings.Trim(line, "\n")
		}
	}()

	return ch
}
