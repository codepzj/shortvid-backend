package logger

import (
	"io"
	"os"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}

// 全量日志
func getFullLogWriter(option *Option) zapcore.WriteSyncer {
	fullLogger := &lumberjack.Logger{
		Filename:   option.FullLogFilename, // 日志文件路径
		MaxSize:    int(option.MaxSize),    // 单个文件最大 10MB
		MaxBackups: int(option.MaxBackups), // 保留 5 个旧文件
		MaxAge:     int(option.MaxAge),     // 保留 30 天
		Compress:   option.Compress,        // 启用压缩
	}
	ws := io.MultiWriter(fullLogger, os.Stdout)
	return zapcore.AddSync(ws)
}
