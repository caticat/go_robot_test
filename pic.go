package main

import (
	"gocv.io/x/gocv"
	"image"
	"os"
	"path"
	"strconv"
	"strings"
)

// 全图
type Pic struct {
	m_fileName string
	m_pic      image.Image
}

func newPic() *Pic {
	return &Pic{}
}

func (this *Pic) init(fileName string) {
	this.m_fileName = path.Base(fileName)

	f, e := os.Open(fileName)
	failOnError(e, "open file %v failed", fileName)
	defer f.Close()

	i, _, e := image.Decode(f)
	failOnError(e, "decode file %v failed", fileName)
	this.m_pic = i
}

func (this *Pic) width() int {
	return this.m_pic.Bounds().Dx()
}

func (this *Pic) height() int {
	return this.m_pic.Bounds().Dy()
}

// 子图
type SubPic struct {
	*Pic
	m_value int
	m_mat   gocv.Mat
}

func newSubPic() *SubPic {
	return &SubPic{
		Pic:   newPic(),
		m_mat: gocv.NewMat(),
	}
}

func (this *SubPic) close() {
	this.m_mat.Close()
}

func (this *SubPic) init(dir, fileName string) {
	fullFileName := dir + "/" + fileName
	this.Pic.init(fullFileName)

	value, e := strconv.Atoi(strings.TrimRight(path.Base(fileName), path.Ext(fileName)))
	failOnError(e, "get file %v value failed", fileName)
	this.m_value = value

	hash := getHash()
	hashCompute(hash, this.m_pic, &this.m_mat)
}
