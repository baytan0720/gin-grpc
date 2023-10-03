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

	e.Use(gin.StoreRequestIntoKeys())

	e.Handle("Ping", func(c *gin.Context) {
		var req *proto.PingReq
		c.BindRequest(&req)
		fmt.Printf("ID: %s\n", req.Id)
		fmt.Printf("ID: %s\n", c.GetString("Id"))
		fmt.Printf("ID: %s\n", c.Req.GetField("Id"))
		c.Response(&proto.PongRes{
			Status: 1,
		})

		// another way to response:
		//
		//c.ResponseField("Status", int32(1))
		//
		//c.ResponseFields(gin.H{
		//	"Status": int32(1),
		//})
	})

	e.Run() // listen and serve on 8080
}

func (s *Server) Ping(ctx context.Context, req *proto.PingReq) (*proto.PongRes, error) {
	return &proto.PongRes{}, nil
}
