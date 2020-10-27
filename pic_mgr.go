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

	for value, ptrSubPic := range this.m_mapPic {
		sliPos := getSliPosInPic(ptrPic, ptrSubPic)
		for _, ptrPos := range sliPos {
			logPrintf("值:%v,坐标:(%v,%v)", value, ptrPos.m_x, ptrPos.m_y)
		}
	}
}
