package gin_grpc

import (
	"os"
	"reflect"
	"runtime"
)

func nameOfFunction(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

// H is a shortcut for map[string]any
type H map[string]any

// StoreRequestIntoKeys stores the all params of request into keys.
// It better be a flattening struct, not a nested struct.
// The keys are the field names of the struct.
func StoreRequestIntoKeys() HandlerFunc {
	return func(c *Context) {
		for k, v := range c.Req.getAllFields() {
			c.Set(k, v)
		}
	}
}
