#!/bin/bash

# This file sets up an testing environment providing a docker Neo4j DMBS and a docker TimescaleDB DMBS which run databases with 
# simplistic testing data 

SCRIPT=$(realpath "$0")
TTDB_SCRIPTS=$(dirname "$SCRIPT")

# stop all running docker containers to avoid conflicts 
docker stop $(docker ps -aq)

# rm testing containers 
docker rm test_neo4j
docker rm test_timescaledb

# cleanup before (if docker-compose was run in other ways than this script there might be some leftovers)
sudo rm -rf $TTDB_SCRIPTS/../docker-test/neo4j/backups/*
sudo rm -rf $TTDB_SCRIPTS/../docker-test/neo4j/data/*
sudo rm -rf $TTDB_SCRIPTS/../docker-test/timescaledb/backups/*
sudo rm -rf $TTDB_SCRIPTS/../docker-test/timescaledb/data/*

docker rm test_timescaledb
docker rm test_neo4j

# prepare testing envionment 
docker-compose -f $TTDB_SCRIPTS/../docker-compose.yml up -d

# wait for testing envionment to be ready 
dockerize -wait tcp://127.0.0.1:7687 -timeout 10s
dockerize -wait tcp://127.0.0.1:5432 -timeout 10s