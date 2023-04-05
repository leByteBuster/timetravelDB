package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	dataadapter "github.com/LexaTRex/timetravelDB/data-adapter"
	dataadapterneo4j "github.com/LexaTRex/timetravelDB/data-adapter-neo4j-only"
	datagenerator "github.com/LexaTRex/timetravelDB/data-generator"
	databaseapi "github.com/LexaTRex/timetravelDB/database-api"
	"github.com/LexaTRex/timetravelDB/parser"
	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/c-bata/go-prompt"
)

var NeoErr error
var TsErr error
var ConfigErr error

var HelpString = `
		Hello, welcome to 

			  _____     _____     ____       ____   
			 |_* *_|  U|_^ ^_|   |  _"\   U | __") u     ______     __         __     
			   | |     \ | |    /| | | |   \|  _ \/     /\  ___\   /\ \       /\ \   
			  /| |\     -| |-u  U| |_| |\   | |_) |     \ \ \____  \ \ \____  \ \ \  
			 u |_|U      |_|     |____/ u   |____/       \ \_____\  \ \_____\  \ \_\ 
			 _// \\_   _// \\_    |||_     _|| \\_        \/_____/   \/_____/   \/_/ 
			(__) (__) (__) (__)  (__)_)   (__) (__)                                


	  To greet the world type 'hello' and return.  	
		To exit the program type 'quit' 'q' or 'exit' and return.
		To query TimeTravelDB type in a valid TTQL query.
		To print helpful information type 'help' or '-h' and return.
		To generate test data type 'Generate Data' and return. 
		To define the structure of the generated data alter data-generator/graph_template.yaml. 
		To load generated data into TTDB type 'Load Data' or 'LD' and return.
		To load generated data only Neo4j type 'Load Data Neo4j' or 'LDN' and return. Note that this data cannot be queried via TTDB CLI 
		   but must be queried via the common Neo4j methods. Refer to the README.md for further information. 
		To clear TTDB from all data type 'Clear Data', 'Clear TTDB' or 'CD' and return.
		To clear only Neo4j from all data type 'Clear Neo4j' and return.
		To enable debug mode type 'Debug=1' or '--debug=true' and return. 
		To disable debug mode type 'Debug=0' or '--debug=false' and return. 
		To run benchmarks type 'Benchmark' or 'benchmark' and return. 

		For more infos: https://github.com/LexaTRex/timetravelDB/
		Note: In case the CLI crashes it can occour that the terminal settings are not reset and becomes unusable because of a bug in 
		      the utilized prompt library. In this case close the terminal or type 'reset' and return. Typing will be invisible - just do it anyways..`

func Api() {

	// counts only one row because the period (from,to] is exclusive for the second value
	// val := getPropertyAggr("2022-12-22T15:33:13Z", "2022-12-29T20:24:36.311106Z", "COUNT", "ts_05318d0f_6a49_4e67_b9a5_62b46af5c209")

	// counts two rows because from is one milisecond larger and is included in the next to of (to,from]
	// val2 := getPropertyAggr("2022-12-22T15:33:13Z", "2022-12-29T20:24:36.311107Z", "COUNT", "ts_05318d0f_6a49_4e67_b9a5_62b46af5c209")
	// fmt.Printf("Aggr: %v\n", val)
	// fmt.Printf("Aggr: %v", val2)

	databaseapi.ConfigNeo, databaseapi.ConfigTS, ConfigErr = databaseapi.LoadConfig()
	if ConfigErr != nil {
		log.Fatalf("Error loading config: %v", ConfigErr)
	}

	ctx := context.Background()

	// connect to Neo4j database

	databaseapi.DriverNeo, NeoErr = databaseapi.ConnectNeo4j()
	if NeoErr != nil {
		log.Fatalf("creating neo4j connection failed: %v", NeoErr)
	} else {
		defer databaseapi.DriverNeo.Close(ctx)
		databaseapi.SessionNeo = databaseapi.SessionNeo4j(ctx, databaseapi.DriverNeo)
		defer databaseapi.SessionNeo.Close(ctx)
	}

	// connect to TimescaleDB database
	databaseapi.SessionTS, TsErr = databaseapi.ConnectTimescale(databaseapi.ConfigTS.Username, databaseapi.ConfigTS.Password, databaseapi.ConfigTS.Port, databaseapi.ConfigTS.Database)
	if TsErr != nil {
		log.Fatalf("creating ts connection failed: %v", TsErr)
	}

	defer databaseapi.SessionTS.Close()

	fmt.Println(HelpString)
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("TTQL> "),
		prompt.OptionTitle("TTQL CLI"),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionPreviewSuggestionTextColor(prompt.Black),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	)
	p.Run()

}

// executor for the prompt
func executor(in string) {
	switch in {
	case "hello":
		fmt.Println("Hello, welcome to TTDB CLI !")
	case "quit", "q", "exit", "Exit", "Quit", "EXIT":
		handleExit() // note: cannot handle this with defer (probably goroutine race condition)
		os.Exit(0)
	case "help", "h", "-h", "--help":
		fmt.Println(HelpString)
	case "Generate Data", "GD", "gd":
		datagenerator.GenerateData()
	case "Load Data", "LD", "ld":
		dataadapter.LoadData()
	case "Load Data Neo4j", "LDN", "ldn":
		dataadapterneo4j.LoadData()
	case "Clear Data", "Clear DB", "CD", "cd", "Clear TTDB":
		databaseapi.ClearNeo4j()
		databaseapi.ClearTimescale()
	case "Clear Neo4j":
		databaseapi.ClearNeo4j()
	case "Debug=1", "--debug=1", "-debug=1", "--debug=true", "-debug=true":
		utils.DEBUG = true
	case "Debug=0", "--debug=0", "-debug=0", "--debug=false", "-debug=false":
		utils.DEBUG = false
	case "Benchmark":
		runBenchmark()
	default:
		runTTQLQuery(in)
	}
}

func runBenchmark() {
	// TODO: implement benchmark
	// - how to differe if only neo4j or TTDB
	// -
	panic("unimplemented")
}

func runTTQLQuery(in string) {
	if ConfigErr != nil {
		log.Printf("\n%v: There occured an error paring db configs. Please provide a valid config.yaml and restart the CLI.", ConfigErr)
		return
	}
	if TsErr != nil {
		log.Printf("\n%v: There occured an error connecting to the timescale database. Please provide a running database with the correct credentials (config.yaml) and restart the CLI", TsErr)
		return
	}
	if NeoErr != nil {
		log.Printf("\n%v: There occured an error connecting to the neo4j database. Please provide a running database with the correct credentials (config.yaml) and restart the CLI.", NeoErr)
		return
	}
	utils.Debugf("\nProcessing Query: %s\n", in)
	in = cleanQuery(in)
	queryInfo, err := parser.ParseQuery(in)
	if err != nil {
		log.Printf("\n%v: error parsing query", err)
		return
	}
	queryRes, err := ProcessQuery(queryInfo)
	if err != nil {
		log.Fatalf("processing query failed: %v", err)
		return
	}
	printResult(queryRes, queryInfo)
}

// auto completion suggestions for the prompt
func completer(in prompt.Document) []prompt.Suggest {
	p := []prompt.Suggest{{Text: "FROM"}, {Text: "TO"}, {Text: "SHALLOW"}, {Text: "MATCH"}, {Text: "WHERE"}, {Text: "RETURN"}}
	return prompt.FilterHasPrefix(p, in.GetWordBeforeCursor(), true)
}

// workaround to get the terminal back to normal after exiting the program
func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "sane", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

// this is a dirty hack which just gives the illusion that we can use any() in the WHERE clause
// like any(a.prop) > 23. The functionality is working but parser wise it basically is the same
// as a.prop > 23. This is just for convinience and better usability. Change this in the future.
// Disallow a.prop > 23 and only allow any(a.prop) > 23 in changing the TTQL grammar accordingly.
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
	whereClause = strings.ReplaceAll(whereClause, ")", "")

	// Reassemble the query
	query = matchClause + "WHERE " + whereClause + returnClause
	return query
}

func printResult(queryRes map[string][]any, queryInfo parser.ParseResult) {
	utils.Debugf("\n\n\n                 		 QUERY RESULT\n						%+v\n\n\n", queryRes)
	if len(queryInfo.ReturnProjections) > 0 {
		utils.Debug("\n\n\n                      Printed ordered\n\n\n\n")
		fmt.Printf("%+v\n", utils.JsonStringFromMapOrdered(queryRes, queryInfo.ReturnProjections))
	} else {
		utils.Debug("\n\n\n                      Printed unordered\n\n\n\n")
		fmt.Printf("%+v\n", utils.JsonStringFromMap(queryRes))
	}

}
