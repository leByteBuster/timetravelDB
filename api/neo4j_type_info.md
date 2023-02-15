			// oC_ComparisonExpression
			// :  oC_StringListNullPredicateExpression ( SP? oC_PartialComparisonExpression )* ;

			// oC_PartialComparisonExpression
			// 													 :  ( '=' SP? oC_StringListNullPredicateExpression )
			// 															 | ( '<>' SP? oC_StringListNullPredicateExpression )
			// 															 | ( '<' SP? oC_StringListNullPredicateExpression )
			// 															 | ( '>' SP? oC_StringListNullPredicateExpression )
			// 															 | ( '<=' SP? oC_StringListNullPredicateExpression )
			// 															 | ( '>=' SP? oC_StringListNullPredicateExpression )




type (
	Point2D       = dbtype.Point2D
	Point3D       = dbtype.Point3D
	Date          = dbtype.Date
	LocalTime     = dbtype.LocalTime
	LocalDateTime = dbtype.LocalDateTime
	Time          = dbtype.Time
	OffsetTime    = dbtype.Time
	Duration      = dbtype.Duration
	Entity        = dbtype.Entity
	Node          = dbtype.Node
	Relationship  = dbtype.Relationship
	Path          = dbtype.Path
	Record        = db.Record
	InvalidValue  = dbtype.InvalidValue
)

	// 	// Deprecated: Id is deprecated and will be removed in 6.0. Use ElementId instead.
	// 	Id        int64          // Id of this Node.
	// 	ElementId string         // ElementId of this Node.
	// 	Labels    []string       // Labels attached to this Node.
	// 	Props     map[string]any // Properties of this Node.
	// }
	//

	// // Node represents a node in the neo4j graph database
	// type Node struct {
	// 	// Deprecated: Id is deprecated and will be removed in 6.0. Use ElementId instead.
	// 	Id        int64          // Id of this Node.
	// 	ElementId string         // ElementId of this Node.
	// 	Labels    []string       // Labels attached to this Node.
	// 	Props     map[string]any // Properties of this Node.
	// }

// For the rest of the types look into neo4j/aliases.go