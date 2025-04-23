package common

import (
	"runtime"
	"strings"
)

func WithSafeFn(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			stackBuf := make([]byte, 1024)
			stackSize := runtime.Stack(stackBuf, false)
			stackTrace := strings.TrimSpace(string(stackBuf[:stackSize]))
			BaseLogger.WithField("panic", r).WithField("stack", stackTrace).Error("调用失败")
		}
	}()
	fn()
}
