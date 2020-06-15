package main

import (
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
	// 监听端口
	listen, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal(err)
	}
	// 创建 grpc 服务
	s := grpc.NewServer()
	// 注册服务
	proto.RegisterHelloServiceServer(s, &service.HelloService{})

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

	// 启动服务
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	log.Println("server graceful stopped")
}
