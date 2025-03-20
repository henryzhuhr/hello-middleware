# Hello Middleware


## 使用

根据当前目录下的 `docker-compose.yml` 文件启动容器
```shell
docker compose up -d
docker compose up -d --build            # 重新构建镜像
docker compose up -d --force-recreate   # 强制重新创建容器
docker compose up -d --build --force-recreate
```

如果希望进入容器内部，可以使用以下命令
```shell
docker compose exec <service-name>  /bin/bash
docker compose exec hello-databse-dev /bin/bash # 进入开发环境容器
docker compose exec mysql_db    /bin/bash       # 进入 mysql    数据库容器
docker compose exec postgres_db /bin/bash       # 进入 postgres 数据库容器

```

停止并删除所有与当前 `docker-compose.yml` 文件关联的容器、网络和卷
```shell
docker compose down
```

组合上述开发命令
```shell
docker compose up -d --build --force-recreate && \
docker compose exec hello-databse-dev /bin/bash && \
docker compose down
```

## 进入开发环境容器


### 连接 MySQL 数据库
```shell
mysql -h mysql_db  -u root -p<PASSWORD> # 容器内连接 p后不加空格
mysql -h 127.0.0.1 -u root -p<PASSWORD> # 外部连接
```

### 连接 PostgreSQL 数据库
```shell
psql -h postgres_db -U postgres -d postgres
# -h, --host=HOSTNAME      database server host or socket directory (default: "local socket")
# -U, --username=USERNAME   database user name
# -d, --dbname=DBNAME       database name to connect to (default: same as username)
```


```shell
\l
\c mydatabase
\dt
select * from users;
\q
```
### 连接 Redis

```shell
redis-cli -h redis-server
```

## GORM

[中文文档](https://gorm.io/zh_CN/docs/index.html)