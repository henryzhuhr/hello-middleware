package model

import (
	"gorm.io/gorm"
)

// ProductSKU SKU 表
// SKU（Stock Keeping Unit，库存量单位） 是电商和零售行业中最核心的概念之一，用于唯一标识一个可销售的商品具体规格组合。
type ProductSKU struct {
	gorm.Model

	ProductID uint    `gorm:"column:product_id;not null" json:"productId"`
	Specs     JSONMap `gorm:"column:specs;type:json" json:"specs"` // map[string]string
	Price     float64 `gorm:"column:price;not null;type:decimal(10,2)" json:"price"`
	Stock     uint    `gorm:"column:stock;not null;default:0" json:"stock"`
	SKUCode   string  `gorm:"column:sku_code;size:100;uniqueIndex" json:"skuCode"`
}

func (ProductSKU) TableName() string {
	return "redis_product_cache_t_product_skus"
}

// ProductSKURepository 定义商品 SKU 的数据访问接口
type ProductSKURepository interface {
	// 创建一个新的商品 SKU
	Create(sku *ProductSKU) error

	// 根据 ID 查找商品 SKU
	FindByID(id uint) (*ProductSKU, error)

	// 更新商品 SKU
	Update(sku *ProductSKU) error

	// 删除商品 SKU
	Delete(id uint) error
}

// GORMProductSKURepository 实现 ProductSKURepository 接口
type GORMProductSKURepository struct {
	db *gorm.DB // 可以传入一个 gorm.DB 实例，或者传入 gorm.Tx 以支持事务
}

// NewProductSKURepository 创建一个新的 Repository 实例
// 可以传入一个 gorm.DB 实例，或者传入 db.Begin() 返回的 gorm.Tx 以支持事务，这叫 透明事务支持（Transparent Transaction Support）
func NewProductSKURepository(db *gorm.DB) ProductSKURepository {
	return &GORMProductSKURepository{db: db}
}
func (r *GORMProductSKURepository) Create(sku *ProductSKU) error {
	if err := r.db.Create(sku).Error; err != nil {
		return err
	}
	return nil
}

func (r *GORMProductSKURepository) FindByID(id uint) (*ProductSKU, error) {
	var sku ProductSKU
	if err := r.db.First(&sku, id).Error; err != nil {
		return nil, err
	}
	return &sku, nil
}

func (r *GORMProductSKURepository) Update(sku *ProductSKU) error {
	if err := r.db.Save(sku).Error; err != nil {
		return err
	}
	return nil
}

func (r *GORMProductSKURepository) Delete(id uint) error {
	if err := r.db.Delete(&ProductSKU{}, id).Error; err != nil {
		return err
	}
	return nil
}
