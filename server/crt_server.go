package main

import (
	"context"
	"github/stone955/go-grpc/internal/auth"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github/stone955/go-grpc/internal/service"
	"github/stone955/go-grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	cred, err := credentials.NewServerTLSFromFile("../crt/server.crt", "../crt/server.key")
	if err != nil {
		log.Fatal(err)
	}

	// 创建 grpc 服务
	//s := grpc.NewServer(grpc.Creds(cred))

	// 添加截取器
	s := grpc.NewServer(grpc.Creds(cred), grpc.UnaryInterceptor(filter))

	// 注册服务
	proto.RegisterHelloServiceServer(s, &service.HelloService{
		// 加入身份验证
		Auth: &auth.Authentication{
			User:     "stone",
			Password: "123456",
		},
	})

	//  监听退出信号，优雅关闭
	var wg sync.WaitGroup
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)

	go func() {
		wg.Add(1)
		defer wg.Done()
		select {
		case <-ch:
			s.GracefulStop()
		}
	}()

	// 监听端口
	listen, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal(err)
	}
	// 启动服务
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	log.Println("server graceful stopped")
}

func filter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("filter  server= %v, fullmethod= %v\n", info.Server, info.FullMethod)
	return handler(ctx, req)
}
