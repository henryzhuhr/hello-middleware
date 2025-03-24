package main

import (
	"github.com/henryzhuhr/hello-sql/internal/logger"
)

func main() {
	// 初始化日志器
	logger.InitLogger("development")

	// 记录一条日志
	logger.Info("This is a test log message",
		logger.LogString("key", "value"),
		logger.LogInt("status_code", 200),
	)
}
