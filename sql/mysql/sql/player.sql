-- MySQL 初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS mydb;
SHOW DATABASES;

-- 选择数据库
USE mydb;

-- 创建 players 表
CREATE TABLE IF NOT EXISTS players (
    `id` INT AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `age` INT COMMENT '年龄',
    `level` INT,
    PRIMARY KEY (id)
);

-- 查看所有表
-- USE mydb;
SHOW TABLES;

-- 查看 users 表结构
DESC users;

INSERT INTO users (`username`, `email`) VALUES ('admin', '');

SELECT * FROM users;



-- 添加一个字段 password 到 users 表  NOT NULL：
ALTER TABLE users ADD COLUMN `password` VARCHAR(64);
-- 2. 更新现有数据，为所有现有行设置默认密码
UPDATE users SET `password` = 'default_password' WHERE `password` IS NULL;
-- 3. 添加 NOT NULL 约束
ALTER TABLE users MODIFY COLUMN `password` VARCHAR(64) NOT NULL;


INSERT INTO users (`username`, `email`, `password`) VALUES ('admin', '', 'admin');
INSERT INTO users (`username`, `email`, `password`) VALUES ('user1', 'user1@gmail.com', 'user1passwd');
INSERT INTO users (`username`, `email`, `password`) VALUES ('user2', 'user2@gmail.com', 'user2passwd');
SELECT * FROM users;

INSERT INTO users (`username`, `email`, `password`) VALUES ('user3', 'user3@gmail.com', 'user3passwd');
INSERT INTO users (`username`, `email`, `password`) VALUES ('user4', 'user4@gmail.com', 'user4passwd');
INSERT INTO users (`username`, `email`, `password`) VALUES ('user5', 'user5@gmail.com', 'user5passwd');


DELETE FROM users WHERE `username` = 'user4';