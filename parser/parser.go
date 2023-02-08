package parser

import (
	"fmt"

	li "github.com/LexaTRex/timetravelDB/parser/listeners"
	tti "github.com/LexaTRex/timetravelDB/parser/ttql_interface"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

func ParseTest() {

	//String cypherQuery1 = "MATCH (a)-[x]->(b) WHERE a.ping > 22.33" + "RETURN a.ping, b";  // should not parse From, TO missing
	//String ttQuery2 = "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33" + "RETURN n.ping, n ";
	//String ttQuery3 = "FROM 2023-02-03T12:34:39Z TO 2023-02-03T14:34:39.2222Z SHALLOW MATCH (a)-[x]->(b) " + "RETURN a.ping, b "; // should parse
	//String ttQuery4 = "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b "; // should parse
	//String ttQuery = "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a "; // should parse
	//String ttQuery6 = "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a) WHERE a.ping > 22" + " RETURN a "; // should parse
	//String ttQuery5 = "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b "; // should parse

	ttQuery4 := "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b " // should parse
	ttQuery5 := "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a "         // should parse
	ttQuery6 := "FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a) WHERE a.ping > 22" + " RETURN a "                  // should parse
	//ttQuery5 := "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b " // should parse

	queries := []string{ttQuery4, ttQuery5, ttQuery6}

	for _, query := range queries {

		qS := antlr.NewInputStream(query)
		fmt.Println("Lexer Tokens: ")

		lexer := tti.NewTTQLLexer(qS)
		tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
		parser := tti.NewTTQLParser(tokens)
		parser.SetErrorHandler(antlr.NewDefaultErrorStrategy())

		treectx := parser.TtQL()

		// HOW TO GET THE ERRORS ?????
		//		parser.GetErrorHandler()
		//		parser.GetErrorListenerDispatch()
		//		parser._SyntaxErrors
		//
		fmt.Println()
		fmt.Println("............................................")
		fmt.Println()

		antlr.ParseTreeWalkerDefault.Walk(li.NewTreeShapeListener(), treectx)

		fmt.Println()
		fmt.Println("............................................")
		fmt.Println()
		fmt.Printf("Tree: %v", treectx.GetText())

		fmt.Println("Query is valid.")

		fmt.Println()
		fmt.Println("............................................")
		fmt.Println()
		timeListener := li.NewTimeClauseListener()
		antlr.ParseTreeWalkerDefault.Walk(timeListener, treectx)
		fmt.Printf("\nTimeClauseInsights: \n from: %v\nto: %v", timeListener.TimePeriod.From, timeListener.TimePeriod.To)
		fmt.Println()
		fmt.Println("............................................")
		fmt.Println()
		//propertyLookupListener := li.NewPropertyLookupListener()
		propertyListener := li.NewPropertyOrLabelsExpressionListener()

		antlr.ParseTreeWalkerDefault.Walk(propertyListener, treectx)
		propertyClauseInsights := propertyListener.Insights
		fmt.Printf("PropertyClauseInsights: %v", propertyClauseInsights)

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
	}

	// String cypherQuery2 = "MATCH (n) WHERE n.ping > 22.33" + " RETURN n.ping, n ";
	// String cypherQuery3 = "MATCH (a)-[x]->(b) " + " RETURN a.ping, b "; // should parse
	// String cypherQuery4 = "MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b "; // should parse

}
