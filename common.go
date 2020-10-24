package main

import (
	"fmt"
	"log"
	"runtime"
)

func failOnError(err error, msg string) {
	if err != nil {
		logFatalf(2, "%s: %s", msg, err)
	}
}

func logPrintf(format string, args ...interface{}) {
	logNormal(2, format, args...)
}

func logNormal(callSkip int, format string, args ...interface{}) {
	_, src, line, _ := runtime.Caller(callSkip)
	log.Printf(fmt.Sprintf("%v:%v ", src, line)+format, args...)
}

func logFatalf(callSkip int, format string, args ...interface{}) {
	_, src, line, _ := runtime.Caller(callSkip)
	log.Fatalf(fmt.Sprintf("%v:%v ", src, line)+format, args...)
}
