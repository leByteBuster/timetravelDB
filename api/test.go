package api

// package test
//
// import "testing"
//
// // test the ParseQuery function
// func TestParseQuery(t *testing.T) {
// 	// test a valid query
// 	_, _, _, _, err := ParseQuery("FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z MATCH (n) RETURN n")
// 	if err != nil {
// 		t.Errorf("ParseQuery(\"FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z MATCH (n) RETURN n\") failed: %v", err)
// 	}
// 	// test a query with an invalid time
// 	_, _, _, _, err = ParseQuery("FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z MATCH (n) RETURN n")
// 	if err == nil {
// 		t.Errorf("ParseQuery(\"FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z MATCH (n) RETURN n\") failed: %v", err)
// 	}
// 	// test a query with an invalid cypher string
// 	_, _, _, _, err = ParseQuery("FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z MATCH (n) RETURN n")
// 	if err == nil {
// 		t.Errorf("ParseQuery(\"FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z MATCH (n) RETURN n\") failed: %v", err)
// 	}
// }
//

// func TestParseQuery(t *testing.T) {
//
// 	// test 1: valid query
// 	if _, _, _, _, err := ParseQuery("FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z SHALLOW MATCH (n) RETURN n"); err != nil {
// 		t.Error("TestParseQuery: valid query failed")
// 	}
// 	// test 2: invalid query
// 	if _, _, _, _, err := ParseQuery("FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z SHALLOW MATCH (n)"); err == nil {
// 		t.Error("TestParseQuery: invalid query failed")
// 	}
// 	// test 3: invalid query
// 	if _, _, _, _, err := ParseQuery("FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z MATCH (n)"); err == nil {
// 		t.Error("TestParseQuery: invalid query failed")
// 	}
// 	// test 4: invalid query
// 	if _, _, _, _, err := ParseQuery("FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z SHALLOW (n)"); err == nil {
// 		t.Error("TestParseQuery: invalid query failed")
// 	}
// 	// test 5: invalid query
// 	if _, _, _, _, err := ParseQuery("FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z MATCH (n) RETURN n"); err == nil {
// 		t.Error("TestParseQuery: invalid query failed")
// 	}
// 	// test 6: invalid query
// 	if _, _, _, _, err := ParseQuery("FROM 2018-01-01T00:00:00Z TO 2018-01-01T00:00:00Z SHALLOW MATCH (n) RETURN
//
