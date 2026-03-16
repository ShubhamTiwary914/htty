package htty_test

import (
	"encoding/json"
	"fmt"
	types "htty/types"
	utils "htty/utils"
	"os"
	"testing"
)

//check config file bytes load
func TestLoadingConfigFile(tt *testing.T) {
	tt.Run("default config path", func(tt *testing.T) {
		os.Unsetenv("CONFIG_FILE")
		_, err := utils.GetConfigFile()
		if err != nil {
			tt.Fatalf("unexpected error: %v", err)
		}
	})

	tt.Run("absolute CONFIG_FILE", func(tt *testing.T) {
		tmp, err := os.CreateTemp("", "htty-config-*.json")
		if err != nil {
			tt.Fatalf("temp file creation failed: %v", err)
		}
		defer os.Remove(tmp.Name())
		expected := []byte("{'k': 'v'}")
		err = os.WriteFile(tmp.Name(), expected, 0644)
		if err != nil {
			tt.Fatalf("write failed: %v", err)
		}
		os.Setenv("CONFIG_FILE", tmp.Name())
		defer os.Unsetenv("CONFIG_FILE")
		data, err := utils.GetConfigFile()
		if err != nil {
			tt.Fatalf("unexpected error: %v", err)
		}
		if string(data) != string(expected) {
			tt.Fatalf("unexpected config contents")
		}
	})

	tt.Run("relative CONFIG_FILE should fail", func(tt *testing.T) {
		os.Setenv("CONFIG_FILE", "./config.conf")
		defer os.Unsetenv("CONFIG_FILE")
		_, err := utils.GetConfigFile()
		if err == nil {
			tt.Fatalf("expected error for relative CONFIG_FILE")
		}
	})
}


func TestConfigParsing(tt *testing.T){
	content, err := utils.GetConfigFile()
	if err != nil {
        tt.Fatal("Error when opening file: ", err)
    }
	var payload types.HttyConfig
	err = json.Unmarshal(content, &payload);
	if err != nil {
        tt.Fatal("Error during Unmarshal(): ", err)
    }
	fmt.Println("PAYLOAD:")
}
