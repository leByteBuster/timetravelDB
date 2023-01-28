# setup
- run the neo4j docker image: 
  `docker run --publish=7474:7474 --publish=7687:7687 --volume=$HOME/neo4j/data:/data neo4j`
- run the timescaledb docker image: 
  `docker run -d --name timescaledb -p 5432:5432 -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg14-latest`
- install peg/leg:
  - necessary dependency for lib 
  - download zip: https://github.com/gpakosz/peg 
  - follow instructions in README.txt
- install libcypher-parser
  - https://github.com/cleishm/libcypher-parser/issues?q=is%3Aissue+is%3Aclosed


setup libcypher-parser:
  - follow install instructions 
  - ?: 
    in ~/.bashrc add: 
    `export LD_LIBRARY_PATH=/usr/local/lib`
  - export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig (this is only
    for the current shell session. If you want to make it permanent, add it to your ~/.bashrc file ?)
  - pkg-config --libs cypher-parser


# access the database manually:
- docker exec -it timescaledb psql -U postgres  


// TODO: irgendwie muss ich die Daten abbilden 


# notes postgres
- DONT forget the ";" for SQL commands or they wont work 
- possible to define functions inside postgresql  
- Grants: access rights
- drop all tables: 
  `DROP SCHEMA public CASCADE;
   CREATE SCHEMA public;
   GRANT ALL ON SCHEMA public TO postgres;
   GRANT ALL ON SCHEMA public TO public;`

# notes neo4j 
- delete all data: MATCH (n) DETACH DELETE n 

# queries to implement:
- question: hide shallow or make it expilicit. ?
  - hiding it would be better 
  - hide it where possible
  - use it where necessary
- Shallow 
  - get all nodes (of a type)
  - get property names
- Not Shallow
  - path queries, other graph pattern queries
    - not so important maybe because nothing changes
      -> not so sure about this. the graph becomes lighter so 
         it might be faster   
  - get all properties of a node
  - get one property of a node
  - get all nodes (of a type)
  - all operations which contain property values
    - average, sum, count, min, max..