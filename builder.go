package astro

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type NodeType string

const (
	FuncDecl NodeType = "Function"
	Unknown  NodeType = "Unknown"
)

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

func ExtractGraphFrmAST(src string) (*Graph, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, fmt.Errorf("Error parsing file: %s", err)
	}

	graph := NewGraph()
	currentFunc := NewNode(Unknown, "")

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// create function node and add it to graph
			funcNode := NewNode(FuncDecl, x.Name.Name)
			graph.Nodes = append(graph.Nodes, funcNode)

			// set current function
			currentFunc = funcNode

		case *ast.CallExpr:
			err := processCall(x, graph, currentFunc)
			if err != nil {
				return false
			}

		default:
			// do nothing
		}

		return true
	})

	return graph, nil
}

func processCall(x *ast.CallExpr, graph *Graph, currentFunc *Node) error {
	var funcName string

	switch call := x.Fun.(type) {
	case *ast.Ident:
		funcName = call.Name
	case *ast.SelectorExpr:
		// handle method or package-level function calls
		if ident, ok := call.X.(*ast.Ident); ok {
			funcName = fmt.Sprintf("%s.%s", ident.Name, call.Sel.Name)
		} else {
			return fmt.Errorf("unhandled call type: %T", call.X)
		}
	default:
		return fmt.Errorf("unknown call type: %T", x.Fun)
	}

	callFunc, exists := graph.NodeMap[funcName]
	if !exists {
		callFunc = NewNode(FuncDecl, funcName)
		graph.AddNode(callFunc)
	}

	// Create an edge from the current function to the called function
	graph.AddEdge(currentFunc, callFunc, Call)

	return nil
}
