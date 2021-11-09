package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
)

var MaxTimes int64
var Probability float64
var TotalMoney int64
var TotalNum int64
var MaxMoney int64
var MinMoney int64

type Configs map[string]json.RawMessage

type MainConfig struct {
	Port string `json:"port"`
	Address string `json:"address"`
}

var conf *MainConfig
var confs Configs

var instanceOnce sync.Once

//从配置文件中载入json字符串
func LoadConfig(path string) (Configs, *MainConfig) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("load config conf failed: ", err)
	}
	mainConfig := &MainConfig{}
	err = json.Unmarshal(buf, mainConfig)
	if err != nil {
		log.Panicln("decode config file failed:", string(buf), err)
	}
	allConfigs := make(Configs, 0)
	err = json.Unmarshal(buf, &allConfigs)
	if err != nil {
		log.Panicln("decode config file failed:", string(buf), err)
	}
	return allConfigs, mainConfig
}

// 初始化，只能运行一次
func Init(path string) *MainConfig {

	instanceOnce.Do(func() {
		allConfigs, mainConfig := LoadConfig(path)
		conf = mainConfig
		confs = allConfigs
	})

	return conf
}

func InitConfigs(path string) Configs {

	Init(path)
	MaxTimes, _ = strconv.ParseInt(string(confs["maxTimes"]), 10, 64)
	Probability, _ = strconv.ParseFloat(string(confs["probability"]), 64)
	TotalMoney, _ = strconv.ParseInt(string(confs["totalMoney"]), 10, 64)
	TotalNum,_ = strconv.ParseInt(string(confs["totalNum"]), 10, 64)
	MaxMoney, _ = strconv.ParseInt(string(confs["maxMoney"]), 10, 64)
	MinMoney,_ = strconv.ParseInt(string(confs["minMoney"]), 10, 64)
	return confs
}

