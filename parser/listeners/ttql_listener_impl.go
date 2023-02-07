package listeners

import (
	"fmt"

	tti "github.com/LexaTRex/timetravelDB/parser/ttql_interface"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type TreeShapeListener struct {
	*tti.BaseTTQLListener
}

type PropertyLookupListener struct {
	*tti.BaseTTQLListener
	Insights []PropertyClauseInsights
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func NewPropertyLookupListener() *PropertyLookupListener {
	return new(PropertyLookupListener)
}

func (treeShapeListener *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

type PropertyClauseInsights struct {
	PropertyLookupClause string
	IsWhere              bool
	IsReturn             bool
	IsValid              bool
}

func (pL *PropertyLookupListener) EnterOC_PropertyLookup(ctx *tti.OC_PropertyLookupContext) {
	fmt.Println(ctx.GetText())
	parent := ctx.GetParent()
	propertyLookupClause := ctx.GetText()
	isWhere := false
	isReturn := false
	for parent != nil && !(isWhere || isReturn) {
		if _, ok := parent.(*tti.OC_WhereContext); ok {
			isWhere = true
		} else if _, ok := parent.(*tti.OC_ReturnContext); ok {
			isReturn = true
		}
		parent = parent.GetParent()
	}
	if pL.Insights == nil {
		pL.Insights = []PropertyClauseInsights{}
	}

	pL.Insights = append(pL.Insights, PropertyClauseInsights{
		PropertyLookupClause: propertyLookupClause,
		IsWhere:              isWhere,
		IsReturn:             isReturn,
		IsValid:              isWhere || isReturn,
	})
}
