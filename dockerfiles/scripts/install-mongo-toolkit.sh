#!/bin/bash

# install guide: https://www.mongodb.com/zh-cn/docs/manual/tutorial/install-mongodb-on-ubuntu/

MONOGODB_SERVER_VERSION=8.0

# 1. Import the public key
curl -fsSL https://www.mongodb.org/static/pgp/server-$MONOGODB_SERVER_VERSION.asc | \
    gpg --yes -o /usr/share/keyrings/mongodb-server-$MONOGODB_SERVER_VERSION.gpg \
    --dearmor

# 2. Create the list file /etc/apt/sources.list.d/mongodb-org-<version>.list for your version of Ubuntu.
echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-${MONOGODB_SERVER_VERSION}.gpg ] \
  https://repo.mongodb.org/apt/ubuntu  \
  $(lsb_release -cs)/mongodb-org/${MONOGODB_SERVER_VERSION} multiverse" | \
tee /etc/apt/sources.list.d/mongodb-org-${MONOGODB_SERVER_VERSION}.list

apt update

# 安装 MongoDB 客户端
apt install -y mongodb-mongosh