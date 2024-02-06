package astro

import (
	"reflect"
	"testing"
)

/* This test case constructs a simple graph and tests the SearchMPP function to find a path
 * from node A to node D.
 *
 * The goal function defines that the goal is reached when node D is found.
 * If the function works correctly, it should return the path A -> B -> C -> D, and the test should pass.
 */
func TestSearchMPP(t *testing.T) {
	t.Parallel()

	nodeA := &Node{Type: "TypeA", Name: "A"}
	nodeB := &Node{Type: "TypeB", Name: "B"}
	nodeC := &Node{Type: "TypeC", Name: "C"}
	nodeD := &Node{Type: "TypeD", Name: "D"} // Goal node for the test

	graph := &Graph{
		Nodes: []*Node{nodeA, nodeB, nodeC, nodeD},
		Edges: []*Edge{
			{From: nodeA, To: nodeB, Relation: Call},
			{From: nodeB, To: nodeC, Relation: Call},
			{From: nodeC, To: nodeD, Relation: Call},
		},
		NodeMap: map[string]*Node{
			"A": nodeA,
			"B": nodeB,
			"C": nodeC,
			"D": nodeD,
		},
	}

	// Define goal function
	goal := func(n *Node) bool {
		return n.Name == "D"
	}

	path, found := MultiPathPruning(graph, nodeA, goal)
	if !found {
		t.Fatalf("Goal node D was not found")
	}

	// Validate the path
	expectedPath := []*Node{nodeA, nodeB, nodeC, nodeD}
	if len(path) != len(expectedPath) {
		t.Fatalf("Expected path length of %d, got %d", len(expectedPath), len(path))
	}

	for i, node := range path {
		if node != expectedPath[i] {
			t.Errorf("Expected node %v at position %d, got %v", expectedPath[i], i, node)
		}
	}
}

func TestSearchMPP2(t *testing.T) {
	t.Parallel()

	nodeA := &Node{Type: "TypeA", Name: "A"}
	nodeB := &Node{Type: "TypeB", Name: "B"}
	nodeC := &Node{Type: "TypeC", Name: "C"}
	nodeD := &Node{Type: "TypeD", Name: "D"}
	nodeE := &Node{Type: "TypeE", Name: "E"}

	// Construct graph
	graph := &Graph{
		Nodes: []*Node{nodeA, nodeB, nodeC, nodeD, nodeE},
		Edges: []*Edge{
			{From: nodeA, To: nodeB, Relation: Call},
			{From: nodeB, To: nodeC, Relation: Call},
			{From: nodeC, To: nodeD, Relation: Call},
			{From: nodeA, To: nodeE, Relation: Call},
			{From: nodeE, To: nodeD, Relation: Call},
		},
		NodeMap: map[string]*Node{
			"A": nodeA,
			"B": nodeB,
			"C": nodeC,
			"D": nodeD,
			"E": nodeE,
		},
	}

	tests := []struct {
		name      string
		start     *Node
		goalName  string
		wantPath  []*Node
		wantFound bool
	}{
		{
			name:      "Path A to D through E",
			start:     nodeA,
			goalName:  "D",
			wantPath:  []*Node{nodeA, nodeE, nodeD},
			wantFound: true,
		},
		{
			name:      "Path A to B",
			start:     nodeA,
			goalName:  "B",
			wantPath:  []*Node{nodeA, nodeB},
			wantFound: true,
		},
		{
			name:      "No Path from D to A",
			start:     nodeD,
			goalName:  "A",
			wantPath:  nil,
			wantFound: false,
		},
		{
			name:      "Path B to D",
			start:     nodeB,
			goalName:  "D",
			wantPath:  []*Node{nodeB, nodeC, nodeD},
			wantFound: true,
		},
		{
			name:      "No Path from E to B",
			start:     nodeE,
			goalName:  "B",
			wantPath:  nil,
			wantFound: false,
		},
		{
			name:      "Path C to D",
			start:     nodeC,
			goalName:  "D",
			wantPath:  []*Node{nodeC, nodeD},
			wantFound: true,
		},
		{
			name:      "No Path from D to B",
			start:     nodeD,
			goalName:  "B",
			wantPath:  nil,
			wantFound: false,
		},
		{
			name:      "Path A to E",
			start:     nodeA,
			goalName:  "E",
			wantPath:  []*Node{nodeA, nodeE},
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goal := func(n *Node) bool {
				return n.Name == tt.goalName
			}
			gotPath, gotFound := MultiPathPruning(graph, tt.start, goal)

			if gotFound != tt.wantFound {
				t.Errorf("SearchMPP(%s) gotFound = %v, want %v", tt.name, gotFound, tt.wantFound)
			}
			if gotFound && !reflect.DeepEqual(gotPath, tt.wantPath) {
				t.Errorf("SearchMPP(%s) gotPath = %v, want %v", tt.name, gotPath, tt.wantPath)
			}
		})
	}
}

func TestDetectCycle(t *testing.T) {
	t.Parallel()

	nodeF := &Node{Type: "TypeF", Name: "F"}
	nodeG := &Node{Type: "TypeG", Name: "G"}
	nodeH := &Node{Type: "TypeH", Name: "H"}

	// Construct graph with a cycle
	graph := &Graph{
		Nodes: []*Node{nodeF, nodeG, nodeH},
		Edges: []*Edge{
			{From: nodeF, To: nodeG, Relation: Call},
			{From: nodeG, To: nodeH, Relation: Call},
			{From: nodeH, To: nodeF, Relation: Call},
		},
		NodeMap: map[string]*Node{
			"F": nodeF,
			"G": nodeG,
			"H": nodeH,
		},
	}

	// Define goal function to stop at node H
	goal := func(n *Node) bool {
		return n.Name == "H"
	}

	// Perform the search starting at node F
	path, found := MultiPathPruning(graph, nodeF, goal)
	if !found {
		t.Errorf("Goal node H was not found")
	}

	// Expected path: F -> G -> H
	expectedPath := []*Node{nodeF, nodeG, nodeH}
	if len(path) != len(expectedPath) {
		t.Errorf("Expected path length of %d, got %d", len(expectedPath), len(path))
	}

	for i, node := range path {
		if node != expectedPath[i] {
			t.Errorf("Expected node %v at position %d, got %v", expectedPath[i], i, node)
		}
	}
}
