// 电商平台商品缓存
package main

import (
	"github.com/henryzhuhr/hello-middleware/internal/plugin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Initialize the Redis client
	rdb := plugin.NewRedisClient()
	if rdb == nil {
		log.Error("failed to create redis client")
		return
	}

}
