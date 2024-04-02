package config

import (
	"encoding/json"
	"fmt"
	"launchpad/logger"
	"os"
)

var Configure Configuration

type Configuration struct {
	Log        logger.Conf `json:"log"`
	PostgreCfg Postgre     `json:"postgre_cfg"`
	Web        Web         `json:"web"`
}

type Postgre struct {
	Conf map[string]string `json:"conf"`
}

type Web struct {
	Address string `json:"address"`
}

func InitConfig(path string) {
	if path == "" {
		path = `./conf.json`
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("fialed to load config file.[err=%v]\n", err)
		panic("load config file failed ")
	}

	stat, _ := file.Stat()

	bytes := make([]byte, stat.Size())
	Configure.PostgreCfg.Conf = make(map[string]string)

	_, err = file.Read(bytes)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &Configure)

	if err != nil {
		fmt.Println("failed to unmarshal config json")
		panic("wrong config json")
	}
	fmt.Println("init config successed")
}
