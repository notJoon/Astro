package astro

import (
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
			expectedNodeCount: 4,
			expectedEdgeCount: 2, // main -> printMore, printMore -> println
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			graph, err := ExtractGraphFromAST(tc.src)
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

func TestExtractGraphFromAST_Variable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		src               string
		expectedNodeCount int
		expectedEdgeCount int
	}{
		{
			name: "variable declaration and usage",
			src: `
package main
func main() {
	var x int
	x = 5
	println(x)
}`,
			expectedNodeCount: 3, // main, x, println
			expectedEdgeCount: 2, // main -> x (Declares), main -> x (Uses)
		},
		{
			name: "variable passed to function",
			src: `
package main
func main() {
	var message string = "hello"
	print(message)
}
func print(msg string) {
	println(msg)
}`,
			expectedNodeCount: 5, // main, message, print, msg, println
			expectedEdgeCount: 3, // main -> message (Declares), main -> print (Calls), print -> msg (PassesTo)
		},
		{
			name: "global variable usage",
			src: `
package main
var globalVar int
func main() {
	globalVar = 10
	println(globalVar)
}`,
			expectedNodeCount: 3, // globalVar, main, println
			expectedEdgeCount: 2, // main -> globalVar (Uses), main -> globalVar (Declares)
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			graph, err := ExtractGraphFromAST(tc.src)
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
