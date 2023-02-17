package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func Api() {
	// counts only one row because the period (from,to] is exclusive for the second value
	// val := getPropertyAggr("2022-12-22T15:33:13Z", "2022-12-29T20:24:36.311106Z", "COUNT", "ts_05318d0f_6a49_4e67_b9a5_62b46af5c209")

	// counts two rows because from is one milisecond larger and is included in the next to of (to,from]
	// val2 := getPropertyAggr("2022-12-22T15:33:13Z", "2022-12-29T20:24:36.311107Z", "COUNT", "ts_05318d0f_6a49_4e67_b9a5_62b46af5c209")
	// fmt.Printf("Aggr: %v\n", val)
	// fmt.Printf("Aggr: %v", val2)

	ctx := context.Background()

	var err error
	DriverNeo, err = neo4j.NewDriverWithContext(UriNeo, neo4j.BasicAuth(UserNeo, PassNeo, ""))
	if err != nil {
		log.Printf("Creating driver failed: %v", err)
		os.Exit(1)
	}
	defer DriverNeo.Close(ctx)

	SessionNeo = DriverNeo.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer SessionNeo.Close(ctx)

	// TEST QUERIES
	//String ttQuery2 = "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33" + "RETURN n.ping, n ";
	//String ttQuery3 = "FROM 2023-02-03T12:34:39Z TO 2023-02-03T14:34:39.2222Z SHALLOW MATCH (a)-[x]->(b) " + "RETURN a.ping, b "; // should parse
	//String ttQuery4 = "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b "; // should parse
	//String ttQuery = "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a "; // should parse
	//String ttQuery6 = "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a) WHERE a.ping > 22" + " RETURN a "; // should parse
	//String ttQuery5 = "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b "; // should parse
	//ttQuery4 := "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b " // should parse
	//ttQuery5 := "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a "         // should parse
	//ttQuery6 := "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a) WHERE a.ping > 22" + " RETURN a "                  // should parse
	//ttQuery5 := "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b " // should parse

	//ttQuery5 := "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b " // should parse

	//query := "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (a)-[x]->(b) WHERE a.ping IS NOT NULL" + " RETURN a.ping, b "

	// ####
	// correct path of the condition tree:
	// ####
	// query := "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z SHALLOW MATCH (a)-[x]->(b) WHERE a.ping IS NOT NULL" + " RETURN a.ping, b "
	// query := "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z SHALLOW MATCH (a)-[x]->(b) WHERE a.ping IS NOT NULL" + " RETURN  b "
	// query := "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z SHALLOW MATCH (a)-[x]->(b) WHERE a.ping" + " RETURN  b "
	// query := "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z SHALLOW MATCH (a)-[x]->(b) WHERE a.ping" + " RETURN  b.st "
	// query := "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z SHALLOW MATCH (a)-[x]->(b) WHERE b.st IS NOT NULL AND a.ping IS NOT NULL" + " RETURN b "
	// query := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu IS NOT NULL" + " RETURN  b.st "

	// ####
	// ####
	// ####
	// ####
	// PRESENTATION QUERIES 15.02
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu IS NOT NULL RETURN *
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu IS NOT NULL RETURN  a.properties_components_cpu
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.notExistingProperty IS NOT NULL RETURN *
	// query := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu IS NOT NULL RETURN  a.properties_components_cpu"
	// ####
	// ####
	// ####

	// ####
	// #### QUERIES NACH CASES: siehe schaubild bzw den condition tree mit den checkpoints (prints)
	// #### queries case 1:
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu = 'UGWJn' RETURN  * | should return 2 match (a)-[x]->(b) with a.properties_components_cpu = 'UGWJn'
	//																																																																												 because the node with a.properties_components_cpu = 'UGWJn' is occouring in a pattern like this twice
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu = 'not available' RETURN  * | should return nothing
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE "aa" = a.properties_components_cpu RETURN  * | should return nothing
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE 40 > b.properties_Risc RETURN  * | should return one pattern where node 104 is a node 105 is b and all 3 entries for the property risc on 105
	// 																																																																			  | because two of the entries are < 40
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE 24 > b.properties_Risc RETURN  * | should return one pattern where node 104 is a node 105 is b and all 3 entries for the property risc on 105
	// 																																																																			  | because one of the entries is < 24
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE 23 > b.properties_Risc RETURN  * | should return no pattern because no of the entries are < 23
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE 23 >= b.properties_Risc RETURN  * | same as two above and below
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE b.properties_Risc =< 23  RETURN  * | same as above and three above
	//    FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE b.properties_Risc =< 23  RETURN  * | same as above and three above
	//    FROM 2023-01-01T14:33:00Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE b.properties_Risc <= 23  RETURN  * | should return nothing since the time range is not including the value with 23
	//    FROM 2023-01-01T14:33:00Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE b.properties_Risc <= 33  RETURN  * | should return one match entry with all propertyies for b.properties_Risc (like above again)
	//    FROM 2023-01-01T14:33:00Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE b.properties_Risc < 33  RETURN  * | should nothing  nothing since the time range lowest value is 33 but we ask for < 33

	// THE LAST ONE IS THE ONLY ONE NOT WORKING YET

	fmt.Print("Enter Query: \n")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("TTQL> "),
		prompt.OptionTitle("TTQL CLI"),
	)

	p.Run()

}

// executor for the prompt
func executor(in string) {
	switch in {
	case "hello":
		fmt.Println("Hello world!")
	case "quit":
		handleExit() // note: cannot handle this with defer (probably goroutine race condition)
		os.Exit(0)
	case "exit":
		handleExit()
		os.Exit(0)
	case "Quit":
		handleExit()
		os.Exit(0)
	case "Exit":
		handleExit()
		os.Exit(0)
	case "help":
		fmt.Println("To query TimeTravelDB type in a valid TTQL query.\nTo exit the program type 'quit' or 'exit'.\nFor more infos: https://github.com/LexaTRex/timetravelDB/")
	default:
		fmt.Printf("Processing Query: %s\n", in)
		err := ProcessQuery(in)
		if err != nil {
			log.Printf("Failed: %v", err)
		}
	}
}

// auto completion suggestions for the prompt
func completer(in prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{{Text: "FROM", Description: "FROM"}, {Text: "TO", Description: "TO"}, {Text: "SHALLOW", Description: "SHALLOW"}, {Text: "MATCH", Description: "MATCH"}, {Text: "WHERE", Description: "WHERE"}, {Text: "RETURN", Description: "RETURN"}}
}

// workaround to get the terminal back to normal after exiting the program
func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "sane", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}
