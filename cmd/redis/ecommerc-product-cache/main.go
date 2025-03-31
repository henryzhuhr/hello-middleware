// 电商平台商品缓存
package main

import (
	"github.com/henryzhuhr/hello-middleware/cmd/redis/ecommerc-product-cache/internal/model"
	"github.com/henryzhuhr/hello-middleware/internal/plugin"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	rdb     *redis.Client // Redis 客户端实例
	mysqlDB *gorm.DB      // MySQL 客户端实例
)

func main() {
	// Initialize the Redis client
	rdb = plugin.NewRedisClient()

	if rdb == nil {
		log.Fatalf("failed to create redis client")
	}

	mysqlDB, err := plugin.NewMysqlClient()
	if err != nil {
		log.Fatalf("failed to create mysql client, error is %s", err)
	}
	if mysqlDB == nil {
		log.Fatalf("failed to create mysql client")
	}
	// 自动迁移（创建表）
	err = mysqlDB.AutoMigrate(&model.Product{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	productRepository := model.NewProductRepositoryWithCache(mysqlDB, rdb)

	productID, err := productRepository.CreateProduct(&model.Product{
		Name:  "Product 1",
		Price: 100.0,
	})
	if err != nil {
		log.Error("Error creating product:", err)
	} else {
		log.Infof("Product created with ID: %s", productID)
	}

	// 查找
	product, err := productRepository.FindProductByID(productID)
	if err != nil {
		log.Error("Error finding product:", err)
	} else {
		log.Infof("Product found: %+v", product)
	}

	productID, err = productRepository.CreateProductWithCache(&model.Product{
		Name:  "Product 1",
		Price: 100.0,
	})
	if err != nil {
		log.Error("Error creating product:", err)
	} else {
		log.Infof("Product created with ID: %s", productID)
	}

	// 查找
	product, err = productRepository.FindProductByID(productID)
	if err != nil {
		log.Error("Error finding product:", err)
	} else {
		log.Infof("Product found: %+v", product)
	}
}
