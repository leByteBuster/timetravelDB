package qpengine

import (
	"strings"

	"github.com/LexaTRex/timetravelDB/query-processor/parser"
	"github.com/LexaTRex/timetravelDB/utils"
)

// buildCondWhereClause manipulates the WHERE clause sent to neo4j to prefilter the result. If there are property lookups
// that compare to a value in the form of "a.prop > 2" then for performance inhancement
// it is replaced by "a.prop IS NOT NULL". So if it does not excist we do not get a pattern back where
// we have to try tro retrieve the non existing property values for which we want to compare to 2 (in the example)
func buildCondWhereClause(lookupsWhereRelevant []parser.LookupInfo, whereClause string) (string, error) {
	for _, lookup := range lookupsWhereRelevant {

		orig := lookup.CompareClause
		var repl strings.Builder
		repl.WriteString(lookup.ElementVariable)
		repl.WriteString(".")
		repl.WriteString(lookup.Property)
		if lookup.IsAppendixOfNullPredicate {
			repl.WriteString(" ")
			repl.WriteString(lookup.AppendixOfNullPredicate)
		} else {
			repl.WriteString(" IS NOT NULL")
			whereClause = strings.ReplaceAll(whereClause, orig, repl.String())
		}
	}
	// andClauses := strings.Split(queryInfo.WhereClause, "AND")
	// orClauses := strings.Split(queryInfo.WhereClause, "OR")

	return whereClause, nil
}

func buildTmpWhereClause(from, to, whereClause string, matchElVars []string) string {

	var sb strings.Builder
	if strings.TrimSpace(whereClause) == "" {
		sb.WriteString(" WHERE")
	} else {
		sb.WriteString(" ")
		sb.WriteString(whereClause)
		sb.WriteString(" AND")
	}
	for i, elVar := range matchElVars {

		if from != "current" && from != "CURRENT" {
			sb.WriteString(" ")
			sb.WriteString(elVar)
			sb.WriteString(".")
			sb.WriteString("end >= ")
			sb.WriteString("datetime('")
			sb.WriteString(from)
			sb.WriteString("')")
		} else {
			// if from == current then only allow elements that elVar.from == current
			sb.WriteString(" ")
			sb.WriteString(elVar)
			sb.WriteString(".")
			sb.WriteString("from = 'current' ")
		}

		// if to != current check elVar.start <= to elseif to = current then the element will start earlier (or same time) anyways
		if to != "current" && to != "CURRENT" {
			sb.WriteString(" AND ")
			sb.WriteString(elVar)
			sb.WriteString(".")
			sb.WriteString("start <= ")
			sb.WriteString("datetime('")
			sb.WriteString(to)
			sb.WriteString("') ")
		}
		if i < len(matchElVars)-1 {
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

	utils.Debugf("\n\n    returnElements: %+v    \n\n", returnElements)
	utils.Debugf("\n\n    whereLookups: %+v    \n\n", whereLookups)
	utils.Debugf("\n\n    returnVariables: %+v    \n\n", returnVariables)
	utils.Debugf("\n\n    varList: %+v    \n\n", varList)

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
	utils.Debug(tmpQuery)
	return tmpQuery
}
