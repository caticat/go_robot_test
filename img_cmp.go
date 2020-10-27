package main

import (
	"bytes"
	"errors"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"time"
)

// 获取在图片中的所有位置
func getSliPosInPic(ptrPic *Pic, ptrSubPic *SubPic) []*Pos {
	// 参数校验
	assert(ptrPic != nil, "ptrPic == nil")
	assert(ptrSubPic != nil, "ptrSubPic == nil")

	// 参数整理
	picWidth := ptrPic.width()
	picHeight := ptrPic.height()
	picSubWidth := g_ptrConfig.BaseW
	picSubHeight := g_ptrConfig.BaseH
	if (picSubWidth == 0) || (picSubHeight == 0) {
		picSubWidth = ptrSubPic.width()
		picSubHeight = ptrSubPic.height()
	}

	sliPos := make([]*Pos, 0)
	hash := getHash()

	hasFound := false
	pic := ptrPic.m_pic.(*image.YCbCr)
	for y := 0; y <= (picHeight - picSubHeight); {
		hasFound = false
		for x := 0; x <= (picWidth - picSubWidth); {
			r := image.Rect(x, y, x+picSubWidth, y+picSubHeight)
			subBaseI := pic.SubImage(r)
			subBaseM := gocv.NewMat()
			defer subBaseM.Close()
			hashCompute(hash, subBaseI, &subBaseM)
			if isSame, _ := isSameMat(hash, &subBaseM, &ptrSubPic.m_mat); isSame {
				//logPrintf("匹配成功起始点:(%v, %v),相似度:%v", x, y, similar)
				ptrPos := newPos()
				ptrPos.init(x, y)
				sliPos = append(sliPos, ptrPos)
				hasFound = true
				x += picSubWidth
			} else {
				x += g_ptrConfig.StepX
			}
		}
		if hasFound {
			y += picSubHeight
		} else {
			y += g_ptrConfig.StepY
		}
	}

	return sliPos
}

// 获取子图片出现的数量 只横移一行
func getSubPicNumSlideX(ptrPic *Pic, ptrSubPic *SubPic, ptrBeginPos *Pos) int {
	// 参数校验
	assert(ptrPic != nil, "ptrPic == nil")
	assert(ptrSubPic != nil, "ptrSubPic == nil")
	assert(ptrBeginPos != nil, "ptrBeginPos == nil")

	// 参数整理
	picWidth := ptrPic.width()
	picSubWidth := g_ptrConfig.BaseW
	picSubHeight := g_ptrConfig.BaseH
	if (picSubWidth == 0) || (picSubHeight == 0) {
		picSubWidth = ptrSubPic.width()
		picSubHeight = ptrSubPic.height()
	}
	hash := getHash()
	pic := ptrPic.m_pic.(*image.YCbCr)
	y := ptrBeginPos.m_y

	counter := 0
	for x := ptrBeginPos.m_x; x < picWidth; x += picSubWidth {
		r := image.Rect(x, y, x+picSubWidth, y+picSubHeight)
		subBaseI := pic.SubImage(r)
		subBaseM := gocv.NewMat()
		defer subBaseM.Close()
		hashCompute(hash, subBaseI, &subBaseM)
		if isSame, _ := isSameMat(hash, &subBaseM, &ptrSubPic.m_mat); isSame {
			//logPrintf("相似度:%v", similar)
			counter++
		}
	}

	return counter
}

// 获取子图片出现的数量 只横移一行 单次匹配全部获取
func getSubPicNumSlideXAll(ptrPic *Pic, mapSubPic map[int]*SubPic, ptrBeginPos *Pos) map[int]int {
	// 参数校验
	assert(ptrPic != nil, "ptrPic == nil")
	assert(ptrBeginPos != nil, "ptrBeginPos == nil")
	assert(len(mapSubPic) > 0, "len(mapSubPic) == 0")
	if _, ok := mapSubPic[1]; !ok {
		failOnError(errors.New("mapSubPic do not contain 1"), "getSubPicNumSlideXAll")
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

// 获取图片中第一个匹配的起始点的精准坐标
func getFirstPos(ptrPic *Pic, ptrSubPic *SubPic) *Pos {
	// 参数校验
	assert(ptrPic != nil, "ptrPic == nil")
	assert(ptrSubPic != nil, "ptrSubPic == nil")

	// 参数整理
	picWidth := ptrPic.width()
	picHeight := ptrPic.height()
	picSubWidth := g_ptrConfig.BaseW
	picSubHeight := g_ptrConfig.BaseH
	if (picSubWidth == 0) || (picSubHeight == 0) {
		picSubWidth = ptrSubPic.width()
		picSubHeight = ptrSubPic.height()
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

// 图片是否包含
func contain(baseI image.Image, partI image.Image) bool {
	timeBegin := time.Now().UnixNano() / 1e6

	baseW := baseI.Bounds().Dx()
	baseH := baseI.Bounds().Dy()
	partW := g_ptrConfig.BaseW
	partH := g_ptrConfig.BaseH

	logPrintf("源图片宽:%v,高:%v;滑动窗口宽:%v,高:%v", baseW, baseH, partW, partH)

	hash := getHash()
	partM := gocv.NewMat()
	defer partM.Close()
	hashCompute(hash, partI, &partM)

	baseY := baseI.(*image.YCbCr)
	for y := 0; y <= (baseH - partH); y += g_ptrConfig.StepY {
		for x := 0; x <= (baseW - partW); x += g_ptrConfig.StepX {
			r := image.Rect(x, y, x+partW, y+partH)
			subBaseI := baseY.SubImage(r)
			subBaseM := gocv.NewMat()
			defer subBaseM.Close()
			hashCompute(hash, subBaseI, &subBaseM)
			if isSame, similar := isSameMat(hash, &subBaseM, &partM); isSame {
				timeEnd := time.Now().UnixNano() / 1e6
				logPrintf("匹配成功起始点:(%v, %v),相似度:%v,花费时间:%v(秒)", x, y, similar, float32(timeEnd-timeBegin)/1000.0)
				return true
			}
		}
	}

	timeEnd := time.Now().UnixNano() / 1e6
	logPrintf("图片不匹配,花费时间:%v(秒)", float32(timeEnd-timeBegin)/1000.0)

	return false
}

// 图片是否相同(相似)
func isSameImg(baseI image.Image, partI image.Image) bool {
	baseM, e := gocv.ImageToMatRGB(baseI)
	failOnError(e, "ImageToMatRGB")
	partM, e := gocv.ImageToMatRGB(partI)
	failOnError(e, "ImageToMatRGB")

	baseMT := gocv.NewMat()
	defer baseMT.Close()
	//hash := contrib.AverageHash{}
	hash := getHash()
	hash.Compute(baseM, &baseMT)
	if baseMT.Empty() {
		log.Fatalf("error computing hash for base")
	}

	partMT := gocv.NewMat()
	defer partMT.Close()
	hash.Compute(partM, &partMT)
	if partMT.Empty() {
		log.Fatalf("error computing hash for base")
	}

	similar := hash.Compare(baseMT, partMT)

	logPrintf("phash:%v", similar)
	return int(similar) <= g_ptrConfig.SimilarLimit
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

// 测试=================================

// 测试指定坐标是否相似
func testPosInPic(ptrPic *Pic, ptrSubPic *SubPic, ptrPos *Pos) []*Pos {
	// 参数校验
	assert(ptrPic != nil, "ptrPic == nil")
	assert(ptrSubPic != nil, "ptrSubPic == nil")
	assert(ptrPos != nil, "ptrPos == nil")

	// 参数整理
	picSubWidth := g_ptrConfig.BaseW
	picSubHeight := g_ptrConfig.BaseH
	if (picSubWidth == 0) || (picSubHeight == 0) {
		picSubWidth = ptrSubPic.width()
		picSubHeight = ptrSubPic.height()
	}

	sliPos := make([]*Pos, 0)
	hash := getHash()

	pic := ptrPic.m_pic.(*image.YCbCr)
	r := image.Rect(ptrPos.m_x, ptrPos.m_y, ptrPos.m_x+picSubWidth, ptrPos.m_y+picSubHeight)
	subBaseI := pic.SubImage(r)
	subBaseM := gocv.NewMat()
	defer subBaseM.Close()
	hashCompute(hash, subBaseI, &subBaseM)
	isSame, similar := isSameMat(hash, &subBaseM, &ptrSubPic.m_mat)
	logPrintf("起始点:(%v,%v),相似度:%v,是否相同:%v", ptrPos.m_x, ptrPos.m_y, similar, isSame)

	testSaveImg(subBaseI)

	return sliPos
}

func testSaveImg(i image.Image) {
	b := bytes.NewBuffer(nil)
	e := jpeg.Encode(b, i, nil)
	failOnError(e, "jpeg encode failed")
	ioutil.WriteFile("a.jpg", b.Bytes(), 0666)
}

//var (
//	useAll            = flag.Bool("all", false, "Compute all hashes")
//	usePHash          = flag.Bool("phash", false, "Compute PHash")
//	useAverage        = flag.Bool("average", false, "Compute AverageHash")
//	useBlockMean0     = flag.Bool("blockmean0", false, "Compute BlockMeanHash mode 0")
//	useBlockMean1     = flag.Bool("blockmean1", false, "Compute BlockMeanHash mode 1")
//	useColorMoment    = flag.Bool("colormoment", false, "Compute ColorMomentHash")
//	useMarrHildreth   = flag.Bool("marrhildreth", false, "Compute MarrHildrethHash")
//	useRadialVariance = flag.Bool("radialvariance", false, "Compute RadialVarianceHash")
//)
//
//func setupHashes() []contrib.ImgHashBase {
//	var hashes []contrib.ImgHashBase
//
//	if *usePHash || *useAll {
//		hashes = append(hashes, getHash())
//	}
//	if *useAverage || *useAll {
//		hashes = append(hashes, contrib.AverageHash{})
//	}
//	if *useBlockMean0 || *useAll {
//		hashes = append(hashes, contrib.BlockMeanHash{})
//	}
//	if *useBlockMean1 || *useAll {
//		hashes = append(hashes, contrib.BlockMeanHash{Mode: contrib.BlockMeanHashMode1})
//	}
//	if *useColorMoment || *useAll {
//		hashes = append(hashes, contrib.ColorMomentHash{})
//	}
//	if *useMarrHildreth || *useAll {
//		// MarrHildreth has default parameters for alpha/scale
//		hashes = append(hashes, contrib.NewMarrHildrethHash())
//	}
//	if *useRadialVariance || *useAll {
//		// RadialVariance has default parameters too
//		hashes = append(hashes, contrib.NewRadialVarianceHash())
//	}
//
//	// If no hashes were selected, behave as if all hashes were selected
//	if len(hashes) == 0 {
//		*useAll = true
//		return setupHashes()
//	}
//
//	return hashes
//}
//
//func main() {
//	flag.Usage = func() {
//		fmt.Println("How to run:\n\timg-similarity [-flags] [image1.jpg] [image2.jpg]")
//		flag.PrintDefaults()
//	}
//
//	printHashes := flag.Bool("print", false, "print hash values")
//	flag.Parse()
//	if flag.NArg() < 2 {
//		flag.Usage()
//		return
//	}
//
//	// read images
//	inputs := flag.Args()
//	images := make([]gocv.Mat, len(inputs))
//
//	for i := 0; i < 2; i++ {
//		img := gocv.IMRead(inputs[i], gocv.IMReadColor)
//		if img.Empty() {
//			fmt.Printf("cannot read image %s\n", inputs[i])
//			return
//		}
//		defer img.Close()
//
//		images[i] = img
//	}
//
//	// construct all of the hash types in a list. normally, you'd only use one of these.
//	hashes := setupHashes()
//
//	// compute and compare the images for each hash type
//	for _, hash := range hashes {
//		results := make([]gocv.Mat, len(images))
//
//		for i, img := range images {
//			results[i] = gocv.NewMat()
//			defer results[i].Close()
//			hash.Compute(img, &results[i])
//			if results[i].Empty() {
//				fmt.Printf("error computing hash for %s\n", inputs[i])
//				return
//			}
//		}
//
//		// compare for similarity; this returns a float64, but the meaning of values is
//		// unique to each algorithm.
//		similar := hash.Compare(results[0], results[1])
//
//		// make a pretty name for the hash
//		name := strings.TrimPrefix(fmt.Sprintf("%T", hash), "contrib.")
//		fmt.Printf("%s: similarity %g\n", name, similar)
//
//		if *printHashes {
//			// print hash result for each image
//			for i, path := range inputs {
//				fmt.Printf("\t%s = %x\n", path, results[i].ToBytes())
//			}
//		}
//	}
//}
