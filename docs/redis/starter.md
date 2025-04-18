

# Redis

## Redis 简介

Redis 是 Remote Dictionary Server（远程字典服务器）的缩写，是一个开源的使用 ANSI C 语言编写、支持网络、可**基于内存**亦可持久化的日志型、Key-Value 数据库，并提供多种语言的 API。


## Redis 安装

在 Ubuntu 上安装 Redis 可以使用以下命令：

```bash
sudo apt update
sudo apt install redis-server
```


## Redis 配置
## Redis 基本命令

Redis 中的变量命令主要有以下几个：

- `SET`：设置指定 key 的值
- `GET`：获取指定 key 的值
- `DEL`：删除指定 key
- `FLUSHALL`：删除所有 key
- `EXISTS`：判断 key 是否存在
- `EXPIRE`：设置 key 的过期时间
- `TTL`：获取 key 的过期时间
- `PERSIST`：移除 key 的过期时间
- `KEYS`：获取所有 key，可以使用通配符
- `RENAME`：重命名 key
- `TYPE`：获取 key 的类型
- `APPEND`：追加字符串
- `STRLEN`：获取字符串长度
- `INCR`：自增

### 基本操作

Redis 中使用 key-value 存储数据，`SET` 命令用于设置指定 key 的值。
```bash
SET key value
# 完整语法
SET key value [NX|XX] [GET] [EX seconds|PX milliseconds|EXAT unix-time-seconds|PXAT unix-time-milliseconds|KEEPTTL]
```

使用 `GET` 命令获取 key 的值：
```bash
GET key
```

Redis 的 key 是区分大小写的，所以 `key` 和 `Key` 是两个不同的 key。

Redis 默认使用字符串存储键和值，而且是二进制安全的，所以 `value` 可以是字符串、数字、布尔值、序列化对象等。

Redis

### 过期时间

Redis 中可以为 key 设置过期时间，使用 `EXPIRE` 命令设置 key 的过期时间：
```bash
EXPIRE key seconds
```

使用 `TTL` 命令获取 key 的过期时间：
```bash
TTL key
```

使用 `PERSIST` 命令移除 key 的过期时间：
```bash
PERSIST key
```

也可以在设置的时候直接设置过期时间：
```bash
SET key value EX seconds
SETEX key seconds value
```



## 十大数据类型
## 事务
## 数据持久化
## 主从复制
## 哨兵模式 