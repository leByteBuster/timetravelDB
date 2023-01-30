package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func main() {

	graph_nodes, err := loadJsonData("graph_nodes.json")
	graph_edges, err := loadJsonData("graph_edges.json")

	fmt.Printf("\n Number nodes: %v", len(graph_nodes))
	fmt.Printf("\n Number edges: %v", len(graph_edges))

	nodeQuereis, timeSeriesNodes := getQuerStringsNodes(graph_nodes)
	edgeQueries, timeSeriesEdges := getQueryStringsEdges(graph_edges)

	queryMultipleNeo4j(context.Background(), "neo4j://localhost:7687", "neo4j", "rhebo", nodeQuereis)
	queryMultipleNeo4j(context.Background(), "neo4j://localhost:7687", "neo4j", "rhebo", edgeQueries)

	fmt.Printf("\n Number time-series map nodes: %v", len(timeSeriesNodes))
	fmt.Printf("\n Number time-series map edges: %v", len(timeSeriesEdges))

	timeSeries := map[uuid.UUID][]map[string]interface{}{}

	for k, v := range timeSeriesNodes {
		timeSeries[k] = v
	}
	for k, v := range timeSeriesEdges {
		timeSeries[k] = v
	}

	loadDataTimeScaleDB(timeSeries)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	// fmt.Printf("Result: %v", res)
}
