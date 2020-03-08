package yamlReader

import (
	"codezexcel/CodeZExcelTool/tool/yaml/model"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

type YamlReader struct {
	Configure yamlConfig.YamlConfigure
}

func Instance() (*YamlReader, error) {
	var filePath, err = filepath.Abs(yamlConfig.YAMLDEFAULTPATH)
	if err != nil {
		return nil, err
	}

	log.Printf("正在获取yaml文件路径[%v] \n", filePath)
	yr := new(YamlReader)
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panicf("yaml文件[%v]获取错误：%v\n", filePath, err.Error())
		return nil, err
	}

	tagConfigure := new(yamlConfig.YamlConfigure)

	err = yaml.Unmarshal(yamlFile, tagConfigure)

	if err != nil {
		log.Panicf("读取数据错误[%v]", err)
		return nil, err
	}

	yr.Configure = *tagConfigure

	return yr, nil
}