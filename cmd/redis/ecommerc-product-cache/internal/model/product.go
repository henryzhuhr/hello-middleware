// 可以参考 https://github.com/8treenet/gcache
package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 商品数据模型
type (
	Product struct {
		gorm.Model
		ID    string  `gorm:"primaryKey"`         // 商品ID, UUID格式，覆盖gorm.Model的ID
		Name  string  `gorm:"size:255"`           // 商品名称
		Price float64 `gorm:"type:decimal(10,2)"` // 商品价格
	}

	// 商品仓库接口
	ProductRepository interface {
		// 创建商品
		CreateProduct(product *Product) (string, error)
		CreateProductWithCache(product *Product) (string, error)

		UpdateProduct(product *Product) error
		DeleteProduct(id string) error

		// 查找商品
		FindProductByID(id string) (*Product, error)
	}

	// 商品仓库实现
	ProductRepositoryImpl struct {
		DB  *gorm.DB
		rdb *redis.Client
	}
)

func (Product) TableName() string {
	return "t_product"
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{DB: db}
}

func NewProductRepositoryWithCache(db *gorm.DB, rdb *redis.Client) ProductRepository {
	return &ProductRepositoryImpl{DB: db, rdb: rdb}
}

func (i *ProductRepositoryImpl) CreateProduct(product *Product) (string, error) {
	// 清理掉 ID 字段，使用 UUID 作为主键
	product.ID = uuid.New().String() // 清空 ID 字段，使用 UUID 作为主键
	err := i.DB.Create(product).Error
	if err != nil {
		return "", err
	}
	return product.ID, nil
}

// CreateProductWithCache implements ProductRepository.
func (i *ProductRepositoryImpl) CreateProductWithCache(product *Product) (string, error) {

	productID, err := i.CreateProduct(product)
	if err != nil {
		return "", err
	}
	if i.rdb == nil {
		log.Warnf("Redis client is not initialized, skipping caching for product with ID %s", product.ID)
		return productID, nil
	}

	// 将商品信息序列化为 JSON
	productJson, err := json.Marshal(product)
	if err != nil {
		return "", err
	}
	// 将商品信息存储到 Redis 中，设置过期时间为 1 小时
	err = i.rdb.Set(context.Background(), product.ID, productJson, 1*time.Hour).Err()
	if err != nil {
		return "", err
	}
	log.Infof("Product with ID %s cached in Redis", product.ID)
	return productID, nil
}

func (i *ProductRepositoryImpl) UpdateProduct(product *Product) error {
	return i.DB.Save(product).Error
}

func (i *ProductRepositoryImpl) DeleteProduct(id string) error {
	return i.DB.Delete(&Product{}, "id = ?", id).Error
}

// FindProductByID implements ProductRepository.
func (i *ProductRepositoryImpl) FindProductByID(id string) (*Product, error) {
	if i.rdb != nil {
		// 先尝试从 Redis 中获取商品信息
		val, err := i.rdb.Get(context.Background(), id).Result()

		if err == redis.Nil {
			log.Warnf("Product with ID %s not found in Redis", id)
		} else if err != nil {
			log.Errorf("Error getting product from Redis: %v", err)
		} else {
			var product Product
			// 反序列化 Redis 中的商品信息
			jsonErr := json.Unmarshal([]byte(val), &product)
			if jsonErr != nil {
				log.Errorf("Error unmarshalling product from Redis: %v", jsonErr)
			} else {
				return &product, nil
			}
		}
	}

	// 如果 Redis 中不存在商品信息，则从 MySQL 中查询
	product := &Product{}
	err := i.DB.First(product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}
