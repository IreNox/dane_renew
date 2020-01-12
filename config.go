package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"os"
)

type config struct {
	URL     	string   `json:"url"`
	AuthToken   string   `json:"authToken"`
	Domains 	[]string `json:"domains"`
}

func configPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, err := filepath.Abs(filepath.Dir(exe))
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.json"), nil
}

func readConfig(config *config) error {
	configPath, err := configPath()
	if err != nil {
		return err
	}

	fmt.Printf("Load config from: %s\n", configPath)

	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}

	defer configFile.Close()

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(configBytes, config)
}
