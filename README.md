Property structure out of data generator : 
node/edge {
  id: _
  property1: {
    property11: value
    property1.2: value
    property1.3: {
      property 1.3.1: value
      property 1.3.2: value
    } 
  }
}

Brainstorming

- there should be one or multiple variables which can be set and describe the 
  connection details of the server

- we need a translator to project api calls in neo4j and TimeScaleDB queries and 
  mabye does additional calculations on the data received

- queries:
  - a query should look something like the following:
    - from..to, query 
    - maybe it is translated to a method query(from, to, cypherQuery)
  - If from..to is not set in a query then the query is applied to all existing data:
    - (nil,nil,query) 

  - Following queries should be available:
    - Give me all data from..to 
    - Give me all data from..to and give me the avg()/min()/max() of:
      - all values of a property of a node/edge
      - property a node

  - Findings:
    - if I want to get the average of a property-timeseries then in 
      (from, to, cypherQuery) the cypherQuery only accesses an reference so
      cypher queries are not alloed to do aggregation calculations like sum, avg and
      things like that. They are only allowed to do:
        - path matching
        - giving back data (references)
        - ..
      the rest need to be conducted by the TimeScaleDB query. So maybe the query
      structure would be better off like:
        - (from, to, aggregations?, cypher_query?)
      The Question here is: should we even ask a cypher query from the user ? 
