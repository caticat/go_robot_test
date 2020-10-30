package main

import "errors"

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

	// 获取所有值出现的数量
	this.calcCounterAll(ptrPic)
}

// 获取所有值出现的数量 单次匹配全部获取
func (this *PicMgr) calcCounterAll(ptrPic *Pic) {
	assert(ptrPic != nil, "ptrPic == nil")

	ptrBeginPos := this.getMatchBeginPosAll(ptrPic)
	if ptrBeginPos == nil {
		logPrintf("没有任何图片匹配")
		return
	}

	mapCounter := getSubPicNumSlideX(ptrPic, this.m_mapPic, ptrBeginPos)

	logPrintf("匹配类型数量:%v", len(mapCounter))
	for value, counter := range mapCounter {
		logPrintf("值:%v,数量:%v", value, counter)
	}
}

// 获取所有子图片匹配的个数(x轴缩小,y轴不变) 单次匹配全部获取
func (this *PicMgr) getMatchBeginPosAll(ptrPic *Pic) *Pos {
	assert(len(this.m_mapPic) > 0, "len(this.m_mapPic) == 0")
	if _, ok := this.m_mapPic[1]; !ok {
		failOnError(errors.New("this.m_mapPic do not testContain 1"), "getMatchBeginPosAll")
	}

	if ptrPos := getFirstPos(ptrPic, this.m_mapPic); ptrPos != nil {
		picSubWidth := g_ptrConfig.BaseW
		if picSubWidth == 0 {
			picSubWidth = this.m_mapPic[1].width()
		}
		assert(picSubWidth > 0, "picSubWidth[%v] <= 0", picSubWidth)
		// 向左平移至最小值
		ptrPos.m_x = ptrPos.m_x - (ptrPos.m_x/picSubWidth)*picSubWidth
		return ptrPos
	} else {
		return nil
	}
}
