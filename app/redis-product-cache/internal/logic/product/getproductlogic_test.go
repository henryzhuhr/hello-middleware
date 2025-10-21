package product

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/model"
	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/svc"
	"github.com/henryzhuhr/hello-middleware/app/redis-product-cache/internal/types"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newTestSvcCtx creates a ServiceContext backed by in-memory sqlite and miniredis
func newTestSvcCtx(t *testing.T) (*svc.ServiceContext, *miniredis.Miniredis) {
	t.Helper()

	// Redis: miniredis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	// SQLite in-memory
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.Product{}, &model.ProductSKU{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return &svc.ServiceContext{DB: db, Redis: rdb}, mr
}

func newLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return NewGetProductLogic(ctx, svcCtx)
}

func mustCreateProduct(t *testing.T, db *gorm.DB, id uint, name string, price float64) *model.Product {
	t.Helper()
	p := &model.Product{
		Model:       gorm.Model{ID: id},
		Name:        name,
		Description: "desc",
		CategoryID:  1,
		Brand:       "brand",
		Price:       price,
		Attributes:  model.JSONMap{"color": "red"},
	}
	if err := db.Create(p).Error; err != nil {
		t.Fatalf("create product failed: %v", err)
	}
	return p
}

func TestGetProduct_CacheHit(t *testing.T) {
	svcCtx, mr := newTestSvcCtx(t)
	defer mr.Close()

	ctx := context.Background()
	l := newLogic(ctx, svcCtx)

	// Prepare cached value
	id := uint(101)
	cacheKey := productCachePrefix + "101"
	cached := types.GetProductResp{Id: id, Name: "cachedName", Price: 9.99, Attributes: map[string]string{"color": "blue"}}
	b, _ := json.Marshal(cached)
	if err := svcCtx.Redis.Set(ctx, cacheKey, b, 10*time.Minute).Err(); err != nil {
		t.Fatalf("prep cache failed: %v", err)
	}

	// Call
	resp, err := l.GetProduct(&types.GetProductReq{Id: id})
	if err != nil {
		t.Fatalf("GetProduct error: %v", err)
	}
	if resp == nil || resp.Name != "cachedName" || resp.Price != 9.99 {
		t.Fatalf("unexpected resp: %+v", resp)
	}
}

func TestGetProduct_EmptyMarkerHit(t *testing.T) {
	svcCtx, mr := newTestSvcCtx(t)
	defer mr.Close()

	ctx := context.Background()
	l := newLogic(ctx, svcCtx)

	id := uint(202)
	cacheKey := productCachePrefix + "202"
	emptyKey := cacheKey + emptyCacheSuffix
	if err := svcCtx.Redis.Set(ctx, emptyKey, "1", 5*time.Minute).Err(); err != nil {
		t.Fatalf("prep empty marker failed: %v", err)
	}

	resp, err := l.GetProduct(&types.GetProductReq{Id: id})
	if err == nil || !strings.Contains(err.Error(), "product not found") {
		t.Fatalf("expected product not found, got resp=%+v err=%v", resp, err)
	}
}

func TestGetProduct_DBFillThenCache(t *testing.T) {
	svcCtx, mr := newTestSvcCtx(t)
	defer mr.Close()

	ctx := context.Background()
	l := newLogic(ctx, svcCtx)

	id := uint(303)
	mustCreateProduct(t, svcCtx.DB, id, "dbName", 12.34)

	resp, err := l.GetProduct(&types.GetProductReq{Id: id})
	if err != nil {
		t.Fatalf("GetProduct error: %v", err)
	}
	if resp == nil || resp.Name != "dbName" || resp.Price != 12.34 {
		t.Fatalf("unexpected resp: %+v", resp)
	}

	// Verify cache populated
	cacheKey := productCachePrefix + "303"
	got, err := svcCtx.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		t.Fatalf("expected cache set, got err: %v", err)
	}
	var cached types.GetProductResp
	if err := json.Unmarshal([]byte(got), &cached); err != nil {
		t.Fatalf("cached json invalid: %v", err)
	}
	if cached.Name != "dbName" || cached.Price != 12.34 {
		t.Fatalf("cached value mismatch: %+v", cached)
	}

	// TTL should be > 25m (base 30m +/-), just ensure non-trivial expiry
	ttl := mr.TTL(cacheKey)
	if ttl <= 25*time.Minute {
		t.Fatalf("unexpected TTL: %v", ttl)
	}
}

func TestGetProduct_NotFound_WritesEmptyCache(t *testing.T) {
	svcCtx, mr := newTestSvcCtx(t)
	defer mr.Close()

	ctx := context.Background()
	l := newLogic(ctx, svcCtx)

	id := uint(404)
	resp, err := l.GetProduct(&types.GetProductReq{Id: id})
	if err == nil || !strings.Contains(err.Error(), "product not found") {
		t.Fatalf("expected product not found, got resp=%+v err=%v", resp, err)
	}

	emptyKey := productCachePrefix + "404" + emptyCacheSuffix
	v, err := mr.Get(emptyKey)
	if err != nil {
		t.Fatalf("expected empty marker key exists: %v", err)
	}
	if v != "1" {
		t.Fatalf("expected empty marker set to '1', got %q", v)
	}
	if mr.TTL(emptyKey) <= 0 {
		t.Fatalf("expected empty marker TTL > 0")
	}
}

func TestGetProduct_CorruptedCache_FallbackToDB(t *testing.T) {
	svcCtx, mr := newTestSvcCtx(t)
	defer mr.Close()

	ctx := context.Background()
	l := newLogic(ctx, svcCtx)

	id := uint(505)
	mustCreateProduct(t, svcCtx.DB, id, "db505", 55.5)

	cacheKey := productCachePrefix + "505"
	// write invalid JSON
	if err := svcCtx.Redis.Set(ctx, cacheKey, "{bad json", time.Minute).Err(); err != nil {
		t.Fatalf("prep bad cache failed: %v", err)
	}

	resp, err := l.GetProduct(&types.GetProductReq{Id: id})
	if err != nil {
		t.Fatalf("GetProduct error: %v", err)
	}
	if resp.Name != "db505" || resp.Price != 55.5 {
		t.Fatalf("unexpected resp: %+v", resp)
	}

	// now cache should be overwritten to valid json
	got, err := svcCtx.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		t.Fatalf("cache get failed: %v", err)
	}
	var cached types.GetProductResp
	if err := json.Unmarshal([]byte(got), &cached); err != nil {
		t.Fatalf("still invalid cache: %v", err)
	}
}
