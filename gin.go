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

// Engine is gin-grpc engine
type Engine struct {
	s   *grpc.Server
	Srv any

	handlers map[string]HandlersChain

	pool sync.Pool
}

// New creates a new gin-grpc engine
func New(srv any) *Engine {
	engine := &Engine{
		Srv: srv,

		handlers: make(map[string]HandlersChain),
	}

	engine.allocateHandlers()

	engine.s = grpc.NewServer(grpc.UnaryInterceptor(engine.handleInterceptor()))

	engine.pool.New = func() any {
		return engine.allocateContext()
	}

	return engine
}

// Default creates a new gin-grpc engine with the Recovery middleware already attached
func Default(srv any) *Engine {
	engine := New(srv)
	engine.Use(Recovery())
	return engine
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

// Run runs the grpc server on address
func (engine *Engine) Run(addr ...string) (err error) {
	defer func() { debugPrintError(err) }()

	address := resolveAddress(addr)

	var l net.Listener
	l, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}

	engine.printHandlers()
	debugPrint("Listening and serving HTTP on %s\n", address)
	err = engine.Serve(l)

	return
}

// Handle registers a handler for the rpc
func (engine *Engine) Handle(funcName string, handler ...HandlerFunc) {
	if _, ok := engine.handlers[funcName]; !ok {
		panic(fmt.Errorf("handle a non-exist rpc: %s", funcName))
	}

	engine.handlers[funcName] = append(engine.handlers[funcName], handler...)
}

func (engine *Engine) HandleFunc(handler ...HandlerFunc) {
	funcName := reflect.TypeOf(handler[len(handler)-1]).Name()
	engine.Handle(funcName, handler...)
}

// Use registers a middleware for all rpc
func (engine *Engine) Use(middleware ...HandlerFunc) {
	for funcName := range engine.handlers {
		engine.handlers[funcName] = append(engine.handlers[funcName], middleware...)
	}
}

func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine, Req: &Request{}, Resp: &Response{}}
}

func (engine *Engine) handleInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		split := strings.Split(info.FullMethod, "/")
		funcName := split[len(split)-1]
		handlers, ok := engine.handlers[funcName]
		if !ok {
			return handler(ctx, req)
		}

		resp, _ = handler(ctx, req)

		c := engine.pool.Get().(*Context)
		defer engine.pool.Put(c)
		c.reset()
		c.ctx = ctx
		c.handlers = handlers
		c.Req.req = req
		c.Resp.resp = resp

		c.Next()

		if len(c.Errors) > 0 {
			return nil, c.Errors[len(c.Errors)-1]
		}

		return c.Resp.resp, nil
	}
}

func (engine *Engine) getFuncNames() []string {
	var funcNames []string
	for i := 0; i < reflect.TypeOf(engine.Srv).NumMethod(); i++ {
		funcNames = append(funcNames, reflect.TypeOf(engine.Srv).Method(i).Name)
	}
	return funcNames
}

func (engine *Engine) allocateHandlers() {
	for _, funcName := range engine.getFuncNames() {
		engine.handlers[funcName] = HandlersChain{}
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

func (engine *Engine) printHandlers() {
	for funcName, handlers := range engine.handlers {
		debugPrintHandler(funcName, handlers)
	}
}

// HandlersChain is gin-grpc handler chain
type HandlersChain []HandlerFunc

func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
