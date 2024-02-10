package astro

import (
	"fmt"
	"reflect"
	"sort"
)

// NodeType represents the type of a AST node in the graph
type NodeType string

const (
	Func    NodeType = "Function"
	Var     NodeType = "Variable"
	Unknown NodeType = "Unknown"
)

// Node holds the information of a AST node in the graph.
//
// For example, a function declaration node will have the type "FuncDecl"
// and the function name (or identifier) as the name.
type Node struct {
	Type NodeType
	Name string
}

func NewNode(t NodeType, name string) *Node {
	return &Node{
		Type: t,
		Name: name,
	}
}

func (n *Node) SetType(t NodeType) {
	n.Type = t
}

func (n *Node) SetName(name string) {
	n.Name = name
}

func (n *Node) String() string {
	return fmt.Sprintf("(%s)", n.Name)
}

// Relation represents the type of directional relationship between two nodes in the graph.
//
// it can be a function call, a variable assignment, etc.
//
// For example, a function call relation will have the type "Call"
type Relation string

const (
	Call            Relation = "Call"     // function call
	Declares        Relation = "Declares" // variable declaration
	Uses            Relation = "Uses"     // variable usage
	PassesTo        Relation = "PassesTo" // variables passed as function parameters
	UnknownRelation Relation = "Unknown"
)

type Edge struct {
	From     *Node
	To       *Node
	Relation Relation
}

func NewEdge(from *Node, to *Node, r Relation) *Edge {
	return &Edge{
		From:     from,
		To:       to,
		Relation: r,
	}
}

func (e *Edge) String() string {
	return fmt.Sprintf("%s-[:%s]->%s", e.From, e.Relation, e.To)
}

// Graph represents the graph of the AST nodes and their relationships.
type Graph struct {
	Nodes   []*Node
	Edges   []*Edge
	NodeMap map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		NodeMap: make(map[string]*Node),
	}
}

func (g *Graph) AddNode(node *Node) {
	if _, exists := g.NodeMap[node.Name]; !exists {
		g.Nodes = append(g.Nodes, node)
		g.NodeMap[node.Name] = node
	}
}

func (g *Graph) AddEdge(from, to *Node, relation Relation) {
	edge := NewEdge(from, to, relation)
	g.Edges = append(g.Edges, edge)
}

// DegreeSequence computes and returns a sorted slice of degree sequence of a given graph.
//
// The degree of a node is defined as the number of edges connected to it. The degree sequence
// is a way to summarize the connectivity pattern of the graph.
//
// It can be used among other things, for graph comaprison, where two graphs with the same
// degree sequence might have similar structures.
func (g *Graph) DegreeSequence() []int {
	dmap := make(map[string]int)

	for _, edge := range g.Edges {
		dmap[edge.From.Name]++
		dmap[edge.To.Name]++
	}

	sequence := make([]int, 0, len(dmap))
	for _, degree := range dmap {
		sequence = append(sequence, degree)
	}

	sort.Ints(sequence)

	return sequence
}

// IsIsomorphic checks if two graphs, g1 and g2, are isomorphic.
//
// Two graphs are considered isomorphic if they contain the same number of nodes and edges,
// and if the nodes can be matched one-to-one based on their NodeType, such that the edges
// also have a one-to-one correspondence in their Relation types between the graph .
//
// This function does not consider the names of the nodes but only their types and the relationships (edges).
//
// TODO: should reduce complexity
func IsIsomorphic(g1, g2 *Graph) bool {
	isSameLen := len(g1.Nodes) != len(g2.Nodes)
	isSameEdgeLen := len(g1.Edges) != len(g2.Edges)

	if isSameLen || isSameEdgeLen {
		return false
	}

	nodeTypeCount1 := make(map[NodeType]int)
	nodeTypeCount2 := make(map[NodeType]int)
	for _, node := range g1.Nodes {
		nodeTypeCount1[node.Type]++
	}
	for _, node := range g2.Nodes {
		nodeTypeCount2[node.Type]++
	}
	if !reflect.DeepEqual(nodeTypeCount1, nodeTypeCount2) {
		return false
	}

	edgeRelationCount1 := make(map[Relation]int)
	edgeRelationCount2 := make(map[Relation]int)
	for _, edge := range g1.Edges {
		edgeRelationCount1[edge.Relation]++
	}
	for _, edge := range g2.Edges {
		edgeRelationCount2[edge.Relation]++
	}
	return reflect.DeepEqual(edgeRelationCount1, edgeRelationCount2)
}
