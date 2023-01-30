package main

import (
	"context"
	"fmt"
)

func main() {

	graph_nodes, err := loadJsonData("graph_nodes.json")
	graph_edges, err := loadJsonData("graph_edges.json")

	fmt.Printf("\n Number nodes: %v", len(graph_nodes))
	fmt.Printf("\n Number edges: %v", len(graph_edges))

	nodeQuereis, timeSeriesNodes := getQuerStringsNodes(graph_nodes, context.Background(), "neo4j://localhost:7687", "neo4j", "rhebo")
	edgeQueries, timeSeriesEdges := getQueryStringsEdges(graph_edges)

	queryMultipleNeo4j(context.Background(), "neo4j://localhost:7687", "neo4j", "rhebo", nodeQuereis)
	queryMultipleNeo4j(context.Background(), "neo4j://localhost:7687", "neo4j", "rhebo", edgeQueries)

	fmt.Printf("\n Number time-series map nodes: %v", len(timeSeriesNodes))
	fmt.Printf("\n Number time-series map edges: %v", len(timeSeriesEdges))

	loadDataTimeScaleDB(timeSeriesNodes, timeSeriesEdges)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	// fmt.Printf("Result: %v", res)
}
