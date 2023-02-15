package api

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func queryNeo4j(query string) (neo4j.ResultWithContext, error) {
	// TODO: maybe introduce some type handling here
	// 			 otherwise i could as well use the lower func queryReadNeo4j
	return queryReadNeo4j(context.Background(), query)
}

// Dont know if I even need this
func queryNodeNeo4j(id int) (neo4j.Node, error) {
	// res, err := queryReadNeo4j(context.Background(), "MATCH (n) WHERE n.nodeid = "+strconv.Itoa(id)+" RETURN n")
	//TODO: check if res contains a single node. if yes return
	return neo4j.Node{}, nil
}

func queryNodesNeo4j(id []int) (neo4j.Node, error) {
	// TODO
	return neo4j.Node{}, nil
}

// func queryPropsNodeNeo4j(id int) (map[string]string, error) {
// 	qval, err := queryReadNeo4j(context.Background(), "MATCH (n) WHERE n.nodeid = "+strconv.Itoa(id)+" RETURN n")
// 	return convertMapStr(qval.(neo4j.Node).Props), err
// }

func queryPropsNodesNeo4j(id []int) (neo4j.Node, error) {
	// TODO
	return neo4j.Node{}, nil
}
