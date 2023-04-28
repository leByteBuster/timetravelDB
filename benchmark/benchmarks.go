package benchmark

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	dataadapter "github.com/LexaTRex/timetravelDB/data-adapter"
	dataadapterneo4j "github.com/LexaTRex/timetravelDB/data-adapter-neo4j-only"
	datagenerator "github.com/LexaTRex/timetravelDB/data-generator"
	databaseapi "github.com/LexaTRex/timetravelDB/database-api"
	"github.com/LexaTRex/timetravelDB/query-processor/parser"
	qpe "github.com/LexaTRex/timetravelDB/query-processor/qpengine"
	"github.com/LexaTRex/timetravelDB/utils"
)

var X = "2022-12-01T00:00:00Z"
var Y = "2022-12-02T00:00:00Z"

var responseTimes [7][100]time.Duration

func RunBenchmark() {

	utils.DEBUG = false

	templates := []string{"graph_template_bm_1.yaml", "graph_template_bm_2.yaml", "graph_template_bm_3.yaml", "graph_template_bm_4.yaml"}

	// generate and test each dataset
	for i, template := range templates {

		// generate next dataset
		datagenerator.GenerateData(template)

		/////////////////////////////
		// cypher query benchmark: //
		/////////////////////////////

		// clear databases
		databaseapi.ClearTTDB()

		// load data to neo4j for cypher queries
		dataadapterneo4j.LoadData()

		RunBenchmarkNeo4j(i)

		/////////////////////////////
		// ttdb query benchmark: ////
		/////////////////////////////

		// clear databases
		databaseapi.ClearTTDB()

		// load data to ttdb for ttql queries
		dataadapter.LoadData()

		// deep ttql queries
		RunBenchmarkTTDB(i, false)

		// shallow ttql queries
		RunBenchmarkTTDB(i, true)
	}
}

// runs the benchmark queries on only Neo4j with the current contained data
func RunBenchmarkNeo4j(dataset int) {

	var results []map[string][]any

	for i, query := range neo4jQueries {

		// warmup run
		for j := 0; j < 10; j++ {
			databaseapi.ReadQueryNeo4j(query)
		}

		// measurement run
		for j := 0; j < 100; j++ {

			start := time.Now()

			res, err := databaseapi.ReadQueryNeo4j(query)

			if err != nil {
				log.Printf("\n error querying neo4j: %v", err)
				return
			}

			result, err := qpe.ResultToMap(res)

			if err != nil {
				log.Printf("\nerror on result: %v", err)
				return
			}

			end := time.Now()

			results = append(results, result)
			utils.UNUSED(results)

			// Calculate the response time
			responseTime := end.Sub(start)
			responseTimes[i][j] = responseTime
		}
	}

	fmt.Println("*********************************************")
	fmt.Println("************* BENCHMARK RESULTS *************")
	fmt.Println("*********************************************")
	fmt.Println()

	writeResultToFile("neo4j-benchmar-result"+strconv.Itoa(dataset), responseTimes)

	for i := range neo4jQueries {
		fmt.Println("Query : 	", i)
		fmt.Print("Response Time : [")
		for j := 0; j < 100; j++ {
			fmt.Print(responseTimes[i][j], ", ")
			// fmt.Println("Result : ", results[i])
		}
		fmt.Print("]")
		fmt.Println()
	}

	fmt.Println("*********************************************")
	fmt.Println("*********** BENCHMARK RESULTS END ***********")
	fmt.Println("*********************************************")
}

// runs the benchmark queries on TTDB (Neo4j & TimescaleDB) with the current contained data
func RunBenchmarkTTDB(dataset int, shallow bool) {

	var results []map[string][]any
	var queries []string

	if shallow {
		queries = ttDBQueriesShallow
	} else {
		queries = ttDBQueries
	}

	for i, query := range queries {

		query = cleanQuery(query)

		// warmup run
		for j := 0; j < 10; j++ {
			// Perform the database query
			queryInfo, _ := parser.ParseQuery(query)
			qpe.ProcessQuery(queryInfo)
		}

		// measurement run
		for j := 0; j < 100; j++ {
			start := time.Now()

			// Perform the database query
			queryInfo, err := parser.ParseQuery(query)
			if err != nil {
				log.Fatalf("\n%v: error parsing query", err)
				return
			}
			res, err := qpe.ProcessQuery(queryInfo)
			if err != nil {
				log.Fatalf("processing query failed: %v", err)
				return
			}

			end := time.Now()

			results = append(results, res)
			utils.UNUSED(results)

			// Calculate the response time
			responseTime := end.Sub(start)
			responseTimes[i][j] = responseTime

		}

	}

	fmt.Println("*********************************************")
	fmt.Println("************* BENCHMARK RESULTS *************")
	fmt.Println("*********************************************")
	fmt.Println()

	if shallow {
		writeResultToFile("ttdb-shallow-benchmark-result"+strconv.Itoa(dataset), responseTimes)
	} else {
		writeResultToFile("ttdb-benchmark-result-dataset"+strconv.Itoa(dataset), responseTimes)
	}

	for i := range ttDBQueries {
		fmt.Println("Query : 	", i)
		fmt.Print("Response Time : [")
		for j := 0; j < 100; j++ {
			fmt.Print(responseTimes[i][j], ", ")
			// fmt.Println("Result : ", results[i])
		}
		fmt.Print("]")
		fmt.Println()
	}

	fmt.Println("*********************************************")
	fmt.Println("*********** BENCHMARK RESULTS END ***********")
	fmt.Println("*********************************************")
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

func cleanQuery(query string) string {
	// Split the query into three parts: MATCH, WHERE, and RETURN
	parts := strings.Split(query, "WHERE")
	if len(parts) != 2 {
		// Query doesn't have a WHERE clause, return the original query
		return query
	}

	matchClause := parts[0]
	whereClause := parts[1]
	returnClause := ""
	returnIndex := strings.Index(whereClause, "RETURN")
	if returnIndex != -1 {
		returnClause = whereClause[returnIndex:]
		whereClause = whereClause[:returnIndex]
	}

	// Replace any() with a.prop
	whereClause = strings.ReplaceAll(whereClause, "any(", "")
	whereClause = strings.ReplaceAll(whereClause, "ANY(", "")
	whereClause = strings.ReplaceAll(whereClause, "Any(", "")
	whereClause = strings.ReplaceAll(whereClause, ")", "")

	// Reassemble the query
	query = matchClause + "WHERE " + whereClause + returnClause
	return query
}

func writeResultToFile(filename string, result [7][100]time.Duration) {
	file, err := os.Create(filename + ".csv")

	// Create a new file for writing
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a new buffered writer
	csvwriter := csv.NewWriter(file)
	defer csvwriter.Flush()

	// Loop through the 2D array and write each element to the file

	// Loop through the 2D array and write each row to the CSV file
	for _, row := range result {
		// Create a new string slice to hold the string values of each element in the row
		var strRow []string
		for _, val := range row {
			// Convert the time.Duration value to a string
			strVal := strconv.Itoa(int(val))
			// Add the string value to the row slice
			strRow = append(strRow, strVal)
		}
		// Write the row of string values to the CSV file
		err := csvwriter.Write(strRow)
		if err != nil {
			panic(err)
		}
	}

}

// queries are  put to the end of the file for better readability

var neo4jQueries []string = []string{
	// todo: check: use n.nodeid or ID(n)

	////////////
	//// #1 ////
	////////////

	// query single node with id 0 if intersection with X to Y
	// this is not the whole query ! we still gotta filter all the properties for the time
	// could not find an cypher equivalent to the ttql query. I would need to know every property name beforehand
	// probably somehow possible with apoc query
	// which adds additional functionality
	// see: ~/Documents/school/master_thesis/research_fetch_single_node_filter_properties_only_neo4j
	// NOTE: time series inside elements not filtered by time interval
	"MATCH (n) WHERE n.nodeid = 0 AND n.end >= datetime('" + X + "') AND n.start <= datetime('" + Y + "') RETURN n",

	////////////
	//// #2 ////
	////////////

	// query single time series property of single node from X to Y
	// with all time series entries that between  X and Y
	// TODO:
	// - [x] ordering should be ordered because the elements are ordered by ts_Risc_1_time, ts_Risc_2_time, ..
	// - [x] addting timestamps to output
	`MATCH (n)
		WHERE n.nodeid = 0 AND n.end >= datetime('` + X + `') AND n.start <= datetime('` + Y + `')` +
		`WITH n, [prop IN keys(n) WHERE prop STARTS WITH 'ts_Risc_'] AS riscProps
  	  WITH n, size(riscProps) / 2 AS ts_Risc_numbers
			WITH n, range(0, ts_Risc_numbers - 1) AS indices
	 RETURN
			REDUCE(acc = [], i IN indices |
				acc + CASE
					WHEN n['ts_Risc_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND n['ts_Risc_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') THEN [n['ts_Risc_' + i + '_time'],n['ts_Risc_' + i + '_value']]
					ELSE null
				END
			) AS props`,

	////////////
	//// #3 ////
	////////////

	// query a time series property of all nodes (that have this property)
	// TODO:
	// - [x] ordering
	// - [x] addting timestamps to output
	`
	MATCH (n)
 	  WHERE n.end >= datetime('` + X + `') AND n.start <= datetime('` + Y + `') 
		WITH n, [prop IN keys(n) WHERE prop STARTS WITH 'ts_Risc_'] AS riscProps
 	   WITH n, size(riscProps) / 2 AS ts_Risc_numbers
		WITH n, range(0, ts_Risc_numbers - 1) AS indices
	RETURN
		REDUCE(acc = [], i IN indices |
			acc + CASE
				WHEN n['ts_Risc_' + i + '_time'] >= datetime('` + X + `') AND n['ts_Risc_' + i + '_time'] < datetime('` + Y + `') THEN [n['ts_Risc_' + i + '_time'],n['ts_Risc_' + i + '_value']]
				ELSE null
			END
		) AS props`,

	////////////
	//// #4 ////
	////////////

	// query multiple time series properties of a single node
	// TODO:
	// - [x] ordering
	// - [] replace time strings with variable as soon as working
	// - [x] add time string to output
	`MATCH (n)
		WHERE n.nodeid = 0 AND n.end >= datetime('2022-12-01T00:00:00Z') AND n.start <= datetime('2022-12-02T00:00:00Z')
		WITH n, 
		  [prop IN keys(n) WHERE prop STARTS WITH 'ts_IP_'] AS ts_IP_props,
		  [prop IN keys(n) WHERE prop STARTS WITH 'ts_Risc_'] AS ts_Risc_props,
		  [prop IN keys(n) WHERE prop STARTS WITH 'ts_Risc1_'] AS ts_Risc1_props,
		  [prop IN keys(n) WHERE prop STARTS WITH 'ts_Risc2_'] AS ts_Risc2_props
		WITH n, 
		  size(ts_IP_props) / 2 AS ts_IP_numbers,
		  size(ts_Risc_props) / 2 AS ts_Risc_numbers,
		  size(ts_Risc1_props) / 2 AS ts_Risc1_numbers,
		  size(ts_Risc2_props) / 2 AS ts_Risc2_numbers,
		  ts_IP_props,
		  ts_Risc_props,
		  ts_Risc1_props,
		  ts_Risc2_props
		WITH n, 
		  range(0, ts_IP_numbers - 1) AS ts_IP_indices,
		  range(0, ts_Risc_numbers - 1) AS ts_Risc_indices,
		  range(0, ts_Risc1_numbers - 1) AS ts_Risc1_indices,
		  range(0, ts_Risc2_numbers - 1) AS ts_Risc2_indices,
		  ts_IP_props,
		  ts_Risc_props,
		  ts_Risc1_props,
		  ts_Risc2_props
		RETURN
		  REDUCE(acc = [], i IN ts_IP_indices |
		    acc + CASE
		      WHEN n['ts_IP_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND n['ts_IP_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') THEN [n['ts_IP_' + i + '_time'], n['ts_IP_' + i + '_value']]
		      ELSE null
		    END
		  ) AS ts_IP_props,
		  REDUCE(acc = [], i IN ts_Risc_indices |
		    acc + CASE
		      WHEN n['ts_Risc_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND n['ts_Risc_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') THEN [n['ts_Risc_' + i + '_time'], n['ts_Risc_' + i + '_value']]
		      ELSE null
		    END
		  ) AS ts_Risc_props,
		  REDUCE(acc = [], i IN ts_Risc1_indices |
		    acc + CASE
		      WHEN n['ts_Risc1_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND n['ts_Risc1_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') THEN [n['ts_Risc1_' + i + '_time'], n['ts_Risc1_' + i + '_value']]
		      ELSE null
		    END
		  ) AS ts_Risc1_props,
		  REDUCE(acc = [], i IN ts_Risc2_indices |
		    acc + CASE
		      WHEN n['ts_Risc2_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND n['ts_Risc2_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') THEN [n['ts_Risc2_' + i + '_time'], n['ts_Risc2_' + i + '_value']]
		      ELSE null
		    END
		  ) AS ts_Risc2_props`,

	////////////
	//// #5 ////
	////////////

	// query multiple time series properties of all nodes n
	// still not ordered !
	// TODO:
	// - [x] ordering
	// - [] replace time strings with variable as soon as working
	// - [x] add time string to output
	`MATCH (n)
		WHERE n.end >= datetime('2022-12-01T00:00:00Z') AND n.start <= datetime('2022-12-02T00:00:00Z')
		WITH n, 
		  [prop IN keys(n) WHERE prop STARTS WITH 'ts_IP_'] AS ts_IP_props,
		  [prop IN keys(n) WHERE prop STARTS WITH 'ts_Risc_'] AS ts_Risc_props,
		  [prop IN keys(n) WHERE prop STARTS WITH 'ts_Risc1_'] AS ts_Risc1_props,
		  [prop IN keys(n) WHERE prop STARTS WITH 'ts_Risc2_'] AS ts_Risc2_props
		WITH n, 
		  size(ts_IP_props) / 2 AS ts_IP_numbers,
		  size(ts_Risc_props) / 2 AS ts_Risc_numbers,
		  size(ts_Risc1_props) / 2 AS ts_Risc1_numbers,
		  size(ts_Risc2_props) / 2 AS ts_Risc2_numbers,
		  ts_IP_props,
		  ts_Risc_props,
		  ts_Risc1_props,
		  ts_Risc2_props
		WITH n, 
		  range(0, ts_IP_numbers - 1) AS ts_IP_indices,
		  range(0, ts_Risc_numbers - 1) AS ts_Risc_indices,
		  range(0, ts_Risc1_numbers - 1) AS ts_Risc1_indices,
		  range(0, ts_Risc2_numbers - 1) AS ts_Risc2_indices,
		  ts_IP_props,
		  ts_Risc_props,
		  ts_Risc1_props,
		  ts_Risc2_props
		RETURN
		  REDUCE(acc = [], i IN ts_IP_indices |
		    acc + CASE
		      WHEN n['ts_IP_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND n['ts_IP_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') THEN [n['ts_IP_' + i + '_time'], n['ts_IP_' + i + '_value']]
		      ELSE null
		    END
		  ) AS ts_IP_props,
		  REDUCE(acc = [], i IN ts_Risc_indices |
		    acc + CASE
		      WHEN n['ts_Risc_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND n['ts_Risc_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') THEN [n['ts_Risc_' + i + '_time'], n['ts_Risc_' + i + '_value']]
		      ELSE null
		    END
		  ) AS ts_Risc_props,
		  REDUCE(acc = [], i IN ts_Risc1_indices |
		    acc + CASE
		      WHEN n['ts_Risc1_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND n['ts_Risc1_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') THEN [n['ts_Risc1_' + i + '_time'], n['ts_Risc1_' + i + '_value']]
		      ELSE null
		    END
		  ) AS ts_Risc1_props,
		  REDUCE(acc = [], i IN ts_Risc2_indices |
		    acc + CASE
		      WHEN n['ts_Risc2_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND n['ts_Risc2_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') THEN [n['ts_Risc2_' + i + '_time'], n['ts_Risc2_' + i + '_value']]
		      ELSE null
		    END
		  ) AS ts_Risc2_props`,

	////////////
	//// #6 ////
	////////////

	//// query a time series property of all nodes if ANY(prop) > 20 (that have this property)
	//// returns all, but still not ordered
	//
	// this one is a list of values  [value,value,...]
	// they are ordered by time already because they came from ts_Risc_1_time, ts_Risc_2_time, ...
	// TODO:
	// - [] time string replacement with variable
	// ` MATCH (n)
	//  		WHERE n.end >= datetime('2022-12-01T00:00:00Z') AND n.start <= datetime('2022-12-02T00:00:00Z')
	//  			WITH n, [prop IN keys(n) WHERE prop STARTS WITH 'ts_Risc_'] AS riscProps
	//  			WITH n, size(riscProps) / 2 AS ts_Risc_numbers
	//  			WITH n, range(0, ts_Risc_numbers - 1) AS indices
	//  			WITH REDUCE(acc = [], i IN indices |
	//  					acc + CASE
	//  							WHEN n['ts_Risc_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND
	//  								   n['ts_Risc_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') THEN
	//  									 n['ts_Risc_' + i + '_value']
	//  							ELSE null
	//  					END
	//  			) AS props
	//  			WITH props, [prop IN props WHERE prop > 20] AS filteredProps
	//  			WHERE size(filteredProps) > 0
	//  	RETURN props`,

	// this one is with timestamps in result list [timestamp,value,timestamp,value,...]
	// they are ordered by time already because the came from ts_Risc_1, ts_Risc_2, ...
	// THIS IS NOT WORKING: somehow condition n['ts_Risc_' + i + '_value'] > 20 is not applied (at least in browser. check it)
	// this is due to that the list contains numeric as well as datetime values
	`MATCH (n)
	WHERE n.end >= datetime('2022-12-01T00:00:00Z') AND n.start <= datetime('2022-12-02T00:00:00Z')
		WITH n, [prop IN keys(n) WHERE prop STARTS WITH 'ts_Risc_'] AS riscProps
		WITH n, size(riscProps) / 2 AS ts_Risc_numbers
		WITH n, range(0, ts_Risc_numbers - 1) AS indices
		WITH REDUCE(acc = [], i IN indices |
						acc + CASE
										WHEN
											n['ts_Risc_' + i + '_time'] >= datetime('2022-12-01T00:00:00Z') AND
											n['ts_Risc_' + i + '_time'] < datetime('2022-12-02T00:00:00Z') 
											THEN [n['ts_Risc_' + i + '_time'], n['ts_Risc_' + i + '_value']]
										ELSE null
						END
		) AS pairs
		WITH pairs, [i IN range(0, size(pairs) - 1, 2) WHERE pairs[i+1] > 20 | pairs[i]] AS filteredProps
		WHERE size(filteredProps) > 0
	RETURN pairs
		`,

	////////////
	//// #7 ////
	////////////

	// query a node b if it occurs in pattern (a)-[r]->(b)
	// missing: temporal checking for all containing properties (note: prob not possible without apoc, see above)
	// 	temporal checking possible on single time series, but on generic i would need complex string operations
	//  using regex to extract all 	interrelated time series values and time values
	// NOTE: time series inside elements are not filetered by the time interval
	`MATCH (a)-[r]->(b) 
	    WHERE a.end >= datetime('` + X + `') AND a.start <= datetime('` + Y + `') AND 
					  r.end >= datetime('` + X + `') AND r.start <= datetime('` + Y + `') AND
					  b.end >= datetime('` + X + `') AND b.start <= datetime('` + Y + `') 
	 RETURN b`,
}

var ttDBQueries []string = []string{
	"FROM " + X + " TO " + Y + " MATCH (n) WHERE n.nodeid = 0 RETURN n",                                            // [x] partly (no specific ts properties)
	"FROM " + X + " TO " + Y + " MATCH (n) WHERE n.nodeid = 0 RETURN n.ts_Risc",                                    // [x]
	"FROM " + X + " TO " + Y + " MATCH (n) RETURN n.ts_Risc",                                                       // [x]
	"FROM " + X + " TO " + Y + " MATCH (n) WHERE  n.nodeid = 0  RETURN n.ts_IP, n.ts_Risc, n.ts_Risc1, n.ts_Risc2", // []
	"FROM " + X + " TO " + Y + " MATCH (n) RETURN n.ts_IP, n.ts_Risc, n.ts_Risc1, n.ts_Risc2",                      // []
	"FROM " + X + " TO " + Y + " MATCH (n) WHERE any(n.ts_Risc > 20) RETURN n.ts_Risc",                             // [x]
	"FROM " + X + " TO " + Y + " MATCH (a)-[r]->(b) RETURN b",                                                      // [x] partly (no specific ts properties)

	// this wont be possible because i need to check the temporality of all  elements
	// so every element must be falling under a variable
	// "FROM X TO Y MATCH (a)-[*2]->(b) WHERE ANY(b.ts_Risc > 20) RETURN b",
}

var ttDBQueriesShallow []string = []string{
	"FROM " + X + " TO " + Y + " SHALLOW MATCH (n) WHERE n.nodeid = 0 RETURN n",
	"FROM " + X + " TO " + Y + " SHALLOW MATCH (n) WHERE n.nodeid = 0 RETURN n.ts_Risc",
	"FROM " + X + " TO " + Y + " SHALLOW MATCH (n) RETURN n.ts_Risc",
	"FROM " + X + " TO " + Y + " SHALLOW MATCH (n) WHERE n.nodeid = 0 RETURN n.ts_IP, n.ts_Risc, n.ts_Risc1, n.ts_Risc2",
	"FROM " + X + " TO " + Y + " SHALLOW MATCH (n) RETURN n.ts_IP, n.ts_Risc, n.ts_Risc1, n.ts_Risc2",
	"FROM " + X + " TO " + Y + " SHALLOW MATCH (n) WHERE any(n.ts_Risc) > 20 RETURN n.ts_Risc",
	"FROM " + X + " TO " + Y + " SHALLOW MATCH (a)-[r]->(b) RETURN b",

	// this wont be possible because i need to check the temporality of all  elements
	// so every element must be falling under a variable
	// "FROM X TO Y SHALLOW MATCH (a)-[*2]->(b) WHERE ANY(b.ts_Risc > 20) RETURN b",
}
