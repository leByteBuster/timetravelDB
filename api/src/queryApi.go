package main

// first write a query for returning all values from one or multiples nodes
// - get all nodes from neo4j ()
// - create a struct from every node (or map)
// - for every node match the references to the property in this struct/map
// - query the data of the references return the references
func queryNodes() {}

// query all listed properties from a node / from multiple nodes
func queryPropertyFromNode(properties ...string) {}

// query  one or multiple nodes with an i from neo4j. For all properties of the node get the value
// which is an uuid. Use this uuid to query the related data in timescaleDB where
// the uuid is a reference to a table
func queryNode(id []string) {

}

func queryNodeShallow() {}

func queryProperty() {}
