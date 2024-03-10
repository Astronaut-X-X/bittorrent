package logger

import (
	"log"
	"os"
)

var (
	File   *os.File
	Logger *log.Logger
)

func Println(v ...any) {
	Logger.Println(v)
}
