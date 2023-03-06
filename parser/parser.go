package parser

import (
	"errors"
	"fmt"
	"log"
	"strings"

	li "github.com/LexaTRex/timetravelDB/parser/listeners"
	tti "github.com/LexaTRex/timetravelDB/parser/ttql_interface"
	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// ParseResult holds all the information about a parsed query to build queries for Neo4j and TimescaleDB
type ParseResult struct {
	IsShallow                 bool             // is the query shallow?
	ContainsPropertyLookup    bool             // contains any property lookup?
	ContainsOnlyNullPredicate bool             // if property lookups - do any of them have a NullPredicates suffix ?
	From                      string           // start time
	To                        string           // end time
	MatchClause               string           // MATCH clause as string
	WhereClause               string           // WHERE clause as string
	ReturnClause              string           // WHERE clause
	GraphElements             li.GraphElements // all element variables occouring in the query
	// LookupsWhere              map[string][]string // all relevant lookups in Where (lookups that are relevant for binary querying)
	ReturnProjections    []string     // all projections in Return, used for ordering and time-series fetching
	LookupsWhereRelevant []LookupInfo // holds all relevant lookups (like above) but with additional information which is relevant for comparisons and information about NullPredicates. relevant: relevant for binary querying
	// Null Predicate lookups are only relevant if they occur in actual comparisons (a.prop IS NOT NULL > 20)
	LookupsReturn          map[string][]string                                                // contains all relevant lookups in Return (lookups that do not have a NullPredicate appendix - but we don't allow this in the RETURN clause yet anyways)
	PropertyClauseInsights map[*tti.OC_ComparisonExpressionContext][]li.PropertyClauseInsight // insights of Comparison expressions / Property Clauses
}

// ParseQuery parses a ttql query and returns a ParseResult which contains the relevant information to build Neo4j and TimecaleDB queries
func ParseQuery(query string) (ParseResult, error) {

	qS := antlr.NewInputStream(query)
	lexer := tti.NewTTQLLexer(qS)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := tti.NewTTQLParser(tokens)

	// remove default error listener and add custom error listener which aggregates all errors
	parser.RemoveErrorListeners()
	errorListener := li.NewErrorListener()
	parser.AddErrorListener(errorListener)

	// parse ttql query
	treectx := parser.TtQL()

	// retrieve parsing errors
	if errorListener.Errors != nil {
		var sb strings.Builder
		for _, error := range errorListener.Errors {
			sb.WriteString(error)
			sb.WriteString("\n")
		}
		return ParseResult{}, errors.New(sb.String())
	}

	listener := li.NewTtqlTreeListener()

	utils.Debug()
	utils.Debug("............................................")
	utils.Debug()
	utils.Debugf("Tree: %v", treectx.GetText())
	utils.Debug("Query is valid.")
	utils.Debug()
	utils.Debug("............................................")
	utils.Debug()

	antlr.ParseTreeWalkerDefault.Walk(listener, treectx)
	parseResult := aggregateParsingInfo(listener)

	return parseResult, nil
}

// this code is outsourced because its still pretty messy because of all the printin. As soon as everything is stable clean this function up
// maybe include error checks here
func aggregateParsingInfo(listener *li.TtqlTreeListener) ParseResult {
	fmt.Printf("\nTimeClauseInsights: \n from: %v\nto: %v\nisShallow: %v", listener.TimePeriod.From,
		listener.TimePeriod.To, listener.IsShallow)

	utils.Debug()
	utils.Debug("............................................")
	utils.Debug()
	utils.Debugf("\nMatchClause variables: %v\nWhereClause variables: %v\nReturnClause variables: %v",
		listener.GraphElements.MatchGraphElements, listener.GraphElements.WhereGraphElements, listener.GraphElements.ReturnGraphElements)
	utils.Debug()
	utils.Debug("............................................")
	utils.Debug()

	propertyClauseInsights := listener.Insights

	utils.Debug("............................................")

	containsPropertyLookup := false
	containsOnlyNullPredicate := true

	lookupsWhere := map[string][]string{}
	lookupsWhereRelevant := []LookupInfo{}
	lookupsReturn := map[string][]string{}
	for comparisonCtx, listOfInsights := range propertyClauseInsights {

		lookupInfo, err := GetRelevantLookupInfoWhere(comparisonCtx.GetText(), listOfInsights)
		if err != nil {
			log.Printf("error - retrieving relevant lookup info: %v", err)
		}

		if lookupInfo != (LookupInfo{}) {
			lookupsWhereRelevant = append(lookupsWhereRelevant, lookupInfo)
		}

		for _, insight := range listOfInsights {
			insightClause := comparisonCtx.GetText() // this should be the part of the string to be cut out
			insightClauseLookup := insight.PropertyLookupClause
			field := insight.Element
			propKeys := insight.PropertyKey
			labels := insight.Labels
			compareOp := insight.CompareOperator

			isWhere := insight.IsWhere
			isReturn := insight.IsReturn
			isComparison := insight.IsComparison
			isPartialComparison := insight.IsPartialComparison
			isPartOfActualComparison := insight.IsPartOfActualComparison
			isPropertyLookup := insight.IsPropertyLookup
			isAppendixOfNullPredicate := insight.IsAppendixOfNullPredicate

			if insight.IsPropertyLookup {
				containsPropertyLookup = true

				// collect property lookups that are relevant for binary databae fetching (neo4j, timescaleDB)
				if !insight.IsAppendixOfNullPredicate {
					containsOnlyNullPredicate = false
					if isWhere {
						lookupsWhere[insight.Element] = append(lookupsWhere[insight.Element], insight.PropertyKey)
					} else if isReturn {
						lookupsReturn[insight.Element] = append(lookupsReturn[insight.Element], insight.PropertyKey)
					}
				}
			}

			isValid := insight.IsValid

			utils.Debugf("\nComparisonWithPropertyLookupQuery: %v\nPropertyLookupinsight: %v \ncomparisonCtx: %v \nfield: %v \npropKeys: %v \nlabels: %v \ncompareOp: %v", insightClause,
				insightClauseLookup, comparisonCtx, field, propKeys, labels, compareOp)

			// print all of the insight insights
			utils.Debugf("\nIsWhere: %v	\nIsReturn: %v	\nIsComparison: %v	\nIsPartialComparison: %v	\nIsPartOfActualComparison: %v	\nIsPropertyLookup: %v \nIsAppendixOfNullPredicate: %v	\nIsValid: %v",
				isWhere, isReturn, isComparison, isPartialComparison, isPartOfActualComparison, isPropertyLookup, isAppendixOfNullPredicate, isValid)

			utils.Debug("")
			utils.Debug("............................................")
			utils.Debug("............................................")
		}
	}

	for _, v := range listener.GraphElements.MatchGraphElements {
		if _, ok := lookupsWhere[v]; !ok {
			lookupsWhere[v] = []string{}
		}
	}

	return ParseResult{
		IsShallow:                 listener.IsShallow,
		ContainsPropertyLookup:    containsPropertyLookup,
		ContainsOnlyNullPredicate: containsOnlyNullPredicate,
		From:                      listener.TimePeriod.From,
		To:                        listener.TimePeriod.To,
		MatchClause:               listener.MatchClause,
		WhereClause:               listener.WhereClause,
		ReturnClause:              listener.ReturnClause,
		ReturnProjections:         listener.ReturnProjections,
		GraphElements:             listener.GraphElements,
		// LookupsWhere:              lookupsWhere,
		LookupsWhereRelevant:   lookupsWhereRelevant,
		LookupsReturn:          lookupsReturn,
		PropertyClauseInsights: propertyClauseInsights,
	}
}

// LookupInfo is a construct describing relevant lookups in the WHERE clause of a query
// (until now: only the case when comparisons are happening)
type LookupInfo struct {
	CompareClause             string
	ElementVariable           string
	Property                  string
	CompareOperator           string
	CompareValue              any
	LookupLeft                bool // a.prop > 5 -> true, 5 > a.prop -> false
	IsAppendixOfNullPredicate bool
	AppendixOfNullPredicate   string
}

// GetRelevantLookupInfoWhere returns the relevant lookup info of a where clause
// relevant: cannot be passed to Neo4j as it is, but needs to be transformed for binary database querying
// TODO: try to integrate this logic into the listener - this aggregation might be unnecessary complicated here
func GetRelevantLookupInfoWhere(compareClause string, insights []li.PropertyClauseInsight) (LookupInfo, error) {
	var elVar string
	var property string
	var compareOperator string // check if this is retrieved the right way in listener. Test if two symbol operators like <= are recognized correctly
	var compareValueStr string
	var compareValue any
	var lookupLeft bool
	var isAppendixOfNullPredicate bool
	var appendixOfNullPredicate string

	switch len(insights) {
	case 0:
		return LookupInfo{}, errors.New("no insights found for comparison. should be impossible if comparison is in list")
	case 1:
		if !insights[0].IsAppendixOfNullPredicate && insights[0].IsWhere {
			return LookupInfo{}, errors.New("single lookups withouth appendix of null predicate (IS NULL / IS NOT NULL) only allowed in return")
		}
		return LookupInfo{}, nil
	// in this case it should be a comparison like "a.prop > 3"
	case 2:
		insightLeft := insights[0]
		insightRight := insights[1]
		if !insightLeft.IsWhere || !insightRight.IsWhere {
			return LookupInfo{}, errors.New("comparison not in WHERE clause")
		}
		if insightLeft.IsPartOfActualComparison {
			compareOperator = insightLeft.CompareOperator
		} else if insightRight.IsPartOfActualComparison {
			compareOperator = insightRight.CompareOperator
		} else {
			return LookupInfo{}, errors.New("comparison expression with two propertylabel expressions that include no partial comparison")
		}

		if insightLeft.IsPropertyLookup {
			lookupLeft = true
			elVar = insightLeft.Element
			property = insightLeft.PropertyKey
			compareValueStr = insightRight.Element
			isAppendixOfNullPredicate = insightLeft.IsAppendixOfNullPredicate
			appendixOfNullPredicate = insightLeft.AppendixOfNullPredicate
		} else if insightRight.IsPropertyLookup {
			elVar = insightRight.Element
			property = insightRight.PropertyKey
			compareValueStr = insightLeft.Element // if insight represents literal then Element is the CompareValue
			lookupLeft = false
			isAppendixOfNullPredicate = insightRight.IsAppendixOfNullPredicate
			appendixOfNullPredicate = insightRight.AppendixOfNullPredicate
		} else {
			return LookupInfo{}, nil
		}
	default:
		return LookupInfo{}, errors.New("chained comparisons are not allowed")
	}

	compareValue = utils.ConvertString(compareValueStr)

	// should only end up here if there is a comparison with a property lookup
	return LookupInfo{CompareClause: compareClause, ElementVariable: elVar, Property: property, CompareOperator: compareOperator, CompareValue: compareValue, LookupLeft: lookupLeft, IsAppendixOfNullPredicate: isAppendixOfNullPredicate, AppendixOfNullPredicate: appendixOfNullPredicate}, nil
}
