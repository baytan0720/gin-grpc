package gin_grpc

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

// HandlerFunc is gin-grpc handler
type HandlerFunc func(*Context)

// HandlersChain is gin-grpc handler chain
type HandlersChain []HandlerFunc

// Engine is gin-grpc engine
type Engine struct {
	s   *grpc.Server
	Srv any

	handlers map[string]HandlersChain

	pool sync.Pool
}

func New(srv any) *Engine {
	e := &Engine{
		Srv: srv,

		handlers: make(map[string]HandlersChain),
	}

	e.s = grpc.NewServer(grpc.UnaryInterceptor(e.handleInterceptor()))
	e.pool.New = func() any {
		return e.allocateContext()
	}

	return e
}

// RegisterService implements grpc.ServiceRegistrar interface
func (engine *Engine) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	engine.s.RegisterService(sd, ss)
}

// Serve runs the grpc server
func (engine *Engine) Serve(l net.Listener) error {
	if err := engine.validateHandlers(); err != nil {
		return err
	}

	return engine.s.Serve(l)
}

func (engine *Engine) Run(addr ...string) error {
	if len(addr) == 0 {
		addr = []string{":8080"}
	}
	l, err := net.Listen("tcp", addr[0])
	if err != nil {
		return err
	}

	return engine.Serve(l)
}

// Handle registers a handler for the rpc
func (engine *Engine) Handle(funcName string, handler ...HandlerFunc) {
	engine.handlers[funcName] = append(engine.handlers[funcName], handler...)
}

// Use registers a handler for all rpc
func (engine *Engine) Use(handler ...HandlerFunc) {
	for funcName := range engine.handlers {
		engine.handlers[funcName] = append(engine.handlers[funcName], handler...)
	}
}

func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine}
}

func (engine *Engine) handleInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		split := strings.Split(info.FullMethod, "/")
		funcName := split[len(split)-1]
		handlers, ok := engine.handlers[funcName]
		if !ok {
			return handler(ctx, req)
		}

		c := engine.pool.Get().(*Context)
		defer engine.pool.Put(c)
		c.reset()
		c.ctx = ctx
		c.handlers = handlers
		c.Req = req

		c.Next()

		resp, _ = handler(ctx, req)

		if len(c.Errors) > 0 {
			return nil, c.Errors[len(c.Errors)-1]
		}

		if c.Resp == nil {
			return
		}

		if reflect.TypeOf(c.Resp) != reflect.TypeOf(resp) {
			return nil, fmt.Errorf("response type mismatch: %s != %s", reflect.TypeOf(c.Resp).Name(), reflect.TypeOf(resp).Name())
		}

		resp = c.Resp

		return
	}
}

func (engine *Engine) validateHandlers() error {
	for funcName, _ := range engine.handlers {
		if _, ok := reflect.TypeOf(engine.Srv).MethodByName(funcName); !ok {
			return fmt.Errorf("rpc func not found: %s", funcName)
		}
	}
	return nil
}
