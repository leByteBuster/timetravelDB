package parser

import (
	"errors"
	"fmt"
	"strings"

	li "github.com/LexaTRex/timetravelDB/parser/listeners"
	tti "github.com/LexaTRex/timetravelDB/parser/ttql_interface"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type ParseResult struct {
	IsShallow              bool
	From                   string
	To                     string
	MatchClause            string
	WhereClause            string
	ReturnClause           string
	PropertyClauseInsights []li.PropertyClauseInsight
}

func ParseQuery(query string) (ParseResult, error) {

	qS := antlr.NewInputStream(query)
	fmt.Println("Lexer Tokens: ")

	lexer := tti.NewTTQLLexer(qS)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := tti.NewTTQLParser(tokens)
	parser.RemoveErrorListeners()
	errorListener := li.NewErrorListener()
	parser.AddErrorListener(errorListener)

	// interrupt parsing as soon as an error occours
	//parser.SetErrorHandler(antlr.NewBailErIrorStrategy())

	treectx := parser.TtQL()

	if errorListener.Errors != nil {
		var sb strings.Builder
		for _, error := range errorListener.Errors {
			sb.WriteString(error)
			sb.WriteString("\n")
		}
		return ParseResult{}, errors.New(sb.String())
	}

	// initialize listeners
	timeListener := li.NewTimeShallowListener()
	propertyListener := li.NewPropertyOrLabelsExpressionListener()
	whereListener := li.NewWhereListener()
	returnListener := li.NewReturnListener()

	// HOW TO GET THE ERRORS ?????
	//		parser.GetErrorHandler()
	//		parser.GetErrorListenerDispatch()
	//		parser._SyntaxErrors
	//
	fmt.Println()
	fmt.Println("............................................")
	fmt.Println()
	fmt.Printf("Tree: %v", treectx.GetText())
	fmt.Println("Query is valid.")
	fmt.Println()
	fmt.Println("............................................")
	fmt.Println()
	antlr.ParseTreeWalkerDefault.Walk(timeListener, treectx)
	fmt.Printf("\nTimeClauseInsights: \n from: %v\nto: %v\nisShallow: %v", timeListener.TimePeriod.From, timeListener.TimePeriod.To, timeListener.IsShallow)
	fmt.Println()
	fmt.Println("............................................")
	fmt.Println()
	antlr.ParseTreeWalkerDefault.Walk(whereListener, treectx)
	antlr.ParseTreeWalkerDefault.Walk(returnListener, treectx)
	fmt.Printf("\nWhereClause: %v\nReturnClause: %v", whereListener.WhereClause, returnListener.ReturnClause)
	fmt.Println()
	fmt.Println("............................................")
	fmt.Println()

	antlr.ParseTreeWalkerDefault.Walk(propertyListener, treectx)
	propertyClauseInsights := propertyListener.Insights

	fmt.Println("............................................")
	for _, subquery := range propertyClauseInsights {
		subqueryClause := subquery.PropertyLookupClause
		comparisonCtx := subquery.ComparisonContext
		field := subquery.Field
		propKeys := subquery.PropertyKeys
		labels := subquery.Labels
		compareOp := subquery.CompareOperator

		isWhere := subquery.IsWhere
		isReturn := subquery.IsReturn
		isComparison := subquery.IsComparison
		isPartialComparison := subquery.IsPartialComparison
		isPropertyLookup := subquery.IsPropertyLookup
		isValid := subquery.IsValid

		fmt.Printf("\nSubquery: %v \ncomparisonCtx: %v \nfield: %v \npropKeys: %v \nlabels: %v \ncompareOp: %v", subqueryClause, comparisonCtx, field, propKeys, labels, compareOp)

		// print all of the subquery insights
		fmt.Printf("\nIsWhere: %v	\nIsReturn: %v	\nIsComparison: %v	\nIsPartialComparison: %v	\nIsPropertyLookup: %v	\nIsValid: %v",
			isWhere, isReturn, isComparison, isPartialComparison, isPropertyLookup, isValid)

		fmt.Println("")
		fmt.Println("............................................")
		fmt.Println("............................................")
	}

	// String cypherQuery2 = "MATCH (n) WHERE n.ping > 22.33" + " RETURN n.ping, n ";
	// String cypherQuery3 = "MATCH (a)-[x]->(b) " + " RETURN a.ping, b "; // should parse
	// String cypherQuery4 = "MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b "; // should parse

	return ParseResult{
		IsShallow:              timeListener.IsShallow,
		From:                   timeListener.TimePeriod.From,
		To:                     timeListener.TimePeriod.To,
		WhereClause:            whereListener.WhereClause,
		ReturnClause:           returnListener.ReturnClause,
		PropertyClauseInsights: propertyClauseInsights,
	}, nil
}
