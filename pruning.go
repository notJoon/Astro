package astro

// MultiPathPruning performs a search on the given graph starting from startNode to find a path that satisfies the goal function.
//
// It utilizes a multi-path pruning strategy to avoid revisiting nodes that have already been explored in different paths.
//
// The function returns the first path that satisfies the goal function, along with a boolean indicating whether such a path was found.
//
// ref: https://artint.info/3e/html/ArtInt3e.Ch3.S7.html $3.7.2
func MultiPathPruning(graph *Graph, startNode *Node, goal func(n *Node) bool) ([]*Node, bool) {
	frontier, explored, visited := initialize(startNode)

	for len(frontier) > 0 {
		path, node := exploreNode(&frontier, explored)
		
		if goal(node) {
			return path, true
		}

		updateFrontier(graph, node, path, &frontier, visited)
	}

	return nil, false
}

// initialize performs the initial setup to start the search.
func initialize(start *Node) ([][]*Node, map[*Node]bool, map[*Node][]*Node) {
	frontier := [][]*Node{{start}}
	explored := make(map[*Node]bool)
	visited := make(map[*Node][]*Node)

	return frontier, explored, visited
}

// exploreNode navigate the current node, and verify that it has found the target node (gola).
func exploreNode(frontier *[][]*Node, visited map[*Node]bool) ([]*Node, *Node) {
	path := (*frontier)[0]
	*frontier = (*frontier)[1:]

	node := path[len(path)-1]
	visited[node] = true

	return path, node
}

func updateFrontier(graph *Graph, node *Node, path []*Node, frontier *[][]*Node, visited map[*Node][]*Node) {
	for _, edge := range graph.Edges {
		if edge.From == node && !contains(path, edge.To) {
			newPath := append([]*Node(nil), path...)
			newPath = append(newPath, edge.To)

			if visitedPath, found := visited[edge.To]; found {
				newPath = append(newPath, visitedPath...)
			} else {
				visited[edge.To] = newPath
			}

			*frontier = append(*frontier, newPath)
		}
	}
}

func contains(path []*Node, node *Node) bool {
	for _, p := range path {
		if p == node {
			return true
		}
	}

	return false
}
