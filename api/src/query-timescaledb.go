package main

import (
	"fmt"
	"log"
	"strings"
)

// get property / single time-series for the period (from,to], apply aggregation on it if not empty
func getPropertyAggr(from, to, aggr, table string) interface{} {

	// build query string
	query := fmt.Sprintf("SELECT %s(value) FROM %s WHERE time >= '%s' AND time < '%s'", aggr, table, from, to)
	fmt.Printf("Query: %v\n", query)

	// res := queryTimeScale(query, username, password, port, dbname)
	val, err := readRowValueTimescale(query, nil, UserTS, PassTS, PortTS, DBnameTS)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}

	return val
}

// get properties / multiple time-series and apply aggrergation on it if not empty
func getProperties(from, to, aggr string, tables []string) (aggrRes interface{}, res []TimeSeriesRow) {

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

	rows, err := readRowsTimescale(query, nil, UserTS, PassTS, PortTS, DBnameTS)

	if err != nil {
		log.Println(err)
	}

	if aggr != "" && len(rows) == 1 {
		return rows[0].value, nil
	} else if aggr != "" && len(rows) > 1 {
		panic("Aggregation returned more than one row")
	}

	return nil, rows
}
