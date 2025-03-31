package plugin

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rdb      *redis.Client // 包级变量，存储 Redis 客户端实例
	rdb_once sync.Once     // 确保 Redis 客户端只被创建一次
)

// RedisConfig 定义 Redis 的连接配置
type RedisConfig struct {
	Addr     string // Redis server address
	Password string // Redis password (optional)
	DB       int    // Redis database number
}

// NewRedisClient 使用单例模式初始化Redis客户端
func NewRedisClient() *redis.Client {
	rdb_once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     "redis-server:6379", // Redis server address
			Password: "",                  // No password set
			DB:       0,                   // Use default DB
		})
	})
	return rdb
}
