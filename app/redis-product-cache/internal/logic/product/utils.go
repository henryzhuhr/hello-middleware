package product

import "time"

const (
	productCachePrefix = "product:" // 商品缓存前缀
	emptyCacheSuffix   = ":empty"   // 空值缓存后缀
)

const (
	baseExpire  = 30 * time.Minute // 基础过期时间
	emptyExpire = 5 * time.Minute  // 空值缓存过期时间

	mutexExpire = 3*time.Second // 互斥锁过期时间
)
