---
outline: deep
---

# MySQL 索引 Index

## 为什么需要索引

在数据库中，索引是一种特殊的查找表，帮助数据库系统高效地获取数据。索引是一种数据结构，可以帮助数据库系统快速地找到表中的数据。索引可以大大提高查询效率，但是索引也会降低插入、删除和更新表的速度。

如果没有索引的表，在 `SELECT` 查询的时候，那么数据会全表扫描（从头查询到结尾），查询效率会非常低。索引可以帮助数据库系统快速地找到表中的数据，提高查询效率。


如果没有索引的表，在 `SELECT` 查询的时候，那么数据会全表扫描（从头查询到结尾），查询效率会非常低。

建立索引，会使查询效率提高，但是会降低插入、删除和更新表的速度。


## 创建和删除索引

创建索引可以使用 `CREATE INDEX` 关键字：
```sql
CREATE [UNIQUE|FULLTEXT|SPATIAL] INDEX index_name
ON table_name (column_name_1, ...);
```

上述的语句中，表示在 `table_name` 表上创建索引 `index_name`，索引的列是 `column_name_1`、`...` 等列。

其中 `UNIQUE` 表示唯一索引，`FULLTEXT` 表示全文索引，`SPATIAL` 表示空间索引。

一般来说，会对经常查询的列创建索引，例如主键、外键、经常用于查询的列等。

删除索引可以使用 `DROP INDEX` 关键字：
```sql
DROP INDEX index_name ON table_name;
```

创建索引的时候也可以在使用 `CREATE TABLE` 语句的时候创建索引：
```sql
CREATE TABLE table_name (
    column_name_1 datatype,
    column_name_2 datatype,
    ...
    INDEX index_name (column_name_1, column_name_2)
);
```
也可以在 `ALTER TABLE` 语句中添加索引：
```sql
ALTER TABLE table_name
ADD INDEX index_name (column_name_1, column_name_2);
```