# running the tests 

# running docker compose

# timescaledb: 

## manually starting setting up docker testing environment

### making a backup 

### restoring the backup 
  - start container: 
    - `docker run --name test_timescaledb -e POSTGRES_PASSWORD=password -d timescale/timescaledb-ha:pg14-latest`
  - copy backup-file into container: 
    - `docker cp /path/to/backup.bak test_timescaledb:/tmp/backup.bak`
  - connect to container and restore the database: 
    - `docker exec -it test_timescaledb bash`
    - `pg_restore -U postgres -d postgres /tmp/postgres.bak`
    - `exit`


# neo4j: 

## manually starting setting up docker testing environment

### starting a container with the current used test data (already restored) - this should be standard
- the other cases only have to be considered if the test-data has been corrupted or if the test data was changed intentionally
- run container with restored database:  
  - `docker run --name testneo4j -p7474:7474 -p7687:7687 -d -v ./docker-test/neo4j/data:/data -v ./docker-test/neo4j/logs:/logs -v ./docker-test/neo4j/backups:/backups --env NEO4J_AUTH=neo4j/test neo4j`



### make neo4j dump backup: https://www.markhneedham.com/blog/2020/01/28/neo4j-database-dump-docker-container/
  - if the container has no bind volume for data, backups:
    - stop container 
    - start container (only works immediately after a container is restarted) 
    - `docker exec -it <container> neo4j-admin database dump <database> --to-path=/var/backups/ --verbose`
      - OR: - `docker exec -it <container> /bin/bash`
            - `neo4j-admin database dump <database> --to-path=/var/backups/ --verbose`
    - copy dump to location on local machine:
      `docker cp <container>:/var/backups/neo4j.dump /path/to/dump`
  - if the container has bind volume data, backups 
    - dump the file like described here: https://neo4j.com/docs/operations-manual/current/docker/maintenance/
    - then its possible to save a copy of this dump file wherever. it is in /path/to/backup which was bound to /backups

### load neo4 dump file: https://neo4j.com/docs/operations-manual/current/docker/maintenance/
  - note: this can only be done with bound volumes !
  - move the dump file to load in some folder /path/to/backup
  - run:  `docker run --interactive --tty --rm \
              --volume=/path/to/data:/data \ 
              --volume=/path/to/backups:/backups \ 
                neo4j neo4j-admin database load neo4j --from-path=/backups --verbose`
    note: the path: `/path/to/data` doesn't has to exist yet. It will be created if it does not exist and database will be loaded by using the .dump file in 
          `/path/to/backups`

create and start a new container based on the loaded backup:
  - `docker run --name testneo4j -p7474:7474 -p7687:7687 -d -v /path/to/data:/data -v ./logs:/logs -v .path/to/backups:/backups --env NEO4J_AUTH=neo4j/test neo4j`


### starting a container with the current used test data WITH NEW RESTORE from a backup dump 
- in ./test-data/neo4j-test-backup is the neo4j.dump file of the test database currently used on all our tests
- volume binding:
    - we want to bind ./docker-test/data to /data of the new container
    - we want to bind ./docker-test/backups to /backups of the new container 
- make sure the ./docker-test/data folder is empty  
- make sure the ./docker-test/backups folder is empty 
- move the test file neo4j.dump in ./test-data/neo4j-test-backup to ./docker-test/backups:
  `cp ./test-data/neo4j-test-backup/neo4j.dump ./docker-test/backups`
- restore the database: 
  - `docker run --interactive --tty --rm --volume=./docker-test/data:/data --volume=./docker-test/backups:/backups neo4j neo4j-admin database load neo4j --from-path=/backups --verbose`
- run container with restored database:  
  - `docker run --name testneo4j -p7474:7474 -p7687:7687 -d -v ./docker-test/data:/data -v ./docker-test/logs:/logs -v ./docker-test/backups:/backups --env NEO4J_AUTH=neo4j/test neo4j`


note THIS WORKS AS WELL 
docker run \
    --name testneo4j \
    -p7474:7474 -p7687:7687 \
    -d \
    -v ./data:/data \
    -v ./backups:/backups \
    neo4j:latest

