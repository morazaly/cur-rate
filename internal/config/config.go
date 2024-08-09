package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	MySQL struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
	} `json:"mysql"`
	AppPort string `json:"appPort"`
	ApiURL  string `json:"apiURL"`
}

func NewConfig() *Config {
	var aconfig Config

	configFile, err := os.Open("..\\internal\\config\\config.json") //
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	configDecoder := json.NewDecoder(configFile)

	if err := configDecoder.Decode(&aconfig); err != nil {
		log.Fatal(err)
	}
	return &aconfig
}
