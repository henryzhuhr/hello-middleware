package svc

import (
	"context"
	"fmt"

	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/config"
	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	DB    *gorm.DB
	Redis *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	fmt.Printf("Initializing service context with config: %+v\n", c)
	svcCtx := &ServiceContext{
		Config: c,
	}

	// 初始化 MySQL 的连接
	svcCtx.DB = initMySQL(c)

	// 初始化 redis 的连接
	svcCtx.Redis = initRedis(c)

	return svcCtx
}
func initMySQL(c config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
		PrepareStmt:            true, // 缓存预编译语句
		SkipDefaultTransaction: true, // 禁用默认事务（提升性能）
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MySQL: %v", err))
	}
	// 测试连接
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get sql.DB from gorm.DB: %v", err))
	}
	err = sqlDB.Ping()
	if err != nil {
		panic(fmt.Sprintf("Failed to ping MySQL: %v", err))
	}
	fmt.Println("Successfully connected to MySQL")

	// 自动迁移表结构（仅创建不存在的表，不会删除字段）（开发环境可用，生产慎用）
	err = db.AutoMigrate(
		&model.Product{},
		&model.ProductSKU{},
		// 如果你有 categories 表，也可加上 &model.Category{}
	)
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}
	fmt.Println("Database migration completed")
	return db
}

func initRedis(c config.Config) *redis.Client {
	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password, // no password set
		DB:       c.Redis.DB,       // use default DB
	}
	rdb := redis.NewClient(opt)
	// 测试连接
	res, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	fmt.Printf("Successfully connected to Redis (Ping response: %s)\n", res)

	return rdb
}
