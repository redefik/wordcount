// Author: Federico Viglietta
package config

import (
	"encoding/json"
	"errors"
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

// Returns the struct encapsulating configuration information reading it from the file given.
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
	if len(config.Master) == 0 {
		return config, errors.New("no master found")
	}
	if len(config.Reducer) == 0 {
		return config, errors.New("no reducer found")
	}
	if len(config.Mapper) == 0 {
		return config, errors.New("no mapper found")
	}
	return config, nil
}
