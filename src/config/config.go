package config

// config 配置文件

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/json-iterator/go"
)

var (
	// GlobalConfig 全局配置文件，在Init调用前为nil
	GlobalConfig *Config
)

// Config 配置
type Config struct {
	Addr        string     `json:"addr"`
	Mongo       mongo      `json:"mongo"`
	Redis       redis      `json:"redis"`
	EnableDebug bool       `json:"debug"`
	Statistics  Statistics `json:"statistics"`
}

type mongo struct {
	Host string `json:"host"`
	Port string `json:"port"`
	DB   string `json:"db"`
}

type redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Statistics struct {
	RecordThreshold int `json:"record_threshold"`
	MinRecordCount  int `json:"min_record_count"`
	MaxRecordCount  int `json:"max_record_count"`
}

func init() {
	configFile := "default.json"

	// 如果有设置 ENV ，则使用ENV中的环境
	if v, ok := os.LookupEnv("ENV"); ok {
		configFile = v + ".json"
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(fmt.Sprintf("config/%s", configFile))

	if err != nil {
		log.Println("Read config error!")
		log.Panic(err)
		return
	}

	config := &Config{}

	err = jsoniter.Unmarshal(data, config)

	if err != nil {
		log.Println("Unmarshal config error!")
		log.Panic(err)
		return
	}

	GlobalConfig = config

	log.Println("Config " + configFile + " loaded.")
	// log.Printf("%+v", GlobalConfig)

}
