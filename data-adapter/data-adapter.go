package dataadapter

import (
	"fmt"
	"log"

	"github.com/LexaTRex/timetravelDB/utils"
)

// NOTE: we do not add the nodeid or relationid. neo4j itself sets a these ids
// starting from 0 with the first node added to the datatbase. This should
// match the ids from the generated data.

// NOTE:
// this data loader stores the element id's as neo4j properties nodeid and relationid
// since neo4j ids cannot be set manually

// TODO: nodes should be able to contain a list of labels. so far they only are able to
//			 include one label

// TODO: use go routines to parallize data loading

// template string allows the loading of generated files other
// than graph_nodes.json and graph_edges.json in the form of
// graph_nodes<template>.json. This way we can keep multiple
// generated data files.
// TODO: This is not the cleanest way. change this
func LoadData(template string) {

	var graph_nodes []map[string]interface{}
	var graph_edges []map[string]interface{}
	var err error

	if template == "" {
		graph_nodes, err = utils.LoadJsonData("data-generator/generated-data/graph_nodes.json")
		if err != nil {
			log.Printf("Error loading nodes from json: %v", err)
		}
		graph_edges, err = utils.LoadJsonData("data-generator/generated-data/graph_edges.json")
		if err != nil {
			log.Printf("Error loading edges from json: %v", err)
			return
		}
	} else {
		fmt.Println("Load template: ", template)
		graph_nodes, err = utils.LoadJsonData("data-generator/generated-data/graph_nodes" + template + ".json")
		if err != nil {
			log.Printf("Error loading nodes from json: %v", err)
		}
		graph_edges, err = utils.LoadJsonData("data-generator/generated-data/graph_edges" + template + ".json")
		if err != nil {
			log.Printf("Error loading edges from json: %v", err)
			return
		}
	}

	// TODO: maybe split preparing data and loading into 2 functions
	// prepare and format data from templates to load in neo4j, load to neo4j database
	// timeSeriesMap holds all time series data which was extracted in the process
	timeSeriesMap, err := loadDataNeo4j(graph_nodes, graph_edges)
	if err != nil {
		log.Printf("Error loading Neo4j data: %v", err)
		return
	}
	// load time series data to TimescaleDB database
	loadDataTimeScaleDB(timeSeriesMap)
}
