package main

import (
	"flag"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:100"`
	Email    string `gorm:"size:100"`
}

var postgres_connect_url string

func testdb() {
	flag.StringVar(&postgres_connect_url, "postgres_connect_url", "", "Postgres connection URL")
	flag.Parse()

	// 连接数据库
	var db *gorm.DB
	db, err := gorm.Open(postgres.Open(postgres_connect_url), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Connecting to Postgres at", postgres_connect_url)

	// 迁移 schema
	db.AutoMigrate(&User{})

	// 获取数据表中的全部数据
	var users []User
	db.Find(&users)
	log.Println("All users:", users)

	// 插入数据
	db.Create(&User{Username: "jinzhu", Email: "111@gmail.com"})

	// 获取数据表中的全部数据
	db.Find(&users)
	log.Println("All users:", users)

}
