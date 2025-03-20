package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379", //"localhost:6379",
		Password: "",                  // 没有密码，默认值
		DB:       0,                   // 默认DB 0
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Errorf("failed to ping redis, error is %s", err)
		return
	}

	// Set a key with a value
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

}
