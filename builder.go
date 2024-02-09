package astro

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// ExtractGraphFromAST extracts a graph from the given source code (go file).
//
// It utilizes the go/ast package to parse the source code and build the graph.
//
// Consideration: use a treesitter parser instead of go/ast to support more languages than just Go
func ExtractGraphFromAST(src string) (*Graph, error) {
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
			funcNode := NewNode(Func, x.Name.Name)
			graph.Nodes = append(graph.Nodes, funcNode)

			// set current function
			currentFunc = funcNode

		case *ast.GenDecl:
			// variable declaration
			if x.Tok == token.VAR {
				for _, spec := range x.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range vs.Names {
							varNode := NewNode(Var, name.Name)
							graph.Nodes = append(graph.Nodes, varNode)

							graph.AddEdge(currentFunc, varNode, Declares)
						}
					}
				}
			}

		case *ast.Ident:
			if currentFunc != nil && x.Obj != nil && x.Obj.Kind == ast.Var {
				// check for variable usage in the current function
				varNode, exists := graph.NodeMap[x.Name]
				if exists {
					graph.AddEdge(currentFunc, varNode, Uses)
				}
			}

		case *ast.CallExpr:
			if err := processCall(x, graph, currentFunc); err != nil {
				return false
			}

			// handling variables passed as parameters in function calls
			for _, arg := range x.Args {
				if ident, ok := arg.(*ast.Ident); ok && ident.Obj != nil && ident.Obj.Kind == ast.Var {
					varNode, exists := graph.NodeMap[ident.Name]
					if exists {
						graph.AddEdge(varNode, currentFunc, PassesTo)
					}
				}
			}

		default:
			// do nothing
		}

		return true
	})

	return graph, nil
}

// processCall extracts node information from the generated AST and
// and converts it into a graph structure.
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
		callFunc = NewNode(Func, funcName)
		graph.AddNode(callFunc)
	}

	// Create an edge from the current function to the called function
	graph.AddEdge(currentFunc, callFunc, Call)

	return nil
}
