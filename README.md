# Gin-gRPC

A framework that uses gRPC like gin

## Getting started

### Prerequisites

- **[Go](https://go.dev/)**: any one of the **three latest major** [releases](https://go.dev/doc/devel/release).

### Getting Gin-gRPC

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/baytan0720/gin-grpc"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `gin-grpc` package:

```
$ go get -u github.com/baytan0720/gin-grpc
```

### Running Gin-gRPC

First you need to create a proto file, one simplest example likes the follow `example/hello.proto`:

```
syntax = "proto3";

package proto;
option go_package="./proto";

message PingReq {
  string id = 1;
}

message PongRes {
  int32 status = 1;
}

service Hello {
  rpc Ping(PingReq) returns (PongRes);
}
```

Then you need to generate the code by using `protoc` command:

```
$ protoc -I . hello.proto --go_out=. --go-grpc_out=.
```

Or use script `example/mk-proto.sh`:

```
$ sh mk-proto.sh
```

Next you need to create a go file to code, one example likes the follow `example/main.go`:

```go
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

	panic(e.Run()) // listen and serve on 0.0.0.0:8080
}

func (s *Server) Ping(ctx context.Context, req *proto.PingReq) (*proto.PongRes, error) {
	return &proto.PongRes{}, nil
}
```

You just need to define a struct and implement the interface `proto.HelloServer` without any logic.

Then use `gin.Default()` to create a Gin-gRPC engine and must call `proto.RegisterHelloServer()` to register the server.

Now you can use `e.Handle()` to handle the rpc, the first parameter is the method name, the second parameter is the handler function.

Finally, use the Go command to run the demo:

```
$ go run main.go
```

You can use any api tools to test the rpc or make a client to call rpc.

## Contributing

We encourage you to contribute to Gin-gRPC!