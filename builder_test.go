package astro_test

import (
	"astro"
	"testing"
)

func TestExtractGraphFromAST(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		src               string
		expectedNodeCount int
		expectedEdgeCount int
	}{
		{
			name: "only main function",
			src: `
package main
func main() {
}`,
			expectedNodeCount: 1, // main
			expectedEdgeCount: 0,
		},
		{
			name: "main and println",
			src: `
package main
import "fmt"
func main() {
	fmt.Println("Hello, World!")
}`,
			expectedNodeCount: 2, // main, fmt.Println
			expectedEdgeCount: 1, // main -> fmt.Println
		},
		{
			name: "simple function",
			src: `
package main
import "fmt"
func main() {
	fmt.Println("Hello, World!")
}`,
			expectedNodeCount: 2, // main, fmt.Println
			expectedEdgeCount: 1, // main -> fmt.Println
		},
		{
			name: "two functions",
			src: `
package main

func main() {
	printMore("Hello")
}
func printMore(msg string) {
	println(msg)
}
`,
			expectedNodeCount: 4, // main, printMore, println
			expectedEdgeCount: 2, // main -> printMore, printMore -> println
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			graph, err := astro.ExtractGraphFrmAST(tc.src)
			if err != nil {
				t.Fatalf("Error extracting graph: %s", err)
			}
			if len(graph.Nodes) != tc.expectedNodeCount {
				t.Errorf("Expected %d nodes, got %d", tc.expectedNodeCount, len(graph.Nodes))
			}
			if len(graph.Edges) != tc.expectedEdgeCount {
				t.Errorf("Expected %d edges, got %d", tc.expectedEdgeCount, len(graph.Edges))
			}
		})
	}
}

func TestConvertASTtoGraphQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		graph          *astro.Graph
		expectedOutput string
	}{
		{
			name: "simple function call",
			graph: func() *astro.Graph {
				g := astro.NewGraph()
				mainNode := astro.NewNode(astro.FuncDecl, "main")
				printNode := astro.NewNode(astro.FuncDecl, "fmt.Println")
				g.AddNode(mainNode)
				g.AddNode(printNode)
				g.AddEdge(mainNode, printNode, astro.Call)
				return g
			}(),
			expectedOutput: "(main)-[:Call]->(fmt.Println)\n",
		},
		{
			name: "two function calls",
			graph: func() *astro.Graph {
				g := astro.NewGraph()

				mainNode := astro.NewNode(astro.FuncDecl, "main")
				printNode := astro.NewNode(astro.FuncDecl, "fmt.Println")
				moreNode := astro.NewNode(astro.FuncDecl, "printMore")

				g.AddNode(mainNode)
				g.AddNode(printNode)
				g.AddNode(moreNode)
				g.AddEdge(mainNode, moreNode, astro.Call)
				g.AddEdge(moreNode, printNode, astro.Call)

				return g
			}(),
			expectedOutput: "(main)-[:Call]->(printMore)\n(printMore)-[:Call]->(fmt.Println)\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := tc.graph.String()
			if actualOutput != tc.expectedOutput {
				t.Errorf("Expected output:\n%s\nGot:\n%s", tc.expectedOutput, actualOutput)
			}
		})
	}
}