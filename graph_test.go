package astro

import (
	"reflect"
	"testing"
)

func TestGraph_DegreeSequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		setup          func(*Graph)
		expectedDegree []int
	}{
		{
			description:    "Empty graph",
			setup:          func(g *Graph) {},
			expectedDegree: []int{},
		},
		{
			description: "Graph with two nodes and one edge",
			setup: func(g *Graph) {
				node1 := NewNode(Func, "function1")
				node2 := NewNode(Var, "variable1")
				g.AddNode(node1)
				g.AddNode(node2)
				g.AddEdge(node1, node2, Call)
			},
			expectedDegree: []int{1, 1},
		},
		{
			description: "Graph with three nodes and three edges",
			setup: func(g *Graph) {
				node1 := NewNode(Func, "function1")
				node2 := NewNode(Var, "variable1")
				node3 := NewNode(Func, "function2")
				g.AddNode(node1)
				g.AddNode(node2)
				g.AddNode(node3)
				g.AddEdge(node1, node2, Call)
				g.AddEdge(node2, node3, Uses)
				g.AddEdge(node3, node1, Declares)
			},
			expectedDegree: []int{2, 2, 2},
		},
	}

	for _, test := range tests {
		g := NewGraph()
		test.setup(g)

		degreeSequence := g.DegreeSequence()
		if !reflect.DeepEqual(degreeSequence, test.expectedDegree) {
			t.Errorf("%s: expected degree sequence %v, got %v", test.description, test.expectedDegree, degreeSequence)
		}
	}
}

func TestIsIsomorphic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setupG1 func(*Graph)
		setupG2 func(*Graph)
		ok      bool
	}{
		{
			name: "Non-isomorphic graphs with different relation types 1",
			setupG1: func(g *Graph) {
				g.AddNode(NewNode(Func, "function1"))
				g.AddNode(NewNode(Var, "variable2"))
				g.AddEdge(NewNode(Func, "function1"), NewNode(Var, "variable2"), Call)
			},
			setupG2: func(g *Graph) {
				g.AddNode(NewNode(Var, "variable1"))
				g.AddNode(NewNode(Var, "variable2"))
				g.AddEdge(NewNode(Var, "variable1"), NewNode(Var, "variable2"), Call)
			},
			ok: false,
		},
		{
			name: "Non-isomorphic graphs with different relation types 2",
			setupG1: func(g *Graph) {
				g.AddNode(NewNode(Func, "function1"))
				g.AddNode(NewNode(Var, "variable2"))
				g.AddEdge(NewNode(Func, "function1"), NewNode(Var, "variable2"), Call)
			},
			setupG2: func(g *Graph) {
				g.AddNode(NewNode(Func, "function1"))
				g.AddNode(NewNode(Var, "variable2"))
				g.AddEdge(NewNode(Func, "function1"), NewNode(Var, "variable2"), Uses)
			},
			ok: false,
		},
		{
			name: "Non-isomorphic graphs with different edge types",
			setupG1: func(g *Graph) {
				g.AddNode(NewNode(Func, "function1"))
				g.AddNode(NewNode(Var, "variable1"))
				g.AddEdge(NewNode(Func, "function1"), NewNode(Var, "variable1"), Call)
			},
			setupG2: func(g *Graph) {
				g.AddNode(NewNode(Func, "function2"))
				g.AddNode(NewNode(Var, "variable2"))
				g.AddEdge(NewNode(Func, "function2"), NewNode(Var, "variable2"), Uses)
			},
			ok: false,
		},
		{
			name: "Isomorphic graphs with same node types and edges",
			setupG1: func(g *Graph) {
				g.AddNode(NewNode(Func, "function1"))
				g.AddNode(NewNode(Var, "variable1"))
				g.AddEdge(NewNode(Func, "function1"), NewNode(Var, "variable1"), Call)
			},
			setupG2: func(g *Graph) {
				g.AddNode(NewNode(Func, "function2"))
				g.AddNode(NewNode(Var, "variable2"))
				g.AddEdge(NewNode(Func, "function2"), NewNode(Var, "variable2"), Call)
			},
			ok: true,
		},
		{
			name: "Non-isomorphic graphs with same node types but different edges",
			setupG1: func(g *Graph) {
				g.AddNode(NewNode(Func, "function1"))
				g.AddNode(NewNode(Var, "variable1"))
				g.AddEdge(NewNode(Func, "function1"), NewNode(Var, "variable1"), Call)
			},
			setupG2: func(g *Graph) {
				g.AddNode(NewNode(Func, "function2"))
				g.AddNode(NewNode(Var, "variable2"))
				g.AddEdge(NewNode(Func, "function2"), NewNode(Var, "variable2"), Uses)
			},
			ok: false,
		},
		{
			name: "Isomorphic graphs with multiple nodes and edges",
			setupG1: func(g *Graph) {
				g.AddNode(NewNode(Func, "function1"))
				g.AddNode(NewNode(Var, "variable1"))
				g.AddNode(NewNode(Var, "variable2"))
				g.AddEdge(NewNode(Func, "function1"), NewNode(Var, "variable1"), Call)
				g.AddEdge(NewNode(Func, "function1"), NewNode(Var, "variable2"), Uses)
			},
			setupG2: func(g *Graph) {
				g.AddNode(NewNode(Func, "function2"))
				g.AddNode(NewNode(Var, "variable3"))
				g.AddNode(NewNode(Var, "variable4"))
				g.AddEdge(NewNode(Func, "function2"), NewNode(Var, "variable3"), Call)
				g.AddEdge(NewNode(Func, "function2"), NewNode(Var, "variable4"), Uses)
			},
			ok: true,
		},
		{
			name: "Non-isomorphic graphs with different number of nodes",
			setupG1: func(g *Graph) {
				g.AddNode(NewNode(Func, "function1"))
				g.AddNode(NewNode(Var, "variable1"))
				g.AddEdge(NewNode(Func, "function1"), NewNode(Var, "variable1"), Call)
			},
			setupG2: func(g *Graph) {
				g.AddNode(NewNode(Func, "function2"))
				g.AddEdge(NewNode(Func, "function2"), NewNode(Var, "variable2"), Call)
				g.AddEdge(NewNode(Func, "function2"), NewNode(Var, "variable3"), Uses)
			},
			ok: false,
		},
		{
			name: "Non-isomorphic graphs with different number of edges",
			setupG1: func(g *Graph) {
				g.AddNode(NewNode(Func, "function1"))
				g.AddNode(NewNode(Var, "variable1"))
				g.AddEdge(NewNode(Func, "function1"), NewNode(Var, "variable1"), Call)
			},
			setupG2: func(g *Graph) {
				g.AddNode(NewNode(Func, "function2"))
				g.AddNode(NewNode(Var, "variable2"))
			},
			ok: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g1 := NewGraph()
			tc.setupG1(g1)

			g2 := NewGraph()
			tc.setupG2(g2)

			if IsIsomorphic(g1, g2) != tc.ok {
				t.Errorf("IsIsomorphic()=%s, expected isomorphic=%v, got %v", tc.name, tc.ok, !tc.ok)
			}
		})
	}
}
