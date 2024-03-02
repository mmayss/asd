package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func GetConfig() config {
	config := config{}
	// read file from config file for user
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	CONFIG_PATH := configDir + "/asd/config.json"

	defaultConfig, _ := os.Open("./config.json")
	configFile, err := os.Open(CONFIG_PATH)

	if err != nil {
		fmt.Println("config file not found, creating one at: ", CONFIG_PATH)
		os.MkdirAll(configDir+"/asd", 0755)
		configFile, err = os.Create(configDir + "/asd/config.json")
		defaultConfigContent, _ := io.ReadAll(defaultConfig)
		configFile.Write(defaultConfigContent)
		panic("")
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}

type config struct {
	Repos struct {
		Useful  []Repo `json:"useful"`
		Useless []Repo `json:"useless"`
	} `json:"repos"`

	Token string `json:"token"`
}

type Repo struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}
