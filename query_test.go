package astro

import (
	"testing"
)

func TestBuildGraphQuery(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name              string
        setupGraph        func() *Graph // Function to setup the graph
        expectedConcrete  string        // Expected output for concrete query
        expectedAbstract  string        // Expected output for abstract query
    }{
        {
            name: "simple function call",
            setupGraph: func() *Graph {
                g := NewGraph()
                mainNode := NewNode(Func, "main")
                printNode := NewNode(Func, "fmt.Println")
                g.AddNode(mainNode)
                g.AddNode(printNode)
                g.AddEdge(mainNode, printNode, Call)
                return g
            },
            expectedConcrete: "(main)-[:Call]->(fmt.Println)",
            expectedAbstract: "(Function)-[:Call]->(Function)",
        },
        {
            name: "two function calls",
            setupGraph: func() *Graph {
                g := NewGraph()

                mainNode := NewNode(Func, "main")
                printNode := NewNode(Func, "fmt.Println")
                moreNode := NewNode(Func, "printMore")

                g.AddNode(mainNode)
                g.AddNode(printNode)
                g.AddNode(moreNode)
                g.AddEdge(mainNode, moreNode, Call)
                g.AddEdge(moreNode, printNode, Call)

                return g
            },
            expectedConcrete: "(main)-[:Call]->(printMore)\n(printMore)-[:Call]->(fmt.Println)",
            expectedAbstract: "(Function)-[:Call]->(Function)\n(Function)-[:Call]->(Function)",
        },
		{
            name: "variable declaration and usage",
            setupGraph: func() *Graph {
                g := NewGraph()

                mainNode := NewNode(Func, "main")
                varNode := NewNode(Var, "x")
                printlnNode := NewNode(Func, "println")

                g.AddNode(mainNode)
                g.AddNode(varNode)
                g.AddNode(printlnNode)
                g.AddEdge(mainNode, varNode, Declares)
                g.AddEdge(mainNode, varNode, Uses)
                g.AddEdge(mainNode, printlnNode, Call)

                return g
            },
            expectedConcrete: "(main)-[:Declares]->(x)\n(main)-[:Uses]->(x)\n(main)-[:Call]->(println)",
            expectedAbstract: "(Function)-[:Declares]->(Variable)\n(Function)-[:Uses]->(Variable)\n(Function)-[:Call]->(Function)",
        },
        {
            name: "variable passed to function",
            setupGraph: func() *Graph {
                g := NewGraph()

                mainNode := NewNode(Func, "main")
                varNode := NewNode(Var, "message")
                printNode := NewNode(Func, "print")
                printlnNode := NewNode(Func, "println")

                g.AddNode(mainNode)
                g.AddNode(varNode)
                g.AddNode(printNode)
                g.AddNode(printlnNode)
                g.AddEdge(mainNode, varNode, Declares)
                g.AddEdge(mainNode, printNode, Call)
                g.AddEdge(printNode, printlnNode, Call)
                g.AddEdge(mainNode, printNode, PassesTo)

                return g
            },
            expectedConcrete: "(main)-[:Declares]->(message)\n(main)-[:Call]->(print)\n(print)-[:Call]->(println)\n(main)-[:PassesTo]->(print)",
            expectedAbstract: "(Function)-[:Declares]->(Variable)\n(Function)-[:Call]->(Function)\n(Function)-[:Call]->(Function)\n(Function)-[:PassesTo]->(Function)",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            graph := tc.setupGraph()
            queryBuilder := &ConcreteQueryBuilder{Graph: graph}

            // Test concrete query building
            actualConcrete := queryBuilder.BuildQuery()
            if actualConcrete != tc.expectedConcrete {
                t.Errorf("Concrete Query - Expected output:\n%s\nGot:\n%s", tc.expectedConcrete, actualConcrete)
            }

            // Test abstract query building
            actualAbstract := queryBuilder.ApplyAbstract()
            if actualAbstract != tc.expectedAbstract {
                t.Errorf("Abstract Query - Expected output:\n%s\nGot:\n%s", tc.expectedAbstract, actualAbstract)
            }
        })
    }
}
