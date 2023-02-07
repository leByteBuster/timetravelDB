package api

import (
	"context"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func queryNeo4j(query string) (any, error) {
	// TODO: maybe introduce some type handling here
	// 			 otherwise i could as well use the lower func queryReadNeo4j
	return queryReadNeo4j(context.Background(), query)
}

func queryNodeNeo4j(id int) (neo4j.Node, error) {
	qval, err := queryReadNeo4j(context.Background(), "MATCH (n) WHERE n.nodeid = "+strconv.Itoa(id)+" RETURN n")
	return qval.(neo4j.Node), err
}

func queryNodesNeo4j(id []int) (neo4j.Node, error) {
	// TODO
	return neo4j.Node{}, nil
}

func queryPropsNodeNeo4j(id int) (map[string]string, error) {
	qval, err := queryReadNeo4j(context.Background(), "MATCH (n) WHERE n.nodeid = "+strconv.Itoa(id)+" RETURN n")
	return convertMapStr(qval.(neo4j.Node).Props), err
}

func queryPropsNodesNeo4j(id []int) (neo4j.Node, error) {
	// TODO
	return neo4j.Node{}, nil
}
