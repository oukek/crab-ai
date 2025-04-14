package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 封装的通用返回
type Response struct {
	Code  int         `json:"code"`
	State bool        `json:"state"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func Success(c *gin.Context, msg string, data interface{}) {
	if c.GetHeader("Accept") == "text/event-stream" || c.GetHeader("accept") == "text/event-stream" {
		if c.Writer.Header().Get("Content-Type") != "text/event-stream" {
			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			c.Header("X-Accel-Buffering", "no") // Nginx 配置
		}
		_, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", "done", "done")
		if err != nil {
			BaseLogger.WithError(err).Error("stream将失败消息传递给客户端失败")
			return
		}
		c.Writer.Flush()
		return
	}
	c.JSON(200, Response{
		Code:  0,
		State: true,
		Msg:   msg,
		Data:  data,
	})
}

func Fail(c *gin.Context, msg string) {
	if c.GetHeader("Accept") == "text/event-stream" || c.GetHeader("accept") == "text/event-stream" {
		if c.Writer.Header().Get("Content-Type") != "text/event-stream" {
			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			c.Header("X-Accel-Buffering", "no") // Nginx 配置
		}
		_, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", "err", msg)
		if err != nil {
			BaseLogger.WithError(err).Error("stream将失败消息传递给客户端失败")
			return
		}
		c.Writer.Flush()
		return
	}
	c.JSON(200, Response{
		Code:  -1,
		State: false,
		Msg:   msg,
		Data:  nil,
	})
}

func Fail2(c *gin.Context, msg string, data interface{}) {
	if c.GetHeader("Accept") == "text/event-stream" || c.GetHeader("accept") == "text/event-stream" {
		if c.Writer.Header().Get("Content-Type") != "text/event-stream" {
			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			c.Header("X-Accel-Buffering", "no") // Nginx 配置
		}
		_, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", "err", msg)
		if err != nil {
			BaseLogger.WithError(err).Error("stream将失败消息传递给客户端失败")
			return
		}
		c.Writer.Flush()
		return
	}
	c.JSON(200, Response{
		Code:  -1,
		State: false,
		Msg:   msg,
		Data:  data,
	})
}

// ReLogin 重新登录，错误码-2
func ReLogin(c *gin.Context) {
	if c.GetHeader("Accept") == "text/event-stream" || c.GetHeader("accept") == "text/event-stream" {
		if c.Writer.Header().Get("Content-Type") != "text/event-stream" {
			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			c.Header("X-Accel-Buffering", "no") // Nginx 配置
		}
		_, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", "err", "relogin")
		if err != nil {
			BaseLogger.WithError(err).Error("stream将失败消息传递给客户端失败")
			return
		}
		c.Writer.Flush()
		return
	}
	c.JSON(200, Response{
		Code:  -2,
		State: false,
		Msg:   "请先进行登录",
		Data:  nil,
	})
}
