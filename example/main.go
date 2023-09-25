package main

import (
	"context"
	"fmt"
	"net"

	gin "github.com/baytan0720/gin-grpc"
	"github.com/baytan0720/gin-grpc/example/proto"
)

func main() {
	panic(serve())
}

type Server struct {
	proto.UnimplementedHelloServer
}

func serve() error {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	e := gin.Default(&Server{})
	proto.RegisterHelloServer(e, e.Srv.(*Server))

	e.Handle("Ping", func(c *gin.Context) {
		var req *proto.PingReq
		c.BindRequest(&req)
		fmt.Printf("ID: %s\n", req.Id)
		c.Response(&proto.PongRes{
			Status: 1,
		})
	})

	return e.Serve(l)
}

func (s *Server) Ping(ctx context.Context, req *proto.PingReq) (*proto.PongRes, error) {
	return &proto.PongRes{}, nil
}
