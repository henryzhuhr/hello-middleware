package plugin

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlDB      *gorm.DB  // 包级变量，存储数据库连接实例
	mysqlDB_once sync.Once // 确保数据库连接只被创建一次
)

const (
	DEFAULT_MYSQL_USER     = "root"      // 默认 MySQL 用户名
	DEFAULT_MYSQL_HOST     = "localhost" // 默认 MySQL 主机
	DEFAULT_MYSQL_PORT     = "3306"      // 默认 MySQL 端口
	DEFAULT_MYSQL_DATABASE = "test"      // 默认 MySQL 数据库名称
)

func NewMysqlClient() (*gorm.DB, error) {
	var err error

	mysqlDB_once.Do(func() {

		// 从环境变量中读取数据库配置
		dbUser := os.Getenv("MYSQL_ROOT_USER")
		if dbUser == "" {
			dbUser = DEFAULT_MYSQL_USER
		}
		dbPassword := os.Getenv("MYSQL_ROOT_PASSWORD")

		dbHost := os.Getenv("MYSQL_HOST")
		if dbHost == "" {
			dbHost = DEFAULT_MYSQL_HOST
		}
		dbPort := os.Getenv("MYSQL_PORT")
		if dbPort == "" {
			dbPort = DEFAULT_MYSQL_PORT
		}
		dbName := os.Getenv("MYSQL_DATABASE")
		if dbName == "" {
			dbName = DEFAULT_MYSQL_DATABASE
		}

		// if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		// 	err = fmt.Errorf("database configuration is incomplete")
		// 	return
		// }

		// 构建 DSN (Data Source Name)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName)

		// 初始化 GORM 数据库连接
		mysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			err = fmt.Errorf("failed to connect to database: %v", err)
			return
		}
	})
	// 如果初始化过程中发生错误，返回错误
	if err != nil {
		return nil, err
	}
	return mysqlDB, nil
}
