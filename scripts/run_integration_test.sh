#!/bin/bash

# This file sets up an testing environment providing a docker Neo4j DMBS and a docker TimescaleDB DMBS which run databases with
# simplistic testing data and runs the Golang SHALLOW and DEEP tests against it

SCRIPT=$(realpath "$0")
TTDB_SCRIPTS=$(dirname "$SCRIPT")

# stop all running docker containers to avoid conflicts 
docker stop $(docker ps -aq)

# cleanup before (if docker-compose was run in other ways than this script there might be some leftovers)
sudo rm -rf $TTDB_SCRIPTS/../docker-test/neo4j/backups/*
sudo rm -rf $TTDB_SCRIPTS/../docker-test/neo4j/data/*
sudo rm -rf $TTDB_SCRIPTS/../docker-test/timescaledb/backups/*
sudo rm -rf $TTDB_SCRIPTS/../docker-test/timescaledb/data/*


# prepare testing data 
cp $TTDB_SCRIPTS/../test-data/neo4j_test_backup/neo4j.dump $TTDB_SCRIPTS/../docker-test/neo4j/backups/
cp $TTDB_SCRIPTS/../test-data/timescaledb_test_backup/postgres.bak $TTDB_SCRIPTS/../docker-test/timescaledb/backups/

# restore neo4j data (was not able to get this to work with docker-compose yet)
docker run --interactive --tty --rm --volume=$TTDB_SCRIPTS/../docker-test/neo4j/data:/data --volume=$TTDB_SCRIPTS/../docker-test/neo4j/backups:/backups neo4j neo4j-admin database load neo4j --from-path=/backups --verbose

# prepare testing envionment 
docker-compose -f $TTDB_SCRIPTS/../docker-compose.yml up -d

# wait for testing envionment to be ready 
dockerize -wait tcp://127.0.0.1:7687 -timeout 10s
dockerize -wait tcp://127.0.0.1:5432 -timeout 10s

# check if neo4j database is ready
docker exec test_neo4j sh -c "while neo4j-admin database info | grep -q 'Database in use:.*false'; do echo 'neo4j database not ready yet'; sleep 2; done"
echo "neo4j database ready"

docker exec test_timescaledb sh -c "while ! psql -lqt | cut -d \| -f 1 | grep -qw postgres; do echo 'timescaledb database not ready yet'; sleep 2; done"
echo "timescaledb database ready"

# Run Golang tests
go test -v -count=1 $TTDB_SCRIPTS/../api -run TestShallowQueries -test.v #> >(tee -a /dev/stdout) 2> >(tee -a /dev/stderr >&2)
gotest1=$?
go test -v -count=1 $TTDB_SCRIPTS/../api -run DeepShallowQueries -test.v #> >(tee -a /dev/stdout) 2> >(tee -a /dev/stderr >&2)
gotest2=$?


# # stop docker compose - removes the containers
# docker-compose -f $TTDB_SCRIPTS/../docker-compose.yml down
# 
# # cleanup  
# sudo rm -rf $TTDB_SCRIPTS/../docker-test/neo4j/backups/*
# sudo rm -rf $TTDB_SCRIPTS/../docker-test/neo4j/data/*
# sudo rm -rf $TTDB_SCRIPTS/../docker-test/timescaledb/backups/*
# sudo rm -rf $TTDB_SCRIPTS/../docker-test/timescaledb/data/*

# if ( $gotest1 -eq 0 && $gotest2 -eq 0 ); then
#     echo "Tests passed"
#     exit 0
# else
#     echo "Tests failed"
#     exit 1
# fi