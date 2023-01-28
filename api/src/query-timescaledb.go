package main

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

// get properties for one node, if aggr != "" apply aggregate function on the result
func getPropertiesFromNode(from, to, aggr, table string) (pgx.Rows, error) {

	// build query string
	query := fmt.Sprint("SELECT %s(value) FROM %s WHERE time >= %s AND time < %s", table, aggr, from, to)

	//TODO: change this so it isnt hard coded anymore. Should be read from config or so
	username := "postgres"
	password := "password"
	port := "5432"
	dbname := "postgres"
	connStr := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s", username, password, port, dbname)
	return sendQuery(connStr, query)
}

// get properties for multiple nodes, if aggr != "" apply aggregate function on the result
func getPropertiesFromNodes(from, to, aggr string, tables []string) {

	var builder strings.Builder

	// build query string
	builder.WriteString("SELECT ")
	builder.WriteString(aggr)
	builder.WriteString("(value)")
	builder.WriteString(" FROM (")
	for i, table := range tables {
		if i > 0 {
			builder.WriteString(" UNION ALL ")
		}
		builder.WriteString("SELECT ")
		builder.WriteString("value FROM ")
		builder.WriteString(table)
		builder.WriteString(" WHERE time >= ")

		//TODO: CHANGE THIS TO DATETIME
		builder.WriteString(from)
		builder.WriteString(" AND time <= ")

		//TODO: CHANGE THIS TO DATETIME
		builder.WriteString(from)
		builder.WriteString(to)
	}
	fmt.Println(builder.String())

	// Get the query string from the StringBuilder
	query := builder.String()
	fmt.Println(query)
}
