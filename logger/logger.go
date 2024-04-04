package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	File   *os.File
	Logger *log.Logger
)

func init() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("无法打开日志文件: ", err)
	}
	File = logFile
	Logger = log.New(logFile, "", log.LstdFlags)
}

func Println(v ...any) {
	Logger.Println(v)
}

func Close() {
	if File == nil {
		return
	}
	if err := File.Sync(); err != nil {
		fmt.Println(err.Error())
	}
	if err := File.Close(); err != nil {
		fmt.Println(err.Error())
	}
}
