package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "shortvid-backend/api/shortvid-service/v1"
)

var (
	name = flag.String("name", "codepzj", "eg: -name=pzj")
)

func main() {
	flag.Parse()
	// 1. 创建连接 无需TLS认证
	conn, err := grpc.Dial("127.0.0.1:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// 2. 创建client对象
	client := v1.NewUsersServiceClient(conn)

	// 3. 调用方法
	resp, err := client.GetUser(context.Background(), &v1.GetUserRequest{Id: 1})
	if err != nil {
		log.Fatalf("failed to call GetUser: %v", err)
	}
	log.Printf("Response: %s", resp.Nickname)
}