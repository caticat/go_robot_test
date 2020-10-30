package main

import (
	"errors"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
	"image"
	"log"
)

// 获取子图片出现的数量 只横移一行 单次匹配全部获取
func getSubPicNumSlideX(ptrPic *Pic, mapSubPic map[int]*SubPic, ptrBeginPos *Pos) map[int]int {
	// 参数校验
	assert(ptrPic != nil, "ptrPic == nil")
	assert(ptrBeginPos != nil, "ptrBeginPos == nil")
	assert(len(mapSubPic) > 0, "len(mapSubPic) == 0")
	if _, ok := mapSubPic[1]; !ok {
		failOnError(errors.New("mapSubPic do not testContain 1"), "getSubPicNumSlideX")
	}

	// 参数整理
	picWidth := ptrPic.width()
	picSubWidth := g_ptrConfig.BaseW
	picSubHeight := g_ptrConfig.BaseH
	if (picSubWidth == 0) || (picSubHeight == 0) {
		picSubWidth = mapSubPic[1].width()
		picSubHeight = mapSubPic[1].height()
	}
	hash := getHash()
	pic := ptrPic.m_pic.(*image.YCbCr)
	y := ptrBeginPos.m_y
	mapCounter := make(map[int]int)

	// 计算
	for x := ptrBeginPos.m_x; x < picWidth; x += picSubWidth {
		r := image.Rect(x, y, x+picSubWidth, y+picSubHeight)
		subBaseI := pic.SubImage(r)
		subBaseM := gocv.NewMat()
		defer subBaseM.Close()
		hashCompute(hash, subBaseI, &subBaseM)
		for _, ptrSubPic := range mapSubPic {
			if isSame, _ := isSameMat(hash, &subBaseM, &ptrSubPic.m_mat); isSame {
				//logPrintf("相似度:%v", similar)
				mapCounter[ptrSubPic.m_value]++
				break
			}
		}
	}

	return mapCounter
}

// 获取图片中第一个匹配的起始点的精准坐标 单次匹配全部获取
func getFirstPos(ptrPic *Pic, mapSubPic map[int]*SubPic) *Pos {
	// 参数校验
	assert(ptrPic != nil, "ptrPic == nil")
	assert(len(mapSubPic) > 0, "len(mapSubPic) == 0")
	if _, ok := mapSubPic[1]; !ok {
		failOnError(errors.New("mapSubPic do not testContain 1"), "getSubPicNumSlideX")
	}

	// 参数整理
	picWidth := ptrPic.width()
	picHeight := ptrPic.height()
	picSubWidth := g_ptrConfig.BaseW
	picSubHeight := g_ptrConfig.BaseH
	if (picSubWidth == 0) || (picSubHeight == 0) {
		picSubWidth = mapSubPic[1].width()
		picSubHeight = mapSubPic[1].height()
	}

	hash := getHash()
	pic := ptrPic.m_pic.(*image.YCbCr)
	for y := 0; y <= (picHeight - picSubHeight); y += g_ptrConfig.StepY {
		for x := 0; x <= (picWidth - picSubWidth); x += g_ptrConfig.StepX {
			r := image.Rect(x, y, x+picSubWidth, y+picSubHeight)
			subBaseI := pic.SubImage(r)
			subBaseM := gocv.NewMat()
			defer subBaseM.Close()
			hashCompute(hash, subBaseI, &subBaseM)
			for _, ptrSubPic := range mapSubPic {
				if isSame, similar := isSameMat(hash, &subBaseM, &ptrSubPic.m_mat); isSame {
					beginPos := newPos()
					beginPos.init(x, y)
					ptrResultPos, similarV := getMostSimilarPos(ptrPic, ptrSubPic, beginPos, picWidth, picHeight, picSubWidth, picSubHeight)
					assert(ptrResultPos != nil, "ptrResultPos == nil")
					logPrintf("值:%v,匹配成功起始点:(%v, %v),相似度:%v,最相似点:(%v, %v),相似度:%v", ptrSubPic.m_value, x, y, similar, ptrResultPos.m_x, ptrResultPos.m_y, similarV)
					return ptrResultPos
				}
			}
		}
	}

	return nil
}

func getMostSimilarPos(ptrPic *Pic, ptrSubPic *SubPic, ptrBeginPos *Pos, picWidth, picHeight, picSubWidth, picSubHeight int) (*Pos, int) {
	// 参数校验
	assert(ptrPic != nil, "ptrPic == nil")
	assert(ptrSubPic != nil, "ptrSubPic == nil")
	assert(ptrBeginPos != nil, "ptrBeginPos == nil")

	beginX := ptrBeginPos.m_x - g_ptrConfig.BlurryPixel
	if beginX < 0 {
		beginX = 0
	}
	beginY := ptrBeginPos.m_y
	endX := ptrBeginPos.m_x + g_ptrConfig.BlurryPixel
	if endX > picWidth {
		endX = picWidth
	}
	endY := ptrBeginPos.m_y + g_ptrConfig.BlurryPixel
	if endY > picHeight {
		endY = picHeight
	}

	hash := getHash()
	pic := ptrPic.m_pic.(*image.YCbCr)
	valuePos := newValuePos()
	valuePos.m_value = 0xEFFFFFFF
	found := false
	for y := beginY; y < endY; y++ {
		for x := beginX; x < endX; x++ {
			r := image.Rect(x, y, x+picSubWidth, y+picSubHeight)
			subBaseI := pic.SubImage(r)
			subBaseM := gocv.NewMat()
			defer subBaseM.Close()
			hashCompute(hash, subBaseI, &subBaseM)
			if isSame, similar := isSameMat(hash, &subBaseM, &ptrSubPic.m_mat); isSame {
				//logPrintf("坐标:(%v,%v),相似度:%v", x, y, similar)
				if valuePos.m_value > int(similar) {
					valuePos.m_value = int(similar)
					valuePos.m_ptrPos.m_x = x
					valuePos.m_ptrPos.m_y = y
					found = true
				}
			}
		}
	}

	if found {
		return valuePos.m_ptrPos, valuePos.m_value
	} else {
		return nil, -1
	}
}

func getHash() contrib.ImgHashBase {
	return contrib.PHash{}
}

func isSameMat(hash contrib.ImgHashBase, ptrBaseM *gocv.Mat, ptrPartM *gocv.Mat) (bool, float64) {

	similar := hash.Compare(*ptrBaseM, *ptrPartM)

	//logPrintf("hash similar:", similar)
	return similar <= float64(g_ptrConfig.SimilarLimit), similar
}

func hashCompute(hash contrib.ImgHashBase, img image.Image, ptrOut *gocv.Mat) {
	imgM, e := gocv.ImageToMatRGB(img)
	failOnError(e, "ImageToMatRGB")

	hash.Compute(imgM, ptrOut)
	if ptrOut.Empty() {
		log.Fatalf("error computing hash for base")
	}
}
