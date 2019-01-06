package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Encapsulates the fields of the configuration file
type Config struct {
	Master []string
	Mapper []string
	Reducer []string
	OutDir string
}

func GetConfiguration(configFile string) (Config, error) {
	var config Config
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return config, err
	}
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}