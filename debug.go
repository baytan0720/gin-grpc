package gin_grpc

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

const ginSupportMinGoVer = 18

// IsDebugging returns true if the framework is running in debug mode.
// Use SetMode(gin.ReleaseMode) to disable debug mode.
func IsDebugging() bool {
	return ginMode == debugCode
}

// DebugPrintHandlerFunc indicates debug log output format.
var DebugPrintHandlerFunc func(funcName, handlerName string, nuHandlers int)

func debugPrintHandler(funcName string, handlers HandlersChain) {
	if IsDebugging() {
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers.Last())
		if DebugPrintHandlerFunc == nil {
			debugPrint("%-25s --> %s (%d handlers)\n", funcName, handlerName, nuHandlers)
		} else {
			DebugPrintHandlerFunc(funcName, handlerName, nuHandlers)
		}
	}
}

func debugPrint(format string, values ...any) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(DefaultWriter, "[GIN-GRPC debug] "+format, values...)
	}
}

func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func debugPrintWARNINGDefault() {
	if v, e := getMinVer(runtime.Version()); e == nil && v < ginSupportMinGoVer {
		debugPrint(`[WARNING] Now Gin requires Go 1.18+.

`)
	}
	debugPrint(`[WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

`)
}

func debugPrintError(err error) {
	if err != nil && IsDebugging() {
		fmt.Fprintf(DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
	}
}
