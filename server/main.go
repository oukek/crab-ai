package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"oukek/crab-ai/app/service/env"
	"oukek/crab-ai/cmd"
)

func main() {
	// 初始化环境变量
	env.InitEnv()

	gin.SetMode(os.Getenv("GIN_MODE"))

	cmd.Execute()

}
