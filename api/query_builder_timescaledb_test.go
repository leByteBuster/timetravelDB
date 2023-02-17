package api

import (
	"testing"
)

func TestBuildQueryString(t *testing.T) {
	// "1"" is in '' because we retrieve it like that from the parsing tree
	query, err := buildQueryString("x", "y", "", "=", "1", true, []string{"table1"})
	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}
	expected := "SELECT * FROM (SELECT time, timestamps, value FROM table1 WHERE time >= 'x' AND time < 'y' AND value = 1) genericAliasName;"
	if query != expected {
		t.Fatalf("\n    Expected: %v\n    Got: %v", expected, query)
	}
	// "'trichter'" is in '' because we retrieve it like that from the parsing tree
	query, err = buildQueryString("x", "y", "", "=", "'trichter'", true, []string{"table1"})
	if err != nil {
		t.Fatalf("Error building query: %v", err)
	}
	expected = "SELECT * FROM (SELECT time, timestamps, value FROM table1 WHERE time >= 'x' AND time < 'y' AND value = 'trichter') genericAliasName;"
	if query != expected {
		t.Fatalf("\n    Expected: %v\n    Got: %v", expected, query)
	}
}
