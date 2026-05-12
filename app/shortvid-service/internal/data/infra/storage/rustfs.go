package storage

import (
	"log"
	"shortvid-backend/app/shortvid-service/internal/conf"

	"github.com/Scorpio69t/rustfs-go"
	"github.com/Scorpio69t/rustfs-go/pkg/credentials"
)

func NewRustFS(conf *conf.RustFs) *rustfs.Client {
	// 初始化客户端
	client, err := rustfs.New(conf.Endpoint, &rustfs.Options{
		Credentials: credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure:      conf.UseSsl, // 设置为 true 使用 HTTPS
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("RustFs connect success...")
	return client
}
