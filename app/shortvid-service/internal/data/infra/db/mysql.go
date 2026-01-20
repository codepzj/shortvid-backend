package db

import (
	"log"
	"os"

	"shortvid-backend/app/shortvid-service/internal/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDB(c *conf.Data) *gorm.DB {
	newLogger := logger.New(
		log.New(getLogOutput(c.Database.LogFormat, c.Database.LogFile), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: c.Database.SlowSqlThreshold.AsDuration(), // 慢查询阈值
			LogLevel:      parseMysqlLogLevel(c.Database.LogLevel),  // 日志级别
		},
	)
	// 初始化gorm
	db, err := gorm.Open(mysql.New(
		mysql.Config{DSN: c.Database.Source}),
		&gorm.Config{Logger: newLogger, NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}},
	)
	if err != nil {
		log.Fatalf("Connect MySQL failed: %v", err)
		panic(err)
	}
	log.Printf("MySQL connect success...")
	return db
}

func getLogOutput(format string, logFile string) *os.File {
	switch format {
	case "text":
		return os.Stdout
	case "json":
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Open log file failed: %v", err)
		}
		return f
	}
	return os.Stdout
}

func parseMysqlLogLevel(logLevel string) logger.LogLevel {
	switch logLevel {
	case "debug":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Info
	}
}
