/*
dealing with config management stuff for loading config
*/
package htty

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	types "htty/types"
)

func GetConfigFile() ([]byte, error) {
	//if path is set -> make sure its absolute
	configpath := os.Getenv("CONFIG_FILE")
	if configpath != "" {
		if !filepath.IsAbs(configpath) {
			return nil, errors.New("CONFIG_FILE must be an absolute path")
		}
		return os.ReadFile(configpath)
	}
	cfgHome := os.Getenv("XDG_CONFIG_HOME")
	//if path not set, go to default path: 
	if cfgHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		cfgHome = filepath.Join(home, ".config")
	}
	configpath = filepath.Join(cfgHome, "htty", "config.json")
	return os.ReadFile(configpath)
}

func GetConfig() (types.HttyConfig, error) {
	contentFile, err := GetConfigFile()	
	if err != nil {
		return types.HttyConfig{}, err
	}
	var configJson types.HttyConfig
	err = json.Unmarshal(contentFile, &configJson);
	if err != nil {
		return types.HttyConfig{}, err
	}
   	return configJson, nil 
}
