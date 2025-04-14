package common

import (
	"errors"
	"fmt"
	"os"
)

// Error 错误检查
// 一般是在service层和model层调用，在controller层直接返回即可
type Error struct {
	Code int
	Msg  string
	Pre  error
	Log  bool
}

var ErrUnKnow = errors.New("服务器异常，请联系客服处理")

func SimpleError(msg string, args ...interface{}) {
	log := false
	if msg == "" {
		log = true
		msg = "服务器异常，请重试"
	}
	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args)
	}
	panic(Error{
		Code: -1,
		Msg:  msg,
		Pre:  nil,
		Log:  log,
	})
}

func SimpleCheck(err error) {
	if err != nil {
		if os.Getenv("ENV") == "local" {
			fmt.Printf("%v", err)
		}
		panic(Error{
			Code: -1,
			Msg:  "服务器异常，请重试",
			Pre:  err,
			Log:  true,
		})
	}
}

func Check(err error, code int, format string, args ...interface{}) {
	if err != nil {
		panic(Error{
			Code: code,
			Msg:  fmt.Sprintf(format, args...),
			Pre:  err,
			Log:  true,
		})
	}
}

var ErrNormal = errors.New("服务异常，请联系管理员")
