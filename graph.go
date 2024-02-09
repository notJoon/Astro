package astro

import (
	"fmt"
	"strings"
)

// NodeType represents the type of a AST node in the graph
type NodeType string

const (
	FuncDecl NodeType = "Function"
	Unknown  NodeType = "Unknown"
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
	Call            Relation = "Call"
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

func (g *Graph) String() string {
	var builder strings.Builder

	for _, edge := range g.Edges {
		builder.WriteString(edge.String() + "\n")
	}

	return builder.String()
}
