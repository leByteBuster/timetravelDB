package databaseapi

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var ConfigTS = TimescaleConfig{}
var ConfigNeo = Neo4jConfig{}
var SessionTS *pgx.Conn

var DriverNeo neo4j.DriverWithContext
var SessionNeo neo4j.SessionWithContext

type TimeSeriesRow struct {
	Timestamp   time.Time
	IsTimestamp bool
	Value       interface{}
}

// send any read query and return the results as a key value map
func QueryReadNeo4j(cypherQueryString string) (neo4j.ResultWithContext, error) {

	res, errReq := SessionNeo.Run(context.Background(), cypherQueryString, map[string]interface{}{})

	if errReq != nil {
		log.Printf("Query failed")
		return nil, errReq
	}

	return res, nil
}

//lint:ignore U1000 Ignore unused function temporarily for debugging

func QueryWriteNeo4j(ctx context.Context, uri, username, password, cypherQueryString string) {

	session := DriverNeo.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err := session.Run(ctx, cypherQueryString, map[string]interface{}{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

//lint:ignore U1000 Ignore unused function temporarily for debugging

func QueryWriteMultipleNeo4j(ctx context.Context, uri, username, password string, cypherQueryStrings []string) {
	// Connect to the Neo4j database

	session := DriverNeo.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	for _, cypherQueryString := range cypherQueryStrings {
		_, err := session.Run(ctx, cypherQueryString, map[string]interface{}{})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func ConnectTimescale(username, password, port, dbname string) (*pgx.Conn, error) {

	var sb strings.Builder
	sb.WriteString("postgresql://")
	sb.WriteString(username)
	sb.WriteString(":")
	sb.WriteString(password)
	sb.WriteString("@localhost:")
	sb.WriteString(port)
	sb.WriteString("/")
	sb.WriteString(dbname)

	// conn, err := pgxpool.Connect(context.Background(), connStr) // use pgxpool for managing multiple connections
	conn, err := pgx.Connect(context.Background(), sb.String())
	return conn, err
}

// send any query string to the database
// this is just for executing, CommandTag only contains status
func ExecTimescale(query string) pgconn.CommandTag {
	res, err := SessionTS.Exec(context.Background(), query)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

// send a list of query strings to the database. not sure if CommandTag contains a result or just status though
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func QueryMultipleTimecale(queries []string, parameters [][]interface{}) []pgconn.CommandTag {
	// create the table according to  the data type
	results := make([]pgconn.CommandTag, 0)
	for i, query := range queries {
		_, err := SessionTS.Exec(context.Background(), query, parameters[i]...)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, pgconn.CommandTag{})
	}
	return results
}

// multiple row READ queries
func ReadRowsTimescale(query string, parameters [][]interface{}) ([]TimeSeriesRow, error) {
	// create the table according to  the data type

	rows, err := SessionTS.Query(context.Background(), query)
	if err != nil {
		log.Println("error querying rows from a table in timescaledb:", err)
		return nil, err
	}
	defer rows.Close()

	res := make([]TimeSeriesRow, 0)

	for rows.Next() {

		var timestamp time.Time
		var isTimestamp bool
		var value interface{}
		if err := rows.Scan(&timestamp, &isTimestamp, &value); err != nil {
			log.Println("Scan failed:", err)
			return nil, err
		}
		res = append(res, TimeSeriesRow{Timestamp: timestamp, IsTimestamp: isTimestamp, Value: value})
	}

	if err := rows.Err(); err != nil {
		log.Println("Error during iteration:", err)
		return nil, err
	}

	return res, nil
}

// single row READ queries
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func ReadRowTimescale(query string, parameters []interface{}) (TimeSeriesRow, error) {

	var timestamp time.Time
	var isTimestamp bool
	var value interface{}

	err := SessionTS.QueryRow(context.Background(), query, parameters...).Scan(&timestamp, &isTimestamp, &value)
	if err != nil {
		return TimeSeriesRow{}, fmt.Errorf("error querying a row from a table in timescaledb: %w", err)
	}

	return TimeSeriesRow{Timestamp: timestamp, IsTimestamp: isTimestamp, Value: value}, nil
}

func ReadRowExistsTimescale(query string) (bool, error) {
	// create the table according to  the data type

	var exists bool
	err := SessionTS.QueryRow(context.Background(), query).Scan(&exists)
	if err != nil {
		return false, nil
		// TODO: reintroduce this as soon as we can be sure that for every UUID in a property in neo4j a
		// table exists in timescaledb
		// return false, fmt.Errorf("error executing an existence check in timescaledb: %w", err)
	}

	return exists, nil
}

// single row only return time READ queries
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func ReadRowTimestampTimescale(query string, parameters []interface{}) (interface{}, error) {
	// create the table according to  the data type

	var timestamp time.Time

	err := SessionTS.QueryRow(context.Background(), query, parameters...).Scan(&timestamp)
	if err != nil {
		return nil, err
	}

	return timestamp, nil
}

// single row only return time READ queries
//lint:ignore U1000 Ignore unused function temporarily for debugging

func ReadRowValueTimescale(query string, parameters []interface{}) (interface{}, error) {
	// create the table according to  the data type

	var value interface{}

	err := SessionTS.QueryRow(context.Background(), query, parameters...).Scan(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}
