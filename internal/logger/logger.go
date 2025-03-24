package logger

import (
	"go.uber.org/zap"
)

// 定义全局日志器
var Logger *zap.Logger

// 导出 zap.Field 的辅助函数
var (
	LogString = zap.String
	LogInt    = zap.Int
	LogBool   = zap.Bool
	LogError  = zap.Error
)

func InitLogger(env string) {
	var cfg zap.Config

	// 根据环境配置日志级别和格式
	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// 设置日志输出到标准错误流
	cfg.OutputPaths = []string{"stderr"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	// 构建日志器
	log, err := cfg.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	Logger = log
}

func init() {
	// 默认初始化为开发环境
	InitLogger("development")
}

// 快捷方法：Info
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// 快捷方法：Warn
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// 快捷方法：Error
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// 快捷方法：Debug
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// 快捷方法：Fatal
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
