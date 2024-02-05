package astro_test

import (
	"astro"
	"testing"
)

func TestExtractGraphFromAST(t *testing.T) {
	src := `
package main
import "fmt"
func main() {
    fmt.Println("Hello, World!")
    printMore("More messages")
}
func printMore(msg string) {
    fmt.Println(msg)
}
`

	graph, err := astro.ExtractGraphFrmAST(src)
	if err != nil {
		t.Errorf("Error extracting graph from AST: %s", err)
	}

	expectedNodeCount := 3 // main, printMore, fmt.Println
	if len(graph.Nodes) != expectedNodeCount {
		t.Errorf("Expected %d nodes, got %d", expectedNodeCount, len(graph.Nodes))
	}

	expectedEdgeCount := 2 // main -> fmt.Println, main -> printMore
	if len(graph.Edges) != expectedEdgeCount {
		t.Errorf("Expected %d edges, got %d", expectedEdgeCount, len(graph.Edges))
	}
}
