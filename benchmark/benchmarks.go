package benchmark

import (
	"fmt"
	"log"
	"time"

	"github.com/LexaTRex/timetravelDB/query-processor/parser"
	qpe "github.com/LexaTRex/timetravelDB/query-processor/qpengine"
	"github.com/LexaTRex/timetravelDB/utils"
)

var X = "2022-12-01T00:00:00Z"
var Y = "2022-12-02T00:00:00Z"

var neo4jQueries []string = []string{
	// todo: check: use n.nodeid or ID(n)

	// query single node with id 0 if intersection with X to Y
	// this is not the whole query ! we still gotta filter all the properties for the time
	"MATCH (n) WHERE n.nodeid = 0 AND datetime(n.end) >= datetime(" + X + ") AND datetime(n.start) <= datetime(" + Y + ") RETURN n",

	// query single property of single node from X to Y
	`MATCH (n)
	WHERE n.nodeid = 0
	WITH n, range(1, 99) AS ts_Risc_numbers
	RETURN
		REDUCE(acc = [], i IN ts_Risc_numbers |
			acc + CASE
				WHEN datetime(n['ts_Risc_' + i + '_time']) >= datetime('2022-12-01T00:00:00Z') AND datetime(n['ts_Risc_' + i + '_time']) < datetime('2022-12-02T00:00:00Z') THEN n['ts_Risc_' + i + '_value']
				ELSE null
			END
		) AS props`,

	// query a time series property of all nodes (that have this property)
	// TODO: n.ts_Risc_0, n.ts_Risc_1, ...
	"MATCH (n) WHERE n.end >= X AND n.start <= Y RETURN n.ts_Risc*",

	// query all time series properties of a single node
	// TODO: maybe impossible
	// "MATCH (n) WHERE ID(n) = 10 AND n.end >= X AND n.start <= Y RETURN n.ts_...",

	// query all ts properties of all nodes
	"MATCH (n) WHERE n.end >= X AND n.start <= Y RETURN properties(n) whihc start with ts ?? also check temporal condition ?h",

	// query a time series property of all nodes if ANY(prop) > 20 (that have this property)
	"MATCH (n) WHERE n.end >= X AND n.start <= Y AND (n.ts_property... > 20 OR n.ts_property.. > 20 OR ... ) RETURN n.ts_property*",

	// query a node b if it occurs in pattern (a)-[r]->(b)
	`MATCH (a)-[r]->(b) 
	    WHERE a.end >= X AND a.start <= Y AND 
					  r.end >= X AND r.start <= Y AND
					  b.end >= X AND b.start <= Y AND
	 RETURN b`,
}

var ttDBQueries []string = []string{
	"FROM X TO Y MATCH (n) WHERE ID(n) = 10 RETURN n",
	"FROM X TO Y MATCH (n) WHERE ID(n) = 10 RETURN n.ts_Risc",
	"FROM X TO Y MATCH (n) RETURN n.ts_Risc",
	"FROM X TO Y MATCH (n) WHERE ID(n) = 10 RETURN n.ts_IP, n.ts_Risc, ts_Risc1, ts_Risc2",
	"FROM X TO Y MATCH (n) RETURN n.ts_IP, n.ts_Risc, ts_Risc1, ts_Risc2",
	"FROM X TO Y MATCH (n) AND ANY(n.ts_Risc > 20) RETURN n.ts_Risc",
	"FROM X TO Y MATCH (a)-[r]->(b) RETURN b",
	"FROM X TO Y MATCH (a)-[*2]->(b) WHERE ANY(b.ts_Risc > 20) RETURN b",
}

var ttDBQueriesShallow []string = []string{
	"FROM X TO Y SHALLOW MATCH (n) WHERE ID(n) = 10 RETURN n",
	"FROM X TO Y SHALLOW MATCH (n) WHERE ID(n) = 10 RETURN n.ts_Risc",
	"FROM X TO Y SHALLOW MATCH (n) RETURN n.ts_Risc",
	"FROM X TO Y SHALLOW MATCH (n) WHERE ID(n) = 10 RETURN n.ts_IP, n.ts_Risc, ts_Risc1, ts_Risc2",
	"FROM X TO Y SHALLOW MATCH (n) RETURN n.ts_IP, n.ts_Risc, ts_Risc1, ts_Risc2",
	"FROM X TO Y SHALLOW MATCH (n) AND ANY(n.ts_Risc > 20) RETURN n.ts_Risc",
	"FROM X TO Y SHALLOW MATCH (a)-[r]->(b) RETURN b",
	"FROM X TO Y SHALLOW MATCH (a)-[*2]->(b) WHERE ANY(b.ts_Risc > 20) RETURN b",
}

var responseTimes []time.Duration

// runs the benchmark queries on only Neo4j with the current contained data
func RunBenchmarkNeo4j() {
	for _, query := range neo4jQueries {
		start := time.Now()

		// Perform the database query
		queryInfo, err := parser.ParseQuery(query)
		if err != nil {
			log.Printf("\n%v: error parsing query", err)
			return
		}
		res, err := qpe.ProcessQuery(queryInfo)
		if err != nil {
			log.Fatalf("processing query failed: %v", err)
			return
		}

		end := time.Now()

		utils.UNUSED(res)

		// Calculate the response time
		responseTime := end.Sub(start)
		responseTimes = append(responseTimes, responseTime)

	}

	fmt.Println("*********************************************")
	fmt.Println("************* BENCHMARK RESULTS *************")
	fmt.Println("*********************************************")
	fmt.Println()

	for i, query := range neo4jQueries {
		fmt.Println("Query Benchmark for query:\n 	", query)
		fmt.Println("Response Time: ", responseTimes[i])
		fmt.Println()
	}

	fmt.Println("*********************************************")
	fmt.Println("*********** BENCHMARK RESULTS END ***********")
	fmt.Println("*********************************************")
}

// runs the benchmark queries on TTDB (Neo4j & TimescaleDB) with the current contained data
func RunBenchmarkTTDB() {
	utils.UNUSED(ttDBQueries, ttDBQueriesShallow)
}

//
//
//
//
// Query a single node
//
// `MATCH (n) WHERE ID(n) = 10 AND n.start = X AND n.start = Y RETURN n`
//
// TTDB
//   Query a single node shallow
//   "FROM X TO Y SHALLOW MATCH (n) WHERE ID(n) = 10 RETURN n"
//   Query a single node
//   "FROM X TO Y MATCH (n) WHERE ID(n) = 10 RETURN n"
//
// ON:
//   Query a time series property of a single node (by ID)
//     MATCH (n) WHERE ID(n) = 10 n.start = X AND n.start = Y RETURN n.ts_property*"
// TTDB
//   Query a time series property of a single node (by ID)
//     FROM X TO Y MATCH (n) WHERE ID(n) = 10 RETURN n.ts_property"
//
// ON:
//   Query a time series property of all nodes (that have this property)
//     MATCH (n) WHERE n.start = X AND n.start = Y RETURN n.ts_property*"
// TTDB
//   Query a time series property of all nodes (that have this property)
//     FROM X TO Y MATCH (n) RETURN n.ts_property"
//
// ON:
//   Query all time series properties of a single node
//   "MATCH (n) WHERE ID(n) = 10 AND n.start = X AND n.start = Y RETURN properties(n)"
// TTDB
//   Query all time series properties of a single node
//   "FROM X TO Y MATCH (n) WHERE ID(n) = 10 RETURN n.prop1, n.prop2, ..."
//
//
// ON:
//   Query all time series properties of all nodes
//   "MATCH (n) WHERE n.start = X AND n.start = Y RETURN properties(n)"
// TTDB
//   Query all time series properties of all nodes
//   "FROM X TO Y MATCH (n) RETURN n.prop1, n.prop2, ..."
//
//
//
// // Querying Time Series Data: the ANY OPERATOR
//
// ON:
//   Query a time series property of all nodes (that have this property)
//     MATCH (n) WHERE n.start = X AND n.start = Y AND (n.ts_property... > 20 OR n.ts_property.. > 20 OR ... ) RETURN n.ts_property*"
// TTDB
//   Query a time series property of all nodes (that have this property)
//     FROM X TO Y MATCH (n) AND ANY(n.ts_property > 20) RETURN n.ts_property"
//
//
