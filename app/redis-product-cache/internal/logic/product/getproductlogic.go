package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/model"
	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/svc"
	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/types"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 商品详情查询服务的 缓存高可用策略，包含了：
//
// - 缓存穿透防护（空值缓存）
// - 缓存击穿防护（互斥锁 + 双检）
// - 缓存雪崩防护（随机过期时间）
// - Redis 异常降级处理（直接查数据库）
//
// 整体流程如下：
//  1. 检查 emptyKey → 是否已标记为“不存在”
//     ├─ 是 → 返回 "not found"
//     └─ 否 → 继续
//  2. 尝试从 Redis 获取 cacheKey
//     ├─ 命中且反序列化成功 → 直接返回
//     ├─ 缓存损坏或未命中 → 进入加锁查 DB 流程
//  3. 加互斥锁（SetNX）尝试获取写权限
//     ├─ 成功获取 → 查 DB → 回填缓存 → 返回结果
//     └─ 失败（拿不到锁）→ 等待 20ms 再试一次缓存 → 若仍失败 → 降级直连 DB
//  4. 降级路径：直接查 DB（不加锁），同时设置 empty 缓存防穿透
func (l *GetProductLogic) GetProduct(req *types.GetProductReq) (resp *types.GetProductResp, err error) {

	// 设置缓存的key
	cacheKey := productCachePrefix + strconv.FormatUint(uint64(req.Id), 10)
	emptyKey := cacheKey + emptyCacheSuffix

	// ===== 1. 防缓存穿透：检查空值缓存 =====
	// 使用独立的 emptyKey 标记“该 ID 商品不存在”，避免频繁访问数据库。
	// ✅ 有效防止恶意请求或无效 ID 导致数据库压力过大。
	val, err := l.svcCtx.Redis.Get(l.ctx, emptyKey).Result()
	if err == nil && val == "1" {
		return nil, errors.New("product not found")
	}
	// 注意：redis.Nil 表示无空标记，继续；其他错误（如网络）则忽略，往下走

	// ===== 2. 查主缓存 =====
	val, err = l.svcCtx.Redis.Get(l.ctx, cacheKey).Result()
	if err == nil {
		// 命中缓存
		var cachedResp types.GetProductResp
		if err := json.Unmarshal([]byte(val), &cachedResp); err == nil {
			return &cachedResp, nil
		}
		// 缓存损坏，继续查数据库
	} else if err != redis.Nil {
		// Redis 其他错误（如连接失败），可选择降级直连数据库
		l.Logger.Errorf("redis get error: %v", err)
	}

	// ===== 3. 查数据库（加互斥锁防缓存击穿）=====
	//  防止缓存击穿（热点 key 失效瞬间大量请求打到 DB）
	// 只有一个协程能拿到锁去查 DB 并回填缓存。
	// 其他竞争者短暂等待后重试读缓存（双检机制）。
	// 锁具有 TTL（mutexExpire），防止死锁。
	// 高并发下保护数据库，避免击穿。
	resp, err = l.getProductFromDBWithMutex(req.Id, cacheKey, emptyKey)
	if err != nil {
		return nil, fmt.Errorf("getProductFromDBWithMutex error: %w", err)
	}

	return resp, nil
}

func (l *GetProductLogic) getProductFromDBWithMutex(id uint, cacheKey, emptyKey string) (*types.GetProductResp, error) {
	mutexKey := cacheKey + ":mutex"
	acquired := false

	// 生成唯一锁标识（防止误删）
	mutexVal := generateLockToken()

	// 尝试获取锁（最多等待 100ms）
	const count = 10
	for i := 0; i < count; i++ {
		ok, err := l.svcCtx.Redis.SetNX(l.ctx, mutexKey, mutexVal, mutexExpire).Result()
		if err != nil {
			break // Redis 错误，跳过加锁，直接查数据库
		}
		if ok {
			acquired = true
			break // 成功获取锁
		}
		// 等待 10ms 后重试
		timer := time.NewTimer(10 * time.Millisecond)
		select {
		case <-l.ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, l.ctx.Err()
		case <-timer.C:
		}
	}

	if !acquired {
		// 未获取到锁：短暂等待后重试读缓存（双检）（避免所有请求都打到数据库）
		select {
		case <-l.ctx.Done():
			return nil, l.ctx.Err()
		case <-time.After(20 * time.Millisecond):
		}

		// 再次尝试从缓存读取
		if val, err := l.svcCtx.Redis.Get(l.ctx, cacheKey).Result(); err == nil {
			var resp types.GetProductResp
			if err := json.Unmarshal([]byte(val), &resp); err == nil {
				return &resp, nil
			}
		}
		// 仍失败，直连数据库（降级）
		return l.getProductFromDB(id, emptyKey)
	}

	// 成功拿到锁：查数据库 + 回填缓存
	defer func() {
		// 释放锁（简单删除，不考虑原子性，因有 TTL）
		if _, err := l.svcCtx.Redis.Del(l.ctx, mutexKey).Result(); err != nil {
			l.Logger.Errorf("redis del mutex error: %v", err)
		}
	}()

	// 查询商品主表
	productRepo := model.NewProductRepository(l.svcCtx.DB)
	product, err := productRepo.FindByID(id)
	if err != nil {
		// 查无此商品，写入空值缓存防穿透
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.svcCtx.Redis.SetEx(l.ctx, emptyKey, "1", emptyExpire)
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("get product from db error: %w", err)
	}

	resp := l.buildProductResp(product)

	// 回填缓存：序列化并写入缓存
	data, err := json.Marshal(resp)
	if err != nil {
		return resp, nil // 序列化失败，不影响正常返回
	}
	// 写入缓存（带随机过期时间防雪崩）
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	expire := baseExpire + time.Duration(rng.Intn(300))*time.Second // +0~5分钟
	if _, err := l.svcCtx.Redis.SetEx(l.ctx, cacheKey, data, expire).Result(); err != nil {
		l.Logger.Errorf("redis setex error: %v", err)
	}
	return resp, nil

}

// generateLockToken 生成唯一 token（简单版可用 timestamp+random）
func generateLockToken() string {
	return time.Now().Format("20060102150405") + "-" + fmt.Sprintf("%06d", rand.Intn(1e6))
}

// 降级：直接查数据库（无锁）
func (l *GetProductLogic) getProductFromDB(id uint, emptyKey string) (*types.GetProductResp, error) {
	product, err := model.NewProductRepository(l.svcCtx.DB).FindByID(id)
	if err != nil {
		// 查无此商品，写入空值缓存防穿透
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.svcCtx.Redis.SetEx(l.ctx, emptyKey, "1", emptyExpire)
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("get product from db error: %w", err)
	}
	return l.buildProductResp(product), nil
}

// 组装响应
func (l *GetProductLogic) buildProductResp(product *model.Product) *types.GetProductResp {
	return &types.GetProductResp{
		Id:         product.ID,
		Name:       product.Name,
		Price:      product.Price,
		Attributes: product.Attributes,
	}
}
