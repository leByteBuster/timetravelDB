package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

// this function takes a map with uiids as keys and lists of maps (with strings as keys and values of arbitrary type) as values
// every of theses lists represent a time-series which can
func loadDataTimeScaleDB(timeSeriesMapNodes map[uuid.UUID][]map[string]interface{}, timeSeriesMapEdges map[uuid.UUID][]map[string]interface{}) {

	// maybe solve "if timestamp" differently. For example per grouping tables or something like this ?
	creatTablePrefix := `CREATE TABLE IF NOT EXISTS `

	for uuid, valueArr := range timeSeriesMapNodes {
		// probably deletable
		uiidUnderscore := strings.ReplaceAll(uuid.String(), "-", "_")
		tablename := "ts_" + uiidUnderscore
		createTableQuery := fmt.Sprintf("%v %v %v", creatTablePrefix, tablename, ` (time TIMESTAMPTZ NOT NULL, timestamps BOOLEAN NOT NULL, `)

		fmt.Printf("\nDatatype of value fields: %v\n", reflect.TypeOf(valueArr[0]["Value"]))
		switch value := valueArr[0]["Value"].(type) {
		// TODO: falls wir komplexe datentypen als time-series values erlauben
		case map[string]interface{}:
			println(value)
			// // datentyp für die map erstellen: TODO if needed
			// typeString := `CREATE TYPE MAP (`
			// for k, v := range value {
			// 	vtype := v.(type)
			// 	typeString += ` ` + key + ` ` + vtype + `, `
			// }
			// // TODO: komma entfernen
			// typeString += `); `
			// // TODO Datentyp in datenbank kreieren
			// // TODO wahrscheinlich braucht der typ "map" einen unique value für jede art von map
			// createTableQuery += "map MAP );"

		// TODO: falls wir listen als datentypen verwenden
		// PS: ggf muss ich dann noch einen switch über den typ der array elemente  machen
		case []interface{}:
			// TODO similar as map[string]interface{}
		case string:
			// put this in function
			createTableQuery += "value TEXT);"
			fmt.Printf("Datatype string")
		case int:
			// put this in function
			createTableQuery += "value INTEGER);"
			fmt.Printf("Datatype int")
		// TODO: maybe change this
		case decimal.Decimal:
			// put this in function
			createTableQuery += "value DECIMAL);"
			fmt.Printf("Datatype decimal")
		case float32:
			// put this in function
			createTableQuery += "value DECIMAL);"
			fmt.Printf("Datatype float32")
		case float64:
			// put this in function
			createTableQuery += "value DECIMAL);"
			fmt.Printf("Datatype float64")
		case bool:
			// put this in function
			createTableQuery += "value BOOLEAN);"
			fmt.Printf("Datatype float64")

		default:
			fmt.Printf("No Datatype fits")
			createTableQuery += "value decimal);"
		}

		// create the table according to  the data type
		conn := connectTimescale("postgres", "password", "5432", "postgres")
		defer conn.Close(context.Background())

		_, err := conn.Exec(context.Background(), createTableQuery)
		if err != nil {
			fmt.Println("\nError creating table:", err)
			fmt.Printf("\n Create Table Query: %v ", createTableQuery)
			panic("lel")
		}

		for _, timeSeriesEntry := range valueArr {
			insertQuery := fmt.Sprintf("INSERT INTO %v (time, timestamps, value) VALUES ($1, false, $2)", tablename)
			fmt.Printf("\n Insert Query String: INSERT INTO %v (time, timestamps, value) VALUES (%v, false, %v)", tablename, timeSeriesEntry["Start"], timeSeriesEntry["Value"])
			ret, err := conn.Exec(context.Background(), insertQuery, timeSeriesEntry["Start"], timeSeriesEntry["Value"])
			fmt.Printf("\nInsert Return: %v", ret)
			if err != nil {
				fmt.Println("\nError inserting data:", err)
				fmt.Printf("\n Insert Query: %v ", insertQuery)
			}
		}
	}

	//   `time TIMESTAMPTZ NOT NULL,
	//   value TEXT NOT NULL,
	//   price DOUBLE PRECISION NULL,
	//   day_volume INT NULL
	// );`

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
