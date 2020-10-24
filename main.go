package main

import (
	"image"
	_ "image/jpeg"
	"os"
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

	baseF, e := os.Open(g_ptrConfig.Base)
	failOnError(e, "open base fail")
	defer baseF.Close()

	partF, e := os.Open(g_ptrConfig.Part)
	failOnError(e, "open part fail")
	defer partF.Close()

	baseI, _, e := image.Decode(baseF)
	failOnError(e, "image decode fail base")
	//logPrintf(baseExt)

	partI, _, e := image.Decode(partF)
	failOnError(e, "image decode fail part")
	//logPrintf(partExt)

	logPrintf("是否包含:%v", contain(baseI, partI))
}
