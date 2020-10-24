package main

import (
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
	"image"
	"log"
	"time"
)

// 图片是否包含
func contain(baseI image.Image, partI image.Image) bool {
	timeBegin := time.Now().UnixNano() / 1e6

	baseW := baseI.Bounds().Dx()
	baseH := baseI.Bounds().Dy()
	partW := g_ptrConfig.BaseW
	partH := g_ptrConfig.BaseH

	logPrintf("源图片宽:%v,高:%v;滑动窗口宽:%v,高:%v", baseW, baseH, partW, partH)

	hash := contrib.PHash{}
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
	hash := contrib.PHash{}
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
	return similar <= 8 // TODO(pan): 这里需要改成配置
}

func isSameMat(hash contrib.ImgHashBase, ptrBaseM *gocv.Mat, ptrPartM *gocv.Mat) (bool, float64) {

	similar := hash.Compare(*ptrBaseM, *ptrPartM)

	//logPrintf("hash similar:", similar)
	return similar <= float64(g_ptrConfig.SimilarLimit), similar // TODO(pan): 这里需要改成配置
}

func hashCompute(hash contrib.ImgHashBase, img image.Image, ptrOut *gocv.Mat) {
	imgM, e := gocv.ImageToMatRGB(img)
	failOnError(e, "ImageToMatRGB")

	hash.Compute(imgM, ptrOut)
	if ptrOut.Empty() {
		log.Fatalf("error computing hash for base")
	}
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
//		hashes = append(hashes, contrib.PHash{})
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
