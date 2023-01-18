package main

import (
	"context"
	"fmt"
)

func main() {

	graph_nodes, err := loadJsonData("graph_nodes.json")
	graph_edges, err := loadJsonData("graph_edges.json")

	timeSeriesMapNodes := loadGraphNodesIntoNeo4jDatabase(graph_nodes, context.Background(), "neo4j://localhost:7687", "neo4j", "rhebo")
	timeSeriesMapEdges := loadGraphEdgesIntoNeo4jDatabase(graph_edges, context.Background(), "neo4j://localhost:7687", "neo4j", "rhebo")

	loadDataTimeScaleDB(timeSeriesMapNodes, timeSeriesMapEdges)

	fmt.Printf("\n Node time-series: %v", timeSeriesMapNodes)
	fmt.Printf("\n Edge time-series: %v", timeSeriesMapEdges)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	// fmt.Printf("Result: %v", res)
}
