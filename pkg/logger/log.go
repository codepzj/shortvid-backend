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
func getLogWriter(opt *Option) zapcore.WriteSyncer {
	logger := &lumberjack.Logger{
		Filename:   opt.LogFile,
		MaxSize:    int(opt.MaxSize),
		MaxBackups: int(opt.MaxBackups),
		MaxAge:     int(opt.MaxAge),
		Compress:   opt.Compress,
	}
	ws := io.MultiWriter(logger, os.Stdout)
	return zapcore.AddSync(ws)
}