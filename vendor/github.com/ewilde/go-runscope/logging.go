package runscope

import (
	"fmt"
	"strings"
)

var (
	debugHandler   func(level int, format string, args ...interface{})
	infoHandler    func(level int, format string, args ...interface{})
	errorHandler   func(level int, format string, args ...interface{})
	defaultHandler = func(level int, format string, args ...interface{}) {
		fmt.Printf("[DEBUG]%s %s\n", strings.Repeat("\t", level-1), fmt.Sprintf(format, args...))
	}
)

func init() {
	RegisterLogHandlers(defaultHandler, defaultHandler, defaultHandler)
}

func RegisterLogHandlers(
	debug func(level int, format string, args ...interface{}),
	info func(level int, format string, args ...interface{}),
	error func(level int, format string, args ...interface{})) {

	debugHandler = debug
	infoHandler = info
	errorHandler = error
}

func DebugF(level int, format string, args ...interface{}) {
	debugHandler(level, format, args...)
}

func InfoF(level int, format string, args ...interface{}) {
	infoHandler(level, format, args...)
}

func ErrorF(level int, format string, args ...interface{}) {
	errorHandler(level, format, args...)
}
