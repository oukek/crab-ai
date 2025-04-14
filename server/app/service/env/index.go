package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var isInit = false

func Get(key string) string {
	if !isInit {
		isInit = true
		InitEnv()
	}
	return os.Getenv(key)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(directory string) string {
	return substr(directory, 0, strings.LastIndex(directory, string(filepath.Separator)))
}

func InitEnv() {
	cur, err := os.Executable()
	if err != nil {
		log.Fatal("获取可执行目录失败", err)
	}
	cur = filepath.Dir(cur)
	log.Println("当前的系统平台：", runtime.GOOS)
	log.Printf("启动参数为：%v\n", os.Args)
	if runtime.GOOS != "windows" {
		log.Println("切换到工作目录：", cur)
		err = os.Chdir(cur)
		if err != nil {
			log.Fatal("切换到工作目录失败", err)
		}
	}

	log.Println("读取环境变量")
	// 读取环境变量
	// 获取当前的工作目录，层层遍历上去，直到找到.env文件
	cur, err = os.Getwd()
	log.Println("获取环境变量目录：", cur)
	if err != nil {
		log.Fatal("获取当前工作目录失败", err)
	}
	cur = filepath.Clean(cur)
	// 遍历
	for {
		// 判断文件是否存在
		_, err := os.Stat(cur + "/.env")
		if err == nil {
			// 读取环境变量
			err = godotenv.Overload(cur + string(filepath.Separator) + ".env")
			if err != nil {
				log.Fatal("Error loading .env file", err)
			}
			log.Println("读取环境变量成功：", cur+string(filepath.Separator)+".env")
			break
		}
		// 判断是否到了根目录
		if cur == "/" {
			log.Fatal("未找到.env文件")
		}
		// 获取父级目录
		cur = getParentDirectory(cur)
	}
}
