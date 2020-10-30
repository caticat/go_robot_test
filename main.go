package main

import (
	_ "image/jpeg"
)

var (
	g_ptrConfig = newConfig("config.yaml")
)

func main() {
	logPrintf("程序开始")

	test()

	logPrintf("程序结束")
}

func test() {
	ptrPicMgr := newPicMgr()
	defer ptrPicMgr.close()
	ptrPicMgr.load(g_ptrConfig.PathPart, getDirFile(g_ptrConfig.PathPart, g_ptrConfig.PartExt))
	ptrPicMgr.calc(g_ptrConfig.Base)
}
