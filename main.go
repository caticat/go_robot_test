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
	//testGetFirstExactPos()

	logPrintf("程序结束")
}

func test() {
	ptrPicMgr := newPicMgr()
	defer ptrPicMgr.close()
	ptrPicMgr.load(g_ptrConfig.PathPart, getDirFile(g_ptrConfig.PathPart, g_ptrConfig.PartExt))
	ptrPicMgr.calc(g_ptrConfig.Base)
}

// 测试获取与子图片最相似的第一个点的坐标
func testGetFirstExactPos() {
	ptrPicMgr := newPicMgr()
	defer ptrPicMgr.close()
	ptrPicMgr.load(g_ptrConfig.PathPart, getDirFile(g_ptrConfig.PathPart, g_ptrConfig.PartExt))

	ptrPic := newPic()
	ptrPic.init(g_ptrConfig.Base)
	getFirstPos(ptrPic, ptrPicMgr.m_mapPic[1])
}

// 测试指定点的相似度
func testPosSimilar() {
	ptrPicMgr := newPicMgr()
	defer ptrPicMgr.close()
	ptrPicMgr.load(g_ptrConfig.PathPart, getDirFile(g_ptrConfig.PathPart, g_ptrConfig.PartExt))

	ptrPic := newPic()
	ptrPic.init(g_ptrConfig.Base)
	ptrPos := newPos()
	ptrPos.init(38, 25)
	testPosInPic(ptrPic, ptrPicMgr.m_mapPic[1], ptrPos)
}

// 获取子图片相关的所有坐标点
func testGetAllSubPicPos() {
	ptrPic := newPic()
	ptrPic.init(g_ptrConfig.Base)
	ptrSubPic := newSubPic()
	ptrSubPic.init(g_ptrConfig.PathPart, "2.jpg")

	sliPos := getSliPosInPic(ptrPic, ptrSubPic)
	logPrintf("匹配点数量:%v", len(sliPos))
	for _, ptrPos := range sliPos {
		logPrintf("匹配点坐标:(%v,%v)", ptrPos.m_x, ptrPos.m_y)
	}
}

// 测试图片是否包含在大图中
func testContain() {
	baseF, e := os.Open(g_ptrConfig.Base)
	failOnError(e, "open base fail")
	defer baseF.Close()

	partF, e := os.Open(g_ptrConfig.PathPart + "/01.jpg") // 子图片路径
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
