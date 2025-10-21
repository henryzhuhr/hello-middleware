package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	Mysql struct {
		DataSource string `json:"DataSource"` // MySQL 数据源
	} `json:"Mysql"`
	Redis Redis `json:"Redis"`
}

type Redis struct {
	Host     string `json:"Host"`     // Redis 服务器地址
	Port     int    `json:"Port"`     // Redis 服务器端口
	Password string `json:"Password"` // Redis 服务器密码
	DB       int    `json:"DB"`       // 使用的数据库编号
}
