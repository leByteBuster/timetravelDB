# setup
- run the neo4j docker image: 
  `docker run --publish=7474:7474 --publish=7687:7687 --volume=$HOME/neo4j/data:/data neo4j`
- run the timescaledb docker image: 
  `docker run -d --name timescaledb -p 5432:5432 -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg14-latest`

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