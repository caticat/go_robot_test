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
}

func newConfig(fileName string) *Config {
	ptrConfig := &Config{
		fileName: fileName,
	}

	sliFile, e := ioutil.ReadFile(fileName)
	failOnError(e, "open config failed")
	e = yaml.Unmarshal(sliFile, ptrConfig)
	failOnError(e, "config unmarshal failed")

	logPrintf("配置:%+v", *ptrConfig)

	return ptrConfig
}
