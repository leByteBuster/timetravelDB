package api

import (
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func Api() {

	// conn := connectTimescale("postgres", "password", "5432", "postgres")
	// defer conn.Close(context.Background())

	//query := `SELECT time, timestamps, value FROM ts_05318d0f_6a49_4e67_b9a5_62b46af5c209 WHERE value='zFCvbu';`
	//query := `SELECT time, timestamps, value FROM ts_05318d0f_6a49_4e67_b9a5_62b46af5c209;`

	// rows, err := readRowsTimescale(query, nil, "postgres", "password", "5432", "postgres")
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }

	// counts only one row because the period (from,to] is exclusive for the second value
	// val := getPropertyAggr("2022-12-22T15:33:13Z", "2022-12-29T20:24:36.311106Z", "COUNT", "ts_05318d0f_6a49_4e67_b9a5_62b46af5c209")

	// counts two rows because from is one milisecond larger and is included in the next to of (to,from]
	// val2 := getPropertyAggr("2022-12-22T15:33:13Z", "2022-12-29T20:24:36.311107Z", "COUNT", "ts_05318d0f_6a49_4e67_b9a5_62b46af5c209")
	// fmt.Printf("Aggr: %v\n", val)
	// fmt.Printf("Aggr: %v", val2)

	var err error
	driverNeo, err = neo4j.NewDriverWithContext(UriNeo, neo4j.BasicAuth(UserNeo, PassNeo, ""))
	if err != nil {
		log.Printf("Creating driver failed: %v", err)
		os.Exit(1)
	}

	// #### TEST QUERY SINGLE NODE NEO4J ####
	// node, err := queryNodeNeo4j(2)

	// if err != nil {
	// 	log.Printf("Querying node failed: %v", err)
	// 	os.Exit(1)
	// }

	// fmt.Printf("Result: %v\n", node)
	// fmt.Println(reflect.TypeOf(node))

	// #### TEST PARSING QUERY STRING ####

}
