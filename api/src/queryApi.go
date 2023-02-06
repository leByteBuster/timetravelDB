package main

// first write a query for returning all values from one or multiples nodes
// - get all nodes from neo4j ()
// - create a struct from every node (or map)
// - for every node match the references to the property in this struct/map
// - query the data of the references return the references
func queryNodes() {}

// query all listed properties from a node / from multiple nodes
func queryPropertyFromNode(properties ...string) {}

// query  one or multiple nodes with an id from neo4j. For all properties of the node get the value
// which is an uuid. Use this uuid to query the related data in timescaleDB where
// the uuid is a reference to a table
//
// at this point the query is supposed to be checked already.
// only queries are supposed to arrive which are parsed,
// examined and devided already
func queryNode(id []string) {

}

// maybe include the following methods into queryNode with using parameters ans zero values
func queryNodeCond(id, conditions []string) {

}

func queryNodeAggr(id, aggregations []string) {

}

func queryNodeCondAggr(id, conditions, aggregations []string) {

}

func queryNodeShallow() {}

// at this point the query is supposed to be checked already.
// only queries are supposed to arrive which are parsed,
// examined and devided already
func queryProperty(query string) {

}
