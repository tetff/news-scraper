package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	APIKey string `json:"apiKey"`
	Conn   string `json:"connection"`
}

func ReadJson(configFile string) (Config, error) {
	var config Config
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return config, err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return config, err
	}
	return config, err
}
