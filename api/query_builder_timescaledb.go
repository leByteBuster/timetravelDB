package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/LexaTRex/timetravelDB/utils"
)

func uuidToTablename(uuid string) string {
	var builder strings.Builder
	builder.WriteString("ts_")
	builder.WriteString(strings.Replace(uuid, "-", "_", -1))
	return builder.String()
}

// get properties / multiple time-series and apply aggrergation on it if not empty
func buildQueryString(from, to, aggr string, cmpOp string, cmpVal any, lookupLeft bool, tables []string) (string, error) {

	var builder strings.Builder
	var err error

	if !lookupLeft && cmpOp != "" && cmpVal != "" {
		if cmpOp, err = invertCmpOperatro(cmpOp); err != nil {
			return "", err
		}
	}

	// build query string
	builder.WriteString("SELECT ")
	// builder.WriteString(aggr) cannot do this here because then i only have one return value
	// have to call a different subfunction aggrTimescale or something
	builder.WriteString("*")
	builder.WriteString(" FROM (")
	for i, tablename := range tables {
		if i > 0 {
			builder.WriteString(" UNION ALL ")
		}
		builder.WriteString("SELECT ")
		builder.WriteString("time, timestamps, value FROM ")
		builder.WriteString(tablename)
		builder.WriteString(" WHERE time >= ")
		builder.WriteString("'")
		builder.WriteString(from)
		builder.WriteString("'")
		builder.WriteString(" AND time < ")

		//TODO: CHANGE THIS TO DATETIME
		builder.WriteString("'")
		builder.WriteString(to)
		builder.WriteString("'")
		if cmpOp != "" && cmpVal != "" {
			builder.WriteString(" AND ")
			builder.WriteString("value ")
			builder.WriteString(cmpOp)
			builder.WriteString(" ")
			switch t := cmpVal.(type) {
			case string:
				builder.WriteString(t) // no quotes needed because already around value (because retrieved like this from parsing tree: 'val')
			default:
				builder.WriteString(strings.Trim(utils.AnyToString(cmpVal), "'"))
			}
		}
	}
	builder.WriteString(") genericAliasName;")
	fmt.Println(builder.String())

	return builder.String(), nil
}

// get properties / multiple time-series and apply aggrergation on it if not empty
func buildQueryStringCmpExists(from, to, aggr string, cmpOp string, cmpVal any, lookupLeft bool, table string) (string, error) {

	var builder strings.Builder
	var err error

	if !lookupLeft && cmpOp != "" && cmpVal != "" {
		if cmpOp, err = invertCmpOperatro(cmpOp); err != nil {
			return "", err
		}
	}

	// build query string
	// builder.WriteString(aggr) cannot do this here because then i only have one return value
	// have to call a different subfunction aggrTimescale or something
	builder.WriteString("SELECT EXISTS (SELECT 1 FROM ")
	builder.WriteString(table)
	builder.WriteString(" WHERE time >= ")
	builder.WriteString("'")
	builder.WriteString(from)
	builder.WriteString("'")
	builder.WriteString(" AND time < ")

	//TODO: CHANGE THIS TO DATETIME
	builder.WriteString("'")
	builder.WriteString(to)
	builder.WriteString("'")
	if cmpOp != "" && cmpVal != "" {
		builder.WriteString(" AND ")
		builder.WriteString("value ")
		builder.WriteString(cmpOp)
		builder.WriteString(" ")
		switch t := cmpVal.(type) {
		case string:
			builder.WriteString(t) // no quotes needed because already around value (because retrieved like this from parsing tree: 'val')
		default:
			builder.WriteString(strings.Trim(utils.AnyToString(cmpVal), "'"))
		}
	}
	builder.WriteString(");")
	fmt.Println(builder.String())

	return builder.String(), nil
}

func invertCmpOperatro(cmpOp string) (string, error) {
	switch cmpOp {
	case ">":
		return "<", nil
	case "<":
		return ">", nil
	case ">=":
		return "<=", nil
	case "<=":
		return ">=", nil
	case "=":
		return "=", nil
	case "!=":
		return "!=", nil
	default:
		log.Printf("\nunsupported compare operator [%s]\n", cmpOp)
		return "", fmt.Errorf("unsupported compare operator %s", cmpOp)
	}
}
