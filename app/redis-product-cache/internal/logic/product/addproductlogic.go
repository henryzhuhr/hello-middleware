package product

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/model"
	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/svc"
	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductLogic {
	return &AddProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 添加商品
func (l *AddProductLogic) AddProduct(req *types.AddProductReq) (resp *types.AddProductResp, err error) {
	// 1. 构建 Product
	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryId,
		Brand:       req.Brand,
		Price:       req.Price,
		Attributes:  model.JSONMap(req.Attributes),
	}
	// ✅ 正确：判断是否传了促销价（无论值是多少）
	// 如果 req.SalePrice == nil，则 product.SalePrice 保持默认（Valid=false）
	if req.SalePrice != nil {
		product.SalePrice = sql.NullFloat64{
			Float64: *req.SalePrice, // 解引用
			Valid:   true,
		}
	}

	// 2. 事务写入
	tx := l.svcCtx.DB.WithContext(l.ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 3. 创建商品
	if err := tx.Create(&product).Error; err != nil {
		return nil, fmt.Errorf("")
	}

	// 4. 创建 SKU
	for _, skuReq := range req.SKUs {
		sku := model.ProductSKU{
			ProductID: product.ID,
			Specs:     model.JSONMap(skuReq.Specs),
			Price:     skuReq.Price,
			Stock:     skuReq.Stock,
			SKUCode:   skuReq.SKUCode,
		}
		if err := tx.Create(&sku).Error; err != nil {
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 5. 【关键】删除缓存（最终一致性）
	cacheKey := productCachePrefix + strconv.FormatUint(uint64(product.ID), 10)
	emptyKey := cacheKey + emptyCacheSuffix

	// 删除主缓存 + 空值缓存
	_, _ = l.svcCtx.Redis.Del(l.ctx, cacheKey, emptyKey).Result()

	return &types.AddProductResp{Id: product.ID}, nil
}
