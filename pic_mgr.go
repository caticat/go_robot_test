package main

type PicMgr struct {
	m_mapPic map[int]*SubPic
}

func newPicMgr() *PicMgr {
	return &PicMgr{
		m_mapPic: make(map[int]*SubPic),
	}
}

func (this *PicMgr) close() {
	for _, pic := range this.m_mapPic {
		pic.close()
	}
}

func (this *PicMgr) load(dir string, sliFileName []string) {
	for _, fileName := range sliFileName {
		ptrSubPic := newSubPic()
		ptrSubPic.init(dir, fileName)
		this.m_mapPic[ptrSubPic.m_value] = ptrSubPic
	}
}

func (this *PicMgr) calc(baseFileName string) {
	ptrPic := newPic()
	ptrPic.init(baseFileName)

	// 获取到所有匹配的坐标
	//for value, ptrSubPic := range this.m_mapPic {
	//	sliPos := getSliPosInPic(ptrPic, ptrSubPic)
	//	for _, ptrPos := range sliPos {
	//		logPrintf("文件名:%v,值:%v,坐标:(%v,%v)", ptrSubPic.m_fileName, value, ptrPos.m_x, ptrPos.m_y)
	//	}
	//}

	// 获取所有值出现的数量
	this.calcCounter(ptrPic)
}

// 获取所有值出现的数量
func (this *PicMgr) calcCounter(ptrPic *Pic) {
	assert(ptrPic != nil, "ptrPic == nil")

	ptrBeginPos := this.getMatchBeginPos(ptrPic)
	if ptrBeginPos == nil {
		logPrintf("没有任何图片匹配")
		return
	}

	mapCounter := make(map[int]int)
	for _, ptrSubPic := range this.m_mapPic {
		counter := getSubPicNumSlideX(ptrPic, ptrSubPic, ptrBeginPos)
		mapCounter[ptrSubPic.m_value] = counter
	}

	logPrintf("匹配类型数量:%v", len(mapCounter))
	for value, counter := range mapCounter {
		logPrintf("值:%v,数量:%v", value, counter)
	}
}

// 获取所有子图片匹配的个数(x轴缩小,y轴不变)
func (this *PicMgr) getMatchBeginPos(ptrPic *Pic) *Pos {
	for _, ptrSubPic := range this.m_mapPic {
		if ptrPos := getFirstPos(ptrPic, ptrSubPic); ptrPos != nil {
			picSubWidth := g_ptrConfig.BaseW
			if picSubWidth == 0 {
				picSubWidth = ptrSubPic.width()
			}
			assert(picSubWidth > 0, "picSubWidth[%v] <= 0", picSubWidth)
			// 向左平移至最小值
			ptrPos.m_x = ptrPos.m_x - (ptrPos.m_x/picSubWidth)*picSubWidth
			return ptrPos
		}
	}
	return nil
}
