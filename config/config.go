package config

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"

	"goSvrLib/log"
)

var (
	Global Config
)

func LoadConfigFile(filename string) error {
	var err error = nil
	Global.jsonObj, err = readFile(filename)

	if err != nil {
		log.Error("Read Json Error")
	}

	return err
}

func readFile(filename string) (j map[string]interface{}, e error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("config", "ReadFile", err.Error())
		return nil, err
	}

	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(bytes, &j); err != nil {
		log.Error("config", "Unmarshal: ", err.Error())
		return nil, err
	}

	return j, nil
}
