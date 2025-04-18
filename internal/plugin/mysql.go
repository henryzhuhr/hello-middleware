package plugin

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysql_db *gorm.DB // 包级变量，存储数据库连接实例

func NewMysqlClient() *gorm.DB {
	once.Do(func() {
		// 这里使用 gorm.Open() 来连接 MySQL 数据库
		var err error
		mysql_db, err = gorm.Open(mysql.Open("user:password@tcp(localhost:3306)/dbname"), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
	})
	return mysql_db
}
