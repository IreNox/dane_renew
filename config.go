package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type config struct {
	URL     string   `json:"url"`
	APIKey  string   `json:"apiKey"`
	Domains []string `json:"domains"`
}

func readConfig(configPath string, config *config) error {
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer configFile.Close()

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return json.Unmarshal(configBytes, config)
}
