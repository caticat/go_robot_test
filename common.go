package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
)

func failOnError(err error, format string, args ...interface{}) {
	if err != nil {
		logFatalf(2, "%s: %s", fmt.Sprintf(format, args...), err)
	}
}

func logPrintf(format string, args ...interface{}) {
	logNormal(2, format, args...)
}

func assert(condition bool, format string, args ...interface{}) {
	if !condition {
		logFatalf(2, "assert failed: %s", fmt.Sprintf(format, args...))
	}
}

func logNormal(callSkip int, format string, args ...interface{}) {
	_, src, line, _ := runtime.Caller(callSkip)
	log.Printf(fmt.Sprintf("%v:%v ", src, line)+format, args...)
}

func logFatalf(callSkip int, format string, args ...interface{}) {
	_, src, line, _ := runtime.Caller(callSkip)
	log.Fatalf(fmt.Sprintf("%v:%v ", src, line)+format, args...)
}

// 获取目录内的文件
func getDirFile(dir string, ext string) (ret []string) {
	sliFileInfo, e := ioutil.ReadDir(dir)
	failOnError(e, "read dir fail")
	//pathSeparator := string(os.PathSeparator)
	for _, fileInfo := range sliFileInfo {
		if fileInfo.IsDir() {
			continue
		}
		fileName := fileInfo.Name()
		if ext != "" {
			if len(fileName) <= len(ext) {
				continue
			}
			//logPrintf("a:%v,b:%v", fileName[len(fileName)-len(ext):], ext)
			if fileName[len(fileName)-len(ext):] != ext {
				continue
			}
		}
		ret = append(ret, fileName)
	}

	return
}
