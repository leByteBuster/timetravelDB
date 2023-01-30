package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func queryNeo4j(ctx context.Context, uri, username, password, cypherQueryString string) {
	// Connect to the Neo4j database
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, errReq := session.Run(ctx, cypherQueryString, map[string]interface{}{})
	if errReq != nil {
		fmt.Println(err)
		return
	}
}

func queryMultipleNeo4j(ctx context.Context, uri, username, password string, cypherQueryStrings []string) {
	// Connect to the Neo4j database
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	for _, cypherQueryString := range cypherQueryStrings {
		_, err := session.Run(ctx, cypherQueryString, map[string]interface{}{})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func queryTimeScale(query, username, password, port, dbname string) {
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Println(err)
	}
}

func queryMultipleTimeScale(queries []string, parameters [][]interface{}, username, password, port, dbname string) {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())
	for i, query := range queries {
		_, err := conn.Exec(context.Background(), query, parameters[i]...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func queryMultipleTimeScaleSameQuery(query string, parameters [][]interface{}, username, password, port, dbname string) {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())
	for i := range parameters {
		_, err := conn.Exec(context.Background(), query, parameters[i]...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func connectTimescale(username, password, port, dbname string) *pgx.Conn {
	connStr := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s", username, password, port, dbname)
	// conn, err := pgxpool.Connect(context.Background(), connStr)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Println("\nUnable to establish connection:", err)
	}
	return conn
}
