package main

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type Config struct {
	fileName     string
	Base         string
	PathPart     string `yaml:"pathPart"`
	PartExt      string `yaml:"partExt"`
	StepX        int    `yaml:"stepX"`
	StepY        int    `yaml:"stepY"`
	BaseW        int    `yaml:"baseW"`
	BaseH        int    `yaml:"baseH"`
	SimilarLimit int    `yaml:"similarLimit"`
	BlurryPixel  int    `yaml:"blurryPixel"`
}

func newConfig(fileName string) *Config {
	ptrConfig := &Config{
		fileName: fileName,
	}

	sliFile, e := ioutil.ReadFile(fileName)
	failOnError(e, "open config failed")
	e = yaml.Unmarshal(sliFile, ptrConfig)
	failOnError(e, "config unmarshal failed")

	ptrConfig.init()

	logPrintf("配置:%+v", *ptrConfig)

	return ptrConfig
}

func (this *Config) init() {
	maxStep := this.StepX
	if maxStep < this.StepY {
		maxStep = this.StepY
	}
	if this.BlurryPixel < maxStep {
		this.BlurryPixel = maxStep
	}
}
