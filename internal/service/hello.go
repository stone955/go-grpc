package service

import (
	"context"
	"io"

	"github/stone955/go-grpc/proto"
)

type HelloService struct{}

func (srv *HelloService) Hello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	resp := proto.HelloResponse{
		FullName: "Si Dong yu",
	}
	return &resp, nil
}

func (srv *HelloService) HelloStream(stream proto.HelloService_HelloStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			// 客户端关闭，服务端关闭流
			if err == io.EOF {
				return nil
			}
			return err
		}

		resp := &proto.HelloResponse{FullName: req.Name + " Si"}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
