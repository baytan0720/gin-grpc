package main

import (
	"context"
	"fmt"

	gin "github.com/baytan0720/gin-grpc"
	"github.com/baytan0720/gin-grpc/example/proto"
)

type Server struct {
	proto.UnimplementedHelloServer
}

func main() {
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

	panic(e.Run()) // listen and serve on 8080
}

func (s *Server) Ping(ctx context.Context, req *proto.PingReq) (*proto.PongRes, error) {
	return &proto.PongRes{}, nil
}
