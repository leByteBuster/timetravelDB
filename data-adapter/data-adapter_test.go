package dataadapter

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	databaseapi "github.com/LexaTRex/timetravelDB/database-api"
	"github.com/LexaTRex/timetravelDB/utils"
)

func setWorkingDirectory() {
	err := os.Chdir("/home/lenovob/Development/timetravelDB")
	if err != nil {
		log.Fatalf("Error changing directory: %v", err)
	}
}

func setupDatabases() {

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)

	var NeoErr error
	var TsErr error
	var ConfigErr error

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
		databaseapi.SessionNeo = databaseapi.SessionNeo4j(ctx, databaseapi.DriverNeo)
	}

	// connect to TimescaleDB database
	databaseapi.SessionTS, TsErr = databaseapi.ConnectTimescale(databaseapi.ConfigTS.Username, databaseapi.ConfigTS.Password, databaseapi.ConfigTS.Port, databaseapi.ConfigTS.Database)
	if TsErr != nil {
		log.Fatalf("creating ts connection failed: %v", TsErr)
	}

}

func BenchmarkLoadData(b *testing.B) {

	utils.DEBUG = false

	setWorkingDirectory()
	setupDatabases()

	ctx := context.Background()
	defer databaseapi.DriverNeo.Close(ctx)
	defer databaseapi.SessionNeo.Close(ctx)
	defer databaseapi.SessionTS.Close()

	iterations := 5

	fmt.Println("benchmark iterations: ", iterations)

	databaseapi.ClearData()

	for i := 0; i < iterations; i++ {
		LoadData("")
	}
}
