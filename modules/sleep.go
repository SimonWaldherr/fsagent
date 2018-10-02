package modules

import (
	"encoding/json"
	"fmt"
	"simonwaldherr.de/go/golibs/file"
	"time"
)

type sleepConfig struct {
	Time int `json:"time"`
}

func Sleep(configName, fileName string) error {
	_ = fileName
	fmt.Printf("file name: %v\n", fileName)
	var c sleepConfig
	fmt.Printf("load config from: %v\n", configName)
	str, _ := file.Read(configName)
	json.Unmarshal([]byte(str), &c)
	time.Sleep(time.Millisecond * time.Duration(c.Time))
	return nil
}
