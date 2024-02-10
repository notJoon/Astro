package astro

import (
	"fmt"
	"strings"
)

type QueryBuilder interface {
	BuildQuery() string
	ApplyAbstract() string
}

type ConcreteQueryBuilder struct {
	Graph *Graph
}

func (q *ConcreteQueryBuilder) BuildQuery() string {
	l := len(q.Graph.Edges)

	if q.Graph == nil || l == 0 {
		return "Empty graph"
	}

	var builder strings.Builder
	for i, edge := range q.Graph.Edges {
		builder.WriteString(edge.String())
		if i < l-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func (q *ConcreteQueryBuilder) ApplyAbstract() string {
	l := len(q.Graph.Edges)

	if q.Graph == nil || l == 0 {
		return "Empty graph"
	}

	var builder strings.Builder
	for i, edge := range q.Graph.Edges {
		from := fmt.Sprintf("(%s)", edge.From.Type)
		to := fmt.Sprintf("(%s)", edge.To.Type)
		relation := fmt.Sprintf("-[:%s]->", edge.Relation)

		builder.WriteString(from + relation + to)
		if i < l-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

type Query struct {
	Builder QueryBuilder
}

func NewQuery(builder QueryBuilder) *Query {
	return &Query{
		Builder: builder,
	}
}
