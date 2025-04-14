package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"oukek/crab-ai/app/common"
	"runtime/debug"
)

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {

			if e, ok := r.(common.Error); ok {
				if e.Log == true {
					common.BaseLogger.WithFields(logrus.Fields{
						"stack": string(debug.Stack()),
						"pre":   e.Pre,
						"panic": r,
					}).Error(e.Msg)
					debug.PrintStack()
				}

				if c.GetHeader("Accept") == "text/event-stream" || c.GetHeader("accept") == "text/event-stream" {
					if c.Writer.Header().Get("Content-Type") != "text/event-stream" {
						c.Header("Content-Type", "text/event-stream")
						c.Header("Cache-Control", "no-cache")
						c.Header("Connection", "keep-alive")
						c.Header("X-Accel-Buffering", "no") // Nginx 配置
					}
					_, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", "err", e.Msg)
					if err != nil {
						fmt.Printf("stream将失败消息传递给客户端失败:%v\r\n", err)
						return
					}
					c.Writer.Flush()
					return
				} else {
					c.JSON(http.StatusOK, common.Response{
						Code:  e.Code,
						State: false,
						Msg:   e.Msg,
						Data:  nil,
					})
				}

			} else {
				if os.Getenv("ENV") == "local" {
					panic(r)
				}
				//打印错误堆栈信息
				common.BaseLogger.WithFields(logrus.Fields{
					"stack": string(debug.Stack()),
					"panic": r,
				}).Error("服务器异常")
				if c.GetHeader("Accept") == "text/event-stream" || c.GetHeader("accept") == "text/event-stream" {
					if c.Writer.Header().Get("Content-Type") != "text/event-stream" {
						c.Header("Content-Type", "text/event-stream")
						c.Header("Cache-Control", "no-cache")
						c.Header("Connection", "keep-alive")
						c.Header("X-Accel-Buffering", "no") // Nginx 配置
					}
					_, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", "err", e.Msg)
					if err != nil {
						fmt.Printf("stream将失败消息传递给客户端失败:%v\r\n", err)
						return
					}
					c.Writer.Flush()
					return
				} else {
					c.JSON(http.StatusOK, common.Response{
						Code:  -1,
						State: false,
						Msg:   "服务器异常，请重试",
						Data:  nil,
					})
				}

			}
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
