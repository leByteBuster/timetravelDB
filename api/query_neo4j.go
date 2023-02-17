package api

import (
	"context"
	"errors"
	"log"

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

// format the result of a neo4j query to a map of arrays
// every entry in the map represents a column in the result (a variable of the match clause)

// maybe instead of assigning the element (neo4j.node or neo4j.relationship - which is elRec here) to the arrays
// of the map we could consider assigning *neo4j.node or *neo4j.relationship so we can access and change the properties of the node/relationship
// by reference once we retrieved them without having to access them through the map and array again
// try this out
func resultToMap(res neo4j.ResultWithContext) (map[string][]interface{}, error) {
	var formRes map[string][]any = map[string][]any{}
	for res.Next(context.Background()) {

		record := res.Record()

		for _, el := range record.Keys {

			// TODO: test if indexed map is faster (see helpers)
			elRec, ok := record.Get(el)
			if !ok {
				return nil, errors.New("Error getting value for column: " + el)
			}
			if formRes[el] == nil {
				formRes[el] = []interface{}{elRec}
			} else {
				formRes[el] = append(formRes[el], elRec)
			}
		}
	}

	if res.Err() != nil {
		log.Fatalf("\nNext error on query result: %v", res.Err())
	}
	return formRes, nil

}
