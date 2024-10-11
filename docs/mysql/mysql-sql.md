---
outline: deep
---

# MySQL SQL


SQL Syntax(语法) 可以分为以下几种类型：

- **DDL（数据定义语言）**：`CREATE`, `DROP`, `ALTER`, `TRUNCATE`，对表结构进行增删改
- **DML（数据操作语言）**：`INSERT`, `UPDATE`, `DELETE`, `CALL`, 对表当中的数据进行增删改
- **DQL（数据查询语言）**：查询语句 `SELECT(WHERER)`
- **TCL（事务控制语言）**：`COMMIT` 提交事务，`ROLLBACK` 回滚事务（TCL中的T是 Transaction）
- **DCL（数据控制语言）**：`GRANT` 授权、`REVOKE` 撤销权限等


## 数据定义

### 创建数据库或数据表

创建数据库使用 `CREATE DATABASE` 语句，

```sql
CREATE DATABASE database_name;
USE database_name;
```

创建数据表使用 `CREATE TABLE` 语句，并定义字段名和数据类型。

```sql
CREATE TABLE `table_name` (
    `column1` datatype,
    `column2` datatype,
    ...
);
```

创建表时，可以指定列的名称、数据类型和约束。例如：

```sql
CREATE TABLE users (
    id INT,
    name VARCHAR(255),
    email VARCHAR(255),
    age INT
);
```

## 修改数据表 ALTER

`ALTER TABLE` 语句用于添加、删除或修改表的列，分别使用 `ADD`、`DROP`、`MODIFY` 关键字。


> 用 `tn` 表示表名，`cn` 表示列名，`dt` 表示数据类型，`dv` 表示默认值。
添加字段
```sql
ALTER TABLE `tn` ADD `cn` datatype;
```

删除字段
```sql
ALTER TABLE `tn` DROP COLUMN `cn`;
```

修改字段数据类型
```sql
ALTER TABLE `tn` MODIFY COLUMN `cn` datatype;
```

修改字段类型和约束
```sql
ALTER TABLE `tn` MODIFY COLUMN `cn` datatype DEFAULT 'default_value';
ALTER TABLE `tn` MODIFY COLUMN `cn` datatype NOT NULL;
```
> - 约束可以是非空 `NOT NULL`、默认值 `DEFAULT`、唯一 `UNIQUE`、主键约束 `PRIMARY KEY`、外键约束 `FOREIGN KEY` 等。
> - 主键约束 `PRIMARY KEY`：唯一标识表中的每一行，主键必顺是唯一的，且不允许有 NULL 值。
> - 外键约束 `FOREIGN KEY`：用于关联两个表的关系，保证数据的一致性，外键约束的列必须是另一个表的主键。

重命名字段
```sql
ALTER TABLE `tn` CHANGE `old_cn` `new_cn` datatype;
-- (MySQL 8.0.13+)
ALTER TABLE `tn` RENAME `old_cn` to `new_cn`;
```

## 数据的增删改查

数据的增删改查操作分别对应 `INSERT`、`DELETE`、`UPDATE` 和 `SELECT` 语句。

### INSERT 语句

插入数据可以使用 `INSERT INTO` 语句，指定表名和字段名，然后指定要插入的值。
```sql
INSERT INTO table_name (column1, column2, ...)
VALUES (value1, value2, ...);
```

如果要插入所有字段的值，可以省略字段名。
```sql
INSERT INTO table_name VALUES (value1, value2, ...);
```

插入多组数据时，可以使用多个 `VALUES` 子句。
```sql
INSERT INTO table_name (column1, column2, ...)
VALUES (value1, value2, ...),
       (value1, value2, ...),
       ...;
```

### SELECT 语句和 WHERE 子句

SELECT 语句用于查询数据，可以指定要查询的字段名，也可以使用 `*` 表示所有字段。`WHERE` 子句用于过滤数据，只返回满足条件的数据。
```sql
SELECT column1, column2, ... 
FROM table_name
WHERE condition
[LIMIT number]; -- 限制返回的行数
```

### UPDATE 语句

```sql
UPDATE table_name
SET column1 = value1, column2 = value2, ...
WHERE condition;
```

### DELETE 语句

一般要和 `WHERE` 子句一起使用，否则会删除所有数据。
```sql
DELETE FROM table_name
WHERE condition;
```

### WHERE 子句

WHERE 子句用于过滤数据，只返回满足条件的数据。可以使用 `AND`、`OR`、`NOT` 等逻辑运算符，也可以使用 `IN`、`BETWEEN`、`LIKE` 等操作符。

优先级上：`NOT` > `AND` > `OR`，可以使用括号来改变优先级。

有一张玩家表 `players`，包含字段 `id`、`name`、`age`、等级`level`，可以使用以下查询语句：

| id  | name | age | level |
| --- | ---- | --- | ----- |
| 1   | 张三 | 18  | 2     |
| 2   | 李四 | 20  | 3     |
| 3   | 王五 | 22  | 1     |
| 4   | 赵六 | 25  | 4     |
| 5   | 钱七 | 30  | 3     |


```sql
SELECT * FROM players WHERE name = '张三'; -- 查询名字为 张三 的用户
SELECT * FROM players WHERE name IN ('张三', '李四'); -- 名字为 张三 或 李四 的用户
SELECT * FROM players WHERE age = 18 OR age = 30; -- 年龄为18或30的用户
SELECT * FROM players WHERE NOT age = 18; -- 年龄不为18的用户
SELECT * FROM players WHERE age > 18 AND (name = '张三' OR name = '李四'); -- 年龄大于18且名字为 Tom 或 Jerry 的用户
```

查找一个范围内的值，例如年龄在18到30之间的用户
```sql
SELECT * FROM players WHERE age >= 18 AND age =< 30;
SELECT * FROM players WHERE age BETWEEN 18 AND 30; -- 闭区间
```

模糊查询，可以用 `LIKE` 和 `%`/`_` 通配符，`%` 表示任意字符，`_` 表示一个字符，例如查询名字以 `Tom` 开头的用户
```sql
SELECT * FROM players WHERE name LIKE '张%';  -- 名字以 张 开头
SELECT * FROM players WHERE name LIKE '%三%'; -- 名字包含 三
SELECT * FROM players WHERE name LIKE '_三';   -- 名字为 三
```

正则表达式查询，可以使用 `REGEXP` 或 `RLIKE`，

- `. ` 匹配任意单个字符
- `^` 匹配字符串的开始
- `$` 匹配字符串的结束
- `[abc]` 匹配 a、b 或 c
- `[a-z]` 匹配 a 到 z 之间的任意字符
- `[^abc]` 匹配除了 a、b 或 c 之外的任意字符
- `a|b` 匹配 a 或 b 
- `*` 匹配前面的子表达式零次或多次


例如查询名字为 Tom 或 Jerry 的用户
```sql
SELECT * FROM players WHERE name REGEXP 'Tom|Jerry';
SELECT * FROM players WHERE name RLIKE 'Tom|Jerry';
```


查找空值或非空值，可以使用 `IS NULL` 或 `IS NOT NULL`，
```sql
SELECT * FROM players WHERE name IS NULL; -- 名字为空的用户
SELECT * FROM players WHERE name IS NOT NULL; -- 名字不为空的用户
```

> `= NULL` 和 `!= NULL` 是无效的，因为 NULL 代表未知的值，不能直接比较。

### ORDER BY 子句

`ORDER BY` 表示按照某个字段进行排序，如果不指定，默认升序排序

例如按照年龄升序和降序排列
```sql
SELECT * FROM players ORDER BY age;      -- 默认升序
SELECT * FROM players ORDER BY age ASC;  -- 升序
SELECT * FROM players ORDER BY age DESC; -- 降序
```

如果希望多个列排序，用 `,` 分开
```sql
SELECT * FROM players ORDER BY age, name;      -- 先按年龄升序，再照名字升序
SELECT * FROM players ORDER BY age, name DESC; -- 先按年龄升序，再照名字降序
```

### 聚合函数

聚合函数用于对一组值进行计算，常用的聚合函数有 `COUNT`、`SUM`、`AVG`、`MAX`、`MIN` 等。

```sql
SELECT COUNT(*) FROM players; -- 计算用户的数量 COUNT(*) 是统计记录条数 
SELECT SUM(age) FROM players; -- 计算年龄的总和
SELECT AVG(age) FROM players; -- 计算年龄的平均值
```

> `COUNT(*)` 和 `COUNT(列名)` 的区别是，`COUNT(*)` 统计所有行数，包括 NULL 值，而 `COUNT(列名)` 统计非空值的行数。



### GROUP BY 分组查询

使用 `GROUP BY` 子句对查询结果进行分组，常用于聚合函数的计算。

例如想知道18岁以上的用户数量，可以使用 `GROUP BY` 子句
```sql
SELECT age, COUNT(*) FROM players WHERE age > 18 GROUP BY age;
```

统计每个年龄有多少人
```sql
SELECT age, COUNT(*) FROM players GROUP BY age;
```
输出的结果
```bash
+------+----------+
| age  | COUNT(*) |
+------+----------+
|  18  |    2     |
|  20  |    1     |
...
```

分组可以与 `HAVING` 子句一起使用，`HAVING` 子句用于过滤分组后的结果。
```sql
SELECT age, COUNT(*) FROM players GROUP BY age HAVING COUNT(*) > 1;
```
与 `WHERE` 子句的区别是，`WHERE` 子句在数据分组前进行过滤，`HAVING` 子句在数据分组后进行过滤，所以 `HAVING` 子句中的条件可以使用聚合函数。

### 集合运算查询 UNION、UNION ALL、INTERSECT、EXCEPT

## 数据的增删改查：子查询

有些情况下需要将一个查询的结果作为另一个查询的条件，这时可以使用子查询。

例如查询玩家等级大于平均等级的玩家
```sql
-- 用在 WHERE 子句中
SELECT * FROM players WHERE level > (SELECT AVG(level) FROM players);
-- 用在 SELECT 子句中 可以起别名
SELECT name, (SELECT AVG(level) FROM players) AS al FROM players;
-- ROUND 函数四舍五入
SELECT name, ROUND((SELECT AVG(level) FROM players)) AS al FROM players; 
```

查询结果可以创建一个新的表，然后再查询
```sql
CREATE TABLE new_table AS SELECT * FROM players WHERE level > (SELECT AVG(level) FROM players);
SELECT * FROM new_table;
```

## 表关联

表关联用于查询多个表的数据，并且两个表之间有相同的字段，通常是主键和外键。关联可以使用 `JOIN` 子句，`JOIN` 子句可以分为：
- 内连接 `INNER JOIN`：只返回两个表中都有的数据
- 外连接 `OUTER JOIN`：返回两个表中的所有数据
  - 左外连接 `LEFT JOIN`：返回左表中的所有数据
  - 右外连接 `RIGHT JOIN`：返回右表中的所有数据
  - 全外连接 `FULL JOIN`：返回两个表中的所有数据





## 参考和资料

- [MySQL语法总结-基础语法 - 辰智的文章 - 知乎](https://zhuanlan.zhihu.com/p/63111767)
- [🔥MySQL一万字深度总结，基础+进阶(一)，建议收藏。✨💖-腾讯云开发者社区-腾讯云](https://cloud.tencent.com/developer/article/1856648)               
- [🔥MySQL一万字深度总结，基础+进阶(二)，建议收藏。✨💖-腾讯云开发者社区-腾讯云](https://cloud.tencent.com/developer/article/1858575)               