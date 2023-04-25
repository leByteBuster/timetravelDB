# BENCHMARKING

## Benchmark Queries 

*graph alternations*

100000 nodes 100000 edges

10000 nodes 10000 edges

1000 nodes 1000 edges

100 nodes 100 edges 

10 nodes 10 edges 




*property alternations*
10 properties per graph element



*time series length alternations (values per property)*

10 properties per graph element

1000 properties per graph element

10000 properties per graph element

100000 properties per graph element

1000000 properties per graph element



*ON*: only Neo4j queries
*TTDB*: TTDB queries 


// Querying Time Series Data



ON: 
  Query a single node
`MATCH (n) WHERE ID(n) = 10 AND n.start = X AND n.start = Y RETURN n`

TTDB
  Query a single node shallow
  "FROM X TO Y SHALLOW MATCH (n) WHERE ID(n) = 10 RETURN n"
  Query a single node
  "FROM X TO Y MATCH (n) WHERE ID(n) = 10 RETURN n"

ON: 
  Query a time series property of a single node (by ID)
    MATCH (n) WHERE ID(n) = 10 n.start = X AND n.start = Y RETURN n.ts_property*"
TTDB
  Query a time series property of a single node (by ID) 
    FROM X TO Y MATCH (n) WHERE ID(n) = 10 RETURN n.ts_property"

ON: 
  Query a time series property of all nodes (that have this property)
    MATCH (n) WHERE n.start = X AND n.start = Y RETURN n.ts_property*"
TTDB
  Query a time series property of all nodes (that have this property)
    FROM X TO Y MATCH (n) RETURN n.ts_property"

ON: 
  Query all time series properties of a single node 
  "MATCH (n) WHERE ID(n) = 10 AND n.start = X AND n.start = Y RETURN properties(n)"
TTDB
  Query all time series properties of a single node 
  "FROM X TO Y MATCH (n) WHERE ID(n) = 10 RETURN n.prop1, n.prop2, ..."


ON: 
  Query all time series properties of all nodes 
  "MATCH (n) WHERE n.start = X AND n.start = Y RETURN properties(n)"
TTDB
  Query all time series properties of all nodes 
  "FROM X TO Y MATCH (n) RETURN n.prop1, n.prop2, ..."



// Querying Time Series Data: the ANY OPERATOR

ON: 
  Query a time series property of all nodes (that have this property)
    MATCH (n) WHERE n.start = X AND n.start = Y AND (n.ts_property... > 20 OR n.ts_property.. > 20 OR ... ) RETURN n.ts_property*"
TTDB
  Query a time series property of all nodes (that have this property)
    FROM X TO Y MATCH (n) AND ANY(n.ts_property > 20) RETURN n.ts_property"

