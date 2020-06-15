package main

import (
	"context"
	"github/stone955/go-grpc/proto"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// 建立连接
	conn, err := grpc.Dial("localhost:8001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = conn.Close()
	}()

	client := proto.NewHelloServiceClient(conn)

	request := &proto.HelloRequest{
		Name: "sdy",
	}

	resp, err := client.Hello(context.TODO(), request)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("response fullname: %v\n", resp.FullName)
}
