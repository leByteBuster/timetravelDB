package dataadapter

import (
	"fmt"
	"log"

	"github.com/LexaTRex/timetravelDB/utils"
)

// NOTE:
// this data loader stores the element id's as neo4j properties nodeid and relationid
// since neo4j ids cannot be set manually

// TODO: nodes should be able to contain a list of labels. so far they only are able to
//			 include one label

type FileLoadResult struct {
	data []map[string]interface{}
	err  error
}

// template string allows the loading of generated files other
// than graph_nodes.json and graph_edges.json in the form of
// graph_nodes<template>.json. This way we can keep multiple
// generated data files. pass empty string to use generated default files
func LoadData(template string) {

	var graphNodesResult FileLoadResult
	var graphEdgesResult FileLoadResult

	nodesChan := make(chan FileLoadResult)
	edgesChan := make(chan FileLoadResult)

	nodePath := "data-generator/generated-data/graph_nodes.json"
	edgePath := "data-generator/generated-data/graph_edges.json"

	if template != "" {

		fmt.Println("Load template: ", template)

		nodePath = "data-generator/generated-data/graph_nodes" + template + ".json"
		edgePath = "data-generator/generated-data/graph_edges" + template + ".json"

	}

	go loadJsonFileAsync(nodePath, nodesChan)
	go loadJsonFileAsync(edgePath, edgesChan)

	graphNodesResult = <-nodesChan
	graphEdgesResult = <-edgesChan

	if graphNodesResult.err != nil {
		log.Printf("Error loading nodes from json: %v", graphNodesResult.err)
	}

	if graphEdgesResult.err != nil {
		log.Printf("Error loading edges from json: %v", graphEdgesResult.err)
		return
	}

	// TODO: maybe split preparing data and loading into 2 functions
	// prepare and format data from templates to load in neo4j, load to neo4j database
	// timeSeriesMap holds all time series data which was extracted in the process
	timeSeriesMap, err := loadDataNeo4j(graphNodesResult.data, graphEdgesResult.data)
	if err != nil {
		log.Printf("Error loading Neo4j data: %v", err)
		return
	}
	// load time series data to TimescaleDB database
	loadDataTimeScaleDB(timeSeriesMap)
}

func loadJsonFileAsync(filePath string, resultChan chan<- FileLoadResult) {
	fileLoadResult := FileLoadResult{}
	data, err := utils.LoadJsonData(filePath)
	fileLoadResult.data = data
	fileLoadResult.err = err
	resultChan <- fileLoadResult
}
