# mysql: https://hub.docker.com/_/mysql
# docker pull mysql
docker pull mysql:9

# redis: https://hub.docker.com/_/redis
# docker pull redis
docker pull redis:7

# mongo: https://hub.docker.com/_/mongo
# docker pull mongo
docker pull mongo:8.0

# kafka: https://hub.docker.com/r/bitnami/kafka
# docker pull bitnami/kafka
docker pull bitnami/kafka:3.9

# etcd: https://hub.docker.com/r/bitnami/etcd
# docker pull bitnami/etcd
docker pull bitnami/etcd:3.5

# # elasticsearch: https://hub.docker.com/_/elasticsearch
# docker pull elasticsearch
# docker pull elasticsearch:8.17.1

# elasticsearch: https://hub.docker.com/r/bitnami/elasticsearch
# docker pull bitnami/elasticsearch
docker pull bitnami/elasticsearch:8.17.1

# clickhouse: https://hub.docker.com/r/bitnami/clickhouse
# docker pull bitnami/clickhouse
docker pull bitnami/clickhouse:25

# milvus: https://hub.docker.com/r/bitnami/milvus
# docker pull --platform linux/amd64 bitnami/milvus
docker pull --platform linux/amd64 milvusdb/milvus:v2.5.4
docker pull quay.io/coreos/etcd:v3.5.16
docker pull minio/minio:RELEASE.2023-03-20T20-16-18Z
docker pull zilliz/attu:v2.4