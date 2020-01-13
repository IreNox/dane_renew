package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"os"
)

type configDomain struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

type config struct {
	URL     	string         `json:"url"`
	AuthToken   string         `json:"authToken"`
	Domains     []configDomain `json:"domains"`
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
