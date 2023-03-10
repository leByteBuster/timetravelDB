# HOW TO USE TimeTravelDB

## setup
- run the neo4j docker image: 
  `docker run --publish=7474:7474 --publish=7687:7687 --volume=$HOME/neo4j/data:/data neo4j`
- install and run the timescaledb docker image (see: https://docs.timescale.com/install/latest/installation-docker/): 
  `docker pull timescale/timescaledb-ha:pg14-latest`
  `docker run -d --name timescaledb -p 5432:5432 -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg14-latest`


## TTQL CLI
  - run go run main/main.go
  - TTQL query in the form of: FROM x TO y MATCH (n)-[r]->(s) RETURN n,r,s
  - help: prints some help text
  - quit/exit: quits the program
  - note: if the program exits unexpectedly then sometimes the terminal settings are not restored correctly 
    and the typed text is not visible anymore. To fix this issue type: `reset` in the terminal and hit enter
    and restart the program


# additional useful intel
 
## access the database manually:
- neo4j: http://localhost:7474/browser/
- timescaledb: docker exec -it timescaledb psql -U postgres  


// TODO: irgendwie muss ich die Daten abbilden 


## notes postgres
- DONT forget the ";" for SQL commands or they wont work 
- possible to define functions inside postgresql  
- Grants: access rights
- drop all tables: 
  `DROP SCHEMA public CASCADE;
   CREATE SCHEMA public;
   GRANT ALL ON SCHEMA public TO postgres;
   GRANT ALL ON SCHEMA public TO public;`
- list all tables:
   SELECT table_name
   FROM information_schema.tables
   WHERE table_schema='public'
   AND table_type='BASE TABLE'; 
    

## notes neo4j 
- delete all data: MATCH (n) DETACH DELETE n 


!!!! Need to merge the READMEs

get back brocken docker container (create new with old database):
    docker run \
    --publish=7474:7474 --publish=7687:7687 \
    --volume=$HOME/neo4j/data:/data \
    --volume=$HOME/neo4j/logs:/logs \
    neo4j:latest



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


# Developer Info

## Testing


### Run Integration Tests

- be aware that running tests including a database setup expect a clean setup with a restore of the test databases. If the tests
  are run manually make sure to clean up the shared volumes of the docker volumes before running the tests. This can be omitted 
  when running go test TestIntegration which will trigger a cleanup and starup of the containers with freshly restored test-databases.