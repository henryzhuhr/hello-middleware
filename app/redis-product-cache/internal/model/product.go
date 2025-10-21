package model

import (
	"database/sql"

	"gorm.io/gorm"
)

// Product 商品主表（SPU）
type Product struct {
	gorm.Model

	Name        string          `gorm:"column:name;not null;size:255" json:"name"`
	Description string          `gorm:"column:description;type:text" json:"description"`
	CategoryID  uint            `gorm:"column:category_id;not null" json:"categoryId"`
	Brand       string          `gorm:"column:brand;size:100" json:"brand"`
	Price       float64         `gorm:"column:price;not null;type:decimal(10,2)" json:"price"`
	SalePrice   sql.NullFloat64 `gorm:"column:sale_price;default:null;type:decimal(10,2)" json:"salePrice"`
	Attributes  JSONMap         `gorm:"column:attributes;type:json" json:"attributes"` // map[string]string
}

func (Product) TableName() string {
	return "redis_product_cache_t_product"
}

// ProductRepository 定义商品的数据访问接口
type ProductRepository interface {
	// 创建一个新的商品
	Create(product *Product) error

	// 根据 ID 查找商品
	FindByID(id uint) (*Product, error)

	// 更新商品信息
	Update(product *Product) error

	// 删除商品
	Delete(id uint) error
}

// GORMProductRepository 实现 ProductRepository 接口
type GORMProductRepository struct {
	db *gorm.DB // 可以传入一个 gorm.DB 实例，或者传入 gorm.Tx 以支持事务
}

// NewProductRepository 创建一个新的 Repository 实例
// 可以传入一个 gorm.DB 实例，或者传入 db.Begin() 返回的 gorm.Tx 以支持事务，这叫 透明事务支持（Transparent Transaction Support）
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &GORMProductRepository{db: db}
}

func (r *GORMProductRepository) Create(product *Product) error {
	if err := r.db.Create(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *GORMProductRepository) FindByID(id uint) (*Product, error) {
	var product Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *GORMProductRepository) Update(product *Product) error {
	if err := r.db.Save(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *GORMProductRepository) Delete(id uint) error {
	if err := r.db.Delete(&Product{}, id).Error; err != nil {
		return err
	}
	return nil
}
