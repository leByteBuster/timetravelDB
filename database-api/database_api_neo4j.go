package databaseapi

import (
	"context"
	"log"
	"strings"

	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var ConfigNeo = Neo4jConfig{}

var DriverNeo neo4j.DriverWithContext
var SessionNeo neo4j.SessionWithContext

func ConnectNeo4j() (neo4j.DriverWithContext, error) {

	// this might be needed for benchmarking. Right now i adjust the config file in the docker container by hand
	// config := neo4j.Config{
	// 	MaxTransactionRetryTime:      0, // disable timeout
	// 	MaxConnectionLifetime:        0,
	// 	ConnectionAcquisitionTimeout: 0,
	// 	SocketConnectTimeout:         0,
	// }

	var sb strings.Builder
	sb.WriteString("neo4j://")
	sb.WriteString(ConfigNeo.Host)
	sb.WriteString(":")
	sb.WriteString(ConfigNeo.Port)

	driver, err := neo4j.NewDriverWithContext(sb.String(), neo4j.BasicAuth(ConfigNeo.Username, ConfigNeo.Password, ""))
	if err != nil {
		return nil, err
	}
	return driver, nil

}

func SessionNeo4j(ctx context.Context, driver neo4j.DriverWithContext) neo4j.SessionWithContext {
	sessionNeo := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	return sessionNeo
}

// send any read query and return the results as a key value map
func ReadQueryNeo4j(query string) (neo4j.ResultWithContext, error) {

	res, err := SessionNeo.Run(context.Background(), query, map[string]interface{}{})

	if err != nil {
		log.Printf("%v: error executing neo4j query: %v", err, query)
		return nil, err
	}

	return res, nil
}

// the following functions are used by the data-adapters

func WriteQueryNeo4j(ctx context.Context, query string, params map[string]interface{}) {
	res, err := SessionNeo.Run(ctx, query, params)

	utils.Debugf("Neo4j write query: %v", query)
	utils.Debugf("res: %v", res)

	if err != nil {
		log.Printf("%v: error executing neo4j query: %v", err, query)
		return
	}
}

func WriteQueryMultipleNeo4j(ctx context.Context, queries []string, params map[string]interface{}) {
	for _, query := range queries {
		WriteQueryNeo4j(ctx, query, params)
	}
}

func ClearNeo4j() {
	query := "MATCH (n) DETACH DELETE n"
	WriteQueryNeo4j(context.Background(), query, map[string]interface{}{})
}
