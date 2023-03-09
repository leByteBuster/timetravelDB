package dataadapter

import (
	"context"
	"fmt"
	"log"

	databaseapi "github.com/LexaTRex/timetravelDB/database-api"
	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/google/uuid"
)

func LoadData() {

	graph_nodes, err := utils.LoadJsonData("data-adapter/graph_nodes.json")
	if err != nil {
		log.Printf("Error loading nodes from json: %v", err)
	}
	graph_edges, err := utils.LoadJsonData("data-adapter/graph_edges.json")
	if err != nil {
		log.Printf("Error loading edges from json: %v", err)
	}
	//log.Printf("\n Number nodes: %v", len(graph_nodes))
	//log.Printf("\n Number edges: %v", len(graph_edges))

	nodeQuereis, timeSeriesNodes := getQueryStringsNodes(graph_nodes)
	edgeQueries, timeSeriesEdges := getQueryStringsEdges(graph_edges)

	databaseapi.WriteQueryMultipleNeo4j(context.Background(), nodeQuereis)
	databaseapi.WriteQueryMultipleNeo4j(context.Background(), edgeQueries)

	// log.Printf("\n Number time-series map nodes: %v", len(timeSeriesNodes))
	// log.Printf("\n Number time-series map edges: %v", len(timeSeriesEdges))

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
