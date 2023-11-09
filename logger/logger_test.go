package logger

import (
	"testing"
)

var setting = Settings{
	Filename:    "im.log",
	Level:       "debug",
	RollingDays: 1024,
	Format:      "text",
}

func TestInitLog(t *testing.T) {
	Init(setting)
	firstLog := "I am the first log"
	StdLog().Infof("hello world %s", firstLog)
}
