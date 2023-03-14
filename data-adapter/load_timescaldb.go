package dataadapter

import (
	"fmt"
	"log"
	"os"
	"strings"

	databaseapi "github.com/LexaTRex/timetravelDB/database-api"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// this function takes a map with uiids as keys and lists of maps (with strings as keys and values of arbitrary type) as values
// every of theses lists represent a time-series which can
func loadDataTimeScaleDB(timeseries map[uuid.UUID][]map[string]interface{}) {

	// maybe solve "if timestamp" differently. For example per grouping tables or something like this ?
	creatTablePrefix := `CREATE TABLE IF NOT EXISTS `

	for uuid, values := range timeseries {

		// probably deletable
		uiidUnderscore := strings.ReplaceAll(uuid.String(), "-", "_")
		tablename := "ts_" + uiidUnderscore
		createTableQuery := fmt.Sprintf("%v %v %v", creatTablePrefix, tablename, ` (time TIMESTAMPTZ NOT NULL, timestamps BOOLEAN NOT NULL, `)

		//fmt.Printf("\nDatatype of value fields: %v\n", reflect.TypeOf(values[0]["Value"]))
		switch values[0]["Value"].(type) {
		// TODO: falls wir komplexe datentypen als time-series values erlauben
		case map[string]interface{}:
			//log.Printf("\rInterface: %v", value)
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
			//log.Printf("Datatype []interface{}\r")
		case string:
			createTableQuery += "value TEXT);"
			//log.Printf("Datatype string\r")
		case int:
			createTableQuery += "value INTEGER);"
			//log.Printf("Datatype int\r")
		// TODO: maybe change this
		case decimal.Decimal:
			createTableQuery += "value DECIMAL);"
			//log.Printf("Datatype decimal\r")
		case float32:
			createTableQuery += "value DECIMAL);"
			//log.Printf("Datatype float32\r")
		case float64:
			createTableQuery += "value DECIMAL);"
			//log.Printf("Datatype float64\r")
		case bool:
			createTableQuery += "value BOOLEAN);"
			//log.Printf("Datatype bool\r")
		default:
			//log.Printf("No Datatype fits\r")
			createTableQuery += "value decimal);"
		}

		// create the table according to  the data type
		databaseapi.WriteQueryTimeScale(createTableQuery, []interface{}{})

		log.SetOutput(os.Stdout)
		parameters := make([][]interface{}, len(values))
		for i := range parameters {
			parameters[i] = make([]interface{}, 0)
		}

		//log.Printf("\n tablename: %v \r", tablename)
		//log.Printf("\n time-series entries: %v \r", values)

		insertQuery := fmt.Sprintf("INSERT INTO %v (time, timestamps, value) VALUES ($1, false, $2);", tablename)
		for i, timeSeriesEntry := range values {
			parameters[i] = append(parameters[i], timeSeriesEntry["Start"], timeSeriesEntry["Value"])
		}

		databaseapi.WriteSameQueryMultipleTimeScale(insertQuery, parameters)
	}

}
