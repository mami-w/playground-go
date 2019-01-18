package logger

import (
	"log"
	"os"
)

var std = log.New(os.Stderr, "", log.LstdFlags)

func Get() (logger *log.Logger) {
	return std
}