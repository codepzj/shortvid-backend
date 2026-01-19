package logger

import (
	"fmt"
	"shortvid-backend/app/shortvid-service/internal/conf"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log    *zap.Logger
	msgKey string
}

type Option struct {
	Format           string
	Level            string
	FullLogFilename  string
	ErrorLogFilename string
	MaxSize          int32
	MaxBackups       int32
	MaxAge           int32
	Compress         bool
}

func NewOption(c *conf.Log) *Option {
	return &Option{
		Format:           c.Format,
		Level:            c.Level,
		FullLogFilename:  c.FullLogFilename,
		ErrorLogFilename: c.ErrorLogFilename,
		MaxSize:          c.MaxSize,
		MaxBackups:       c.MaxBackups,
		MaxAge:           c.MaxAge,
		Compress:         c.Compress,
	}
}

func NewZapLogger(opts *Option) *Logger {
	var encoder zapcore.Encoder
	if opts.Format == "text" {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	core := zapcore.NewCore(encoder, getFullLogWriter(opts), parseLogLevel(opts.Level))

	return &Logger{log: zap.New(core), msgKey: log.DefaultMessageKey}
}

func (l *Logger) Log(level log.Level, keyvals ...any) error {
	// If logging at this level is completely disabled, skip the overhead of
	// string formatting.
	if zapcore.Level(level) < zapcore.DPanicLevel && !l.log.Core().Enabled(zapcore.Level(level)) {
		return nil
	}
	var (
		msg    = ""
		keylen = len(keyvals)
	)
	if keylen == 0 || keylen%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		if keyvals[i].(string) == l.msgKey {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}
