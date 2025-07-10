#!/bin/bash

docker compose \
    -f docker-compose.yml \
    -f dockerfiles/docker-compose.es.yml \
    -f dockerfiles/docker-compose.etcd.yml \
    -f dockerfiles/docker-compose.mongo.yml \
    -f dockerfiles/docker-compose.mysql.yml \
    -f dockerfiles/docker-compose.pulsar-cluster.yml \
    -f dockerfiles/docker-compose.redis.yml \
    up --build
