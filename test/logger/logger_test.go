package main

import (
	"bytes"
	"testing"

	"github.com/henryzhuhr/hello-sql/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	// 创建一个内存缓冲区用于捕获日志输出
	var buf bytes.Buffer

	// 配置 zap 将日志输出到缓冲区
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(&buf),
		zap.DebugLevel,
	)
	logger.Logger = zap.New(core)

	// 记录一条日志
	logger.Info("This is a test log message",
		logger.LogString("key", "value"),
		logger.LogInt("status_code", 200),
	)

	// // 检查日志内容
	// logOutput := buf.String()
	// expectedMessage := "This is a test log message"
	// if !bytes.Contains(buf.Bytes(), []byte(expectedMessage)) {
	// 	t.Errorf("Expected log to contain %q, but got:\n%s", expectedMessage, logOutput)
	// }

	// // 检查结构化字段
	// if !bytes.Contains(buf.Bytes(), []byte(`"key":"value"`)) {
	// 	t.Errorf("Expected log to contain key-value pair %q, but got:\n%s", `"key":"value"`, logOutput)
	// }
	// if !bytes.Contains(buf.Bytes(), []byte(`"status_code":200`)) {
	// 	t.Errorf("Expected log to contain key-value pair %q, but got:\n%s", `"status_code":200`, logOutput)
	// }
}
