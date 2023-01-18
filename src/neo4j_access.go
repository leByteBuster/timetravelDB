package main

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func exampleAccessNeo4j(ctx context.Context, uri, username, password string, cypherQueryString string) {
	// Connect to the Neo4j database
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	res, err := session.Run(ctx, cypherQueryString, map[string]interface{}{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Result: ", res)
}
