package main

import (
	"flag"
	"shortvid-backend/app/shortvid-service/internal/conf"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

var (
	db  *gorm.DB
	err error
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../../internal/data/query",                   // 生成的代码路径
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface, // 生成模式
		FieldNullable: true,                                          // 数据库字段设置为空值列，当更新或插入时，设置为nil而不是零值
	})

	// 读取配置
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	switch bc.Data.Database.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(bc.Data.Database.Source), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("连接%s失败, err:%v", bc.Data.Database.Driver, err)
	}
	g.UseDB(db)
	g.ApplyBasic(
		g.GenerateAllTable()..., // 生成所有表
	)
	// 生成代码
	g.Execute()
}
