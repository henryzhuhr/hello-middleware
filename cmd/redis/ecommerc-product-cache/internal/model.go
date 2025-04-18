package ecommercproductcache

import (
	"gorm.io/gorm"
)

// 商品数据模型
type Product struct {
	gorm.Model
	ID    string  `gorm:"primaryKey"`
	Name  string  `gorm:"size:255"`
	Price float64 `gorm:"type:decimal(10,2)"`
}

// 商品仓库接口
type ProductRepository interface {
	CreateProduct(product *Product) error
	UpdateProduct(product *Product) error
	DeleteProduct(id string) error
}

// 商品仓库实现
type ProductRepositoryImpl struct {
	DB *gorm.DB
}

func (Product) TableName() string {
	return "t_product"
}

func (r *ProductRepositoryImpl) CreateProduct(product *Product) error {
	return r.DB.Create(product).Error
}

func (r *ProductRepositoryImpl) UpdateProduct(product *Product) error {
	return r.DB.Save(product).Error
}

func (r *ProductRepositoryImpl) DeleteProduct(id string) error {
	return r.DB.Delete(&Product{}, "id = ?", id).Error
}
