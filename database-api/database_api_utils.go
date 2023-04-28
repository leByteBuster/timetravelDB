package databaseapi

func ClearTTDB() {
	ClearNeo4j()
	ClearTimescale()
}
