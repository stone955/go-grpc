package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github/stone955/go-grpc/internal/auth"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"github/stone955/go-grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	certificate, err := tls.LoadX509KeyPair("../crt/client.crt", "../crt/client.key")
	if err != nil {
		log.Fatal(err)
	}

	ca, err := ioutil.ReadFile("../crt/ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal(err)
	}

	cred := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
		ServerName:   "localhost",
	})

	// 客户端信息
	a := &auth.Authentication{
		User:     "stone",
		Password: "123456",
	}

	// 建立连接
	conn, err := grpc.Dial("localhost:8001", grpc.WithTransportCredentials(cred), grpc.WithPerRPCCredentials(a))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = conn.Close()
	}()

	// 获取客户端
	client := proto.NewHelloServiceClient(conn)

	// 单向流调用
	request := &proto.HelloRequest{
		Name: "sdy",
	}

	resp, err := client.Hello(context.TODO(), request)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("response fullname: %v\n", resp.FullName)

	// 双向流调用
	// 获取流
	stream, err := client.HelloStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 确保发送线程，接收线程都结束才退出
	var wg sync.WaitGroup
	wg.Add(2)

	// 发送线程
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			if err := stream.Send(&proto.HelloRequest{Name: "stone "}); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// 接收线程
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			log.Printf("stream response fullname: %v\n", resp.FullName)
		}
	}()

	wg.Wait()
}
