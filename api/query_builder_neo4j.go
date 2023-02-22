package api

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/LexaTRex/timetravelDB/parser"
	"github.com/LexaTRex/timetravelDB/utils"
)

// manipulate the WHERE clause sent to neo4j to prefilter the result. If there are property lookups
// that compare to a value in the form of "a.prop > 2" then for performance inhancement
// it is replaced by "a.prop IS NOT NULL". So if it does not excist we do not get a pattern back where
// we have to try tro retrieve the non existing property values for which we want to compare to 2 (in the example)
func manipulateWhereClause(queryInfo parser.ParseResult, whereClause string) (string, error) {
	// note: all Lookups in LookupsWhere are property lookups that are not of the type IS NULL / IS NOT NULL

	for compCtx, insights := range queryInfo.PropertyClauseInsights {
		log.Printf("\nInsights for compareContext: %v \ninsights: %+v\n", compCtx.GetText(), insights)
		// iterate over the insights of eache ComparisonExpressions. Each insight contains information about one
		// PropertyOrLabelExpressions (x.prop/x) inside the ComparisonExpression. If one insight isPartialComparison==true
		// then this PropertyOrLabelExpression is part of a comparison with a value (even though the other insight
		// won't be isPartialComparison==true). We do not alloed chained comparisons (yet).
		// I need to find out which of the PropertyLabelExpressions (insights) contains the property lookup
		// because i need the variable and the lookup property for manipulating the WHERE clause.
		// example:
		// allowed: WHERE a.prop1 > 1 AND a.prop2 < 2
		// not allowed: WHERE a.prop1 > 1 > a.prop3
		// not allowed: WHERE a.prop1 >  a.prop3      (this does not make sense because properties are time-series)

		var isRelevant = false
		var el string
		var lookup string

		if len(insights) > 2 {
			return "", errors.New("chained comparisons are not allowed")
		}
		for i, insight := range insights {
			log.Printf("\ni: %v \ninsight: %+v\n", i, insight)
			isRelevant = !insight.IsAppendixOfNullPredicate && insight.IsWhere && insight.IsPartialComparison
			if isRelevant {
				if insight.IsPropertyLookup {
					el = insight.Element
					lookup = insight.PropertyKey
					break
				} else {
					idx := (i + 1) % 2
					// attention ! this is error prone ! only works if if there are only two Insights which are
					// PartialComparisons. Make this more robust when extending the functionality
					if !insights[idx].IsPropertyLookup {
						return "", errors.New("error partial comparison: one of both must be a property lookup")
					}
					el = insights[idx].Element
					lookup = insights[idx].PropertyKey
					break
				}
			}
		}
		if isRelevant {
			log.Println("isRelevant!")
			orig := compCtx.GetText()
			var repl strings.Builder
			repl.WriteString(el)
			repl.WriteString(".")
			repl.WriteString(lookup)
			repl.WriteString(" IS NOT NULL")
			whereClause = strings.ReplaceAll(whereClause, orig, repl.String())
		}
	}
	// andClauses := strings.Split(queryInfo.WhereClause, "AND")
	// orClauses := strings.Split(queryInfo.WhereClause, "OR")

	return whereClause, nil
}

func buildTmpWhereClause(from, to, whereClause string, matchElements []string) string {

	var sb strings.Builder
	if strings.TrimSpace(whereClause) == "" {
		sb.WriteString(" WHERE")
	} else {
		sb.WriteString(" ")
		sb.WriteString(whereClause)
		sb.WriteString(" AND")
	}
	for i, el := range matchElements {
		sb.WriteString(" ")
		sb.WriteString(el)
		sb.WriteString(".")
		sb.WriteString("end >= '") // TODO
		sb.WriteString(from)
		sb.WriteString("' AND ")
		sb.WriteString(el)
		sb.WriteString(".")
		sb.WriteString("start <= '")
		sb.WriteString(to)
		sb.WriteString("' ")
		if i < len(matchElements)-1 {
			sb.WriteString("AND")
		}
	}
	return sb.String()
}

func buildReturnClause(whereLookups []parser.LookupInfo, returnElements []string) string {

	var sb strings.Builder
	returnVariables := utils.NewSet()

	// if there are no return elements specified, return all variables
	if len(returnElements) == 0 {
		return " RETURN *"
	}
	for _, lookup := range whereLookups {
		// some lookupInfos do not contain an elementVar (if they represent the literal side of the lookup)
		if strings.Trim(lookup.ElementVariable, " ") != "" {
			returnVariables.Add(lookup.ElementVariable)
		}
	}
	for _, element := range returnElements {
		returnVariables.Add(element)
	}

	sb.WriteString("RETURN ")

	varList := returnVariables.GetElements()
	size := len(varList)

	fmt.Printf("\n\n    returnElements: %+v    \n\n", returnElements)
	fmt.Printf("\n\n    whereLookups: %+v    \n\n", whereLookups)
	fmt.Printf("\n\n    returnVariables: %+v    \n\n", returnVariables)
	fmt.Printf("\n\n    varList: %+v    \n\n", varList)

	for i, el := range varList {
		sb.WriteString(el)
		if i < size-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func buildFinalQuery(matchClause, whereClause, returnClause string) string {
	var sb strings.Builder
	sb.WriteString(matchClause)
	sb.WriteString(whereClause)
	sb.WriteString(returnClause)

	tmpQuery := sb.String()
	fmt.Println(tmpQuery)
	return tmpQuery
}
