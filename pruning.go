package astro

// ref: https://artint.info/3e/html/ArtInt3e.Ch3.S7.html $3.7.2
func MultiPathPruning(graph *Graph, startNode *Node, goal func(n *Node) bool) ([]*Node, bool) {
	frontier := [][]*Node{{startNode}}
	explored := make(map[*Node]bool)
	visited := make(map[*Node][]*Node)

	for len(frontier) > 0 {
		path := frontier[0]
		frontier = frontier[1:]
		node := path[len(path)-1]

		if !explored[node] {
			explored[node] = true
			if goal(node) {
				return path, true
			}

			for _, edge := range graph.Edges {
				if edge.From == node {
					if contains(path, edge.To) {
						continue
					}

					if visitPath, found := visited[edge.To]; found {
						newPath := append(path, visitPath...)
						frontier = append(frontier, newPath)
					} else {
						newPath := make([]*Node, len(path)+1)
						copy(newPath, path)
						newPath[len(newPath)-1] = edge.To
						frontier = append(frontier, newPath)
					}
				}
			}
		}
	}

	return nil, false
}

func contains(path []*Node, node *Node) bool {
	for _, n := range path {
		if n == node {
			return true
		}
	}

	return false
}
