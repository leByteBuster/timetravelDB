#!/bin/bash

SCRIPT=$(realpath "$0")
HOME=$(dirname "$SCRIPT")

echo $HOME

# stop all running docker containers to avoid conflicts 
docker stop $(docker ps -aq)

# cleanup before (if docker-compose was run manually there might be data, in this case restore is not possible)
sudo rm -rf ./docker-test/neo4j/backups/*
sudo rm -rf ./docker-test/neo4j/data/*

# prepare testing data  
cp $HOME/test-data/neo4j_test_backup/neo4j.dump $HOME/docker-test/neo4j/backups/
cp $HOME/test-data/timescaledb_test_backup/postgres.bak $HOME/docker-test/timescaledb/backups/

# restore neo4j data (was not able to get this to work with docker-compose yet)k
docker run --interactive --tty --rm --volume=$HOME/docker-test/neo4j/data:/data --volume=$HOME/docker-test/neo4j/backups:/backups neo4j neo4j-admin database load neo4j --from-path=/backups --verbose

# prepare testing envionment 
docker-compose up -d

# wait for testing envionment to be ready 
dockerize -wait tcp://127.0.0.1:7687 -timeout 10s
dockerize -wait tcp://127.0.0.1:5432 -timeout 10s

# check if neo4j database is ready
docker exec test_neo4j sh -c "while neo4j-admin database info | grep -q 'Database in use:.*false'; do echo 'neo4j database not ready yet'; sleep 2; done"
echo "neo4j database ready"

# Run Golang tests
go test -v $HOME/api -run TestShallowQueries -test.v #> >(tee -a /dev/stdout) 2> >(tee -a /dev/stderr >&2)
gotest1=$?
go test -v $HOME/api -run TestNonShallowQueries -test.v #> >(tee -a /dev/stdout) 2> >(tee -a /dev/stderr >&2)
gotest2=$?


Stop Docker Compose
docker-compose down

# cleanup  
sudo rm -rf $HOME/docker-test/neo4j/backups/*
sudo rm -rf $HOME/docker-test/neo4j/data/*
sudo rm -rf $HOME/docker-test/timescaledb/backups/*
sudo rm -rf $HOME/docker-test/timescaledb/data/*

# if ( $gotest1 -eq 0 && $gotest2 -eq 0 ); then
#     echo "Tests passed"
#     exit 0
# else
#     echo "Tests failed"
#     exit 1
# fi
