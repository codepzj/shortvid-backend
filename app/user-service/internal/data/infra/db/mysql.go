package db

import (
	"log"
	"os"

	"shortvid-backend/app/user-service/internal/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDB(c *conf.Data) *gorm.DB {
	log.Println("MySQL connect start...")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: c.Mysql.SlowSqlThreshold.AsDuration(), // 慢查询阈值
			LogLevel:      logger.Info,                           // 日志级别
		},
	)
	// 初始化gorm
	db, err := gorm.Open(mysql.New(
		mysql.Config{DSN: c.Mysql.Dsn}),
		&gorm.Config{Logger: newLogger, NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}},
	)

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(int(c.Mysql.MaxOpenConns)) // 设置最大打开连接数(最多允许n个连接去使用mysql, 超过会被阻塞)
	sqlDB.SetMaxIdleConns(int(c.Mysql.MaxIdleConns)) // 设置最大空闲连接数(最多允许创建n个连接但是不使用mysql的, 超过会被丢弃, 防止频繁创建对象造成的开销)

	if err != nil {
		log.Fatalf("Connect MySQL failed: %v", err)
	}
	log.Println("MySQL connect success...")
	return db
}
