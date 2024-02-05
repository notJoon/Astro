package astro

import (
	"errors"
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

type Relation string

const (
	Call            Relation = "Call"
	UnknownRelation Relation = "Unknown"
)

type Edge struct {
	From     *Node
	to       *Node
	Relation Relation
}

func NewEdge(from *Node, to *Node, r Relation) *Edge {
	return &Edge{
		From:     from,
		to:       to,
		Relation: r,
	}
}

type Graph struct {
	Nodes []*Node
	Edges []*Edge
}

func NewGraph() *Graph {
	return &Graph{}
}

func ExtractGraphFrmAST(src string) (*Graph, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing file: %s", err))
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
        funcName = fmt.Sprintf("%s.%s", call.X, call.Sel)

        // find the node of the called function
        var callFunc *Node
        for _, node := range graph.Nodes {
            if node.Name == strings.TrimPrefix(funcName, "*") && node.Type == FuncDecl {
                callFunc = node
                break
            }
        }

        // create edge and add it to graph
        if callFunc != nil {
            edge := NewEdge(currentFunc, callFunc, Call)
            graph.Edges = append(graph.Edges, edge)
        }
	default:
		return errors.New("Unknown call type")
    }

	return nil
}