package service

import (
	"context"

	"github/stone955/go-grpc/proto"
)

type HelloService struct{}

func (srv *HelloService) Hello(context.Context, *proto.HelloRequest) (*proto.HelloResponse, error) {
	resp := proto.HelloResponse{
		FullName: "stone",
	}
	return &resp, nil
}
