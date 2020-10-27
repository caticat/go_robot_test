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

	//ptrPic := newPic()
	//ptrPic.init(g_ptrConfig.Base)
	//ptrSubPic := newSubPic()
	//ptrSubPic.init(g_ptrConfig.PathPart, "2.jpg")
	//
	//sliPos := getSliPosInPic(ptrPic, ptrSubPic)
	//logPrintf("匹配点数量:%v", len(sliPos))
	//for _, ptrPos := range sliPos {
	//	logPrintf("匹配点坐标:(%v,%v)", ptrPos.m_x, ptrPos.m_y)
	//}
	//

	//baseF, e := os.Open(g_ptrConfig.Base)
	//failOnError(e, "open base fail")
	//defer baseF.Close()
	//
	//partF, e := os.Open(g_ptrConfig.Part)
	//failOnError(e, "open part fail")
	//defer partF.Close()
	//
	//baseI, _, e := image.Decode(baseF)
	//failOnError(e, "image decode fail base")
	////logPrintf(baseExt)
	//
	//partI, _, e := image.Decode(partF)
	//failOnError(e, "image decode fail part")
	////logPrintf(partExt)
	//
	//logPrintf("是否包含:%v", contain(baseI, partI))
}
