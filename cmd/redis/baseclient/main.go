// https://redis.uptrace.dev/guide/go-redis.html
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/uuid"
	"github.com/henryzhuhr/hello-middleware/internal/plugin"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	rdb := plugin.NewRedisClient()
	if rdb == nil {
		log.Error("failed to create redis client")
		return
	}

	// Set a key with a value.
	// Every Redis command accepts a context that you can use to set timeouts 
	// or propagate some information, for example, tracing context.
	ctx := context.Background()

	uid := uuid.New().String()
	keyName := fmt.Sprintf("USER_%s", uid[:8])
	keyValue := uid
	if err := rdb.Set(ctx, keyName, keyValue, 0).Err(); err != nil {
		log.Errorf("failed to set key, error is %s", err)
	}
	val, err := rdb.Get(ctx, keyName).Result()
	switch {
	case err == redis.Nil:
		log.Warnf("key %s does not exist", keyName)
	case err == nil && val != "":
		log.Infof("key %s exists, value is %s", keyName, val)
	case err != nil:
		log.Errorf("failed to get key, error is %s", err)
	}
	if err != nil {
		log.Errorf("failed to get key, error is %s", err)
	}
	log.Infof("[%s] %s", keyName, val)

	get := rdb.Get(ctx, "abc")
	log.Infof("val = %s, err = %v", get.Val(), get.Err())

	RedisListDemo(rdb)
}

// RedisListDemo 列表的使用
func RedisListDemo(rdb *redis.Client) {
	// List operations
	ctx := context.Background()

	// Push values to the list
	if err := rdb.LPush(ctx, "mylist", "value1", "value2").Err(); err != nil {
		log.Errorf("failed to lpush, error is %s", err)
	}

	// Pop a value from the list
	val, err := rdb.LPop(ctx, "mylist").Result()
	if err != nil {
		log.Errorf("failed to lpop, error is %s", err)
	} else {
		log.Infof("popped value: %s", val)
	}
}
