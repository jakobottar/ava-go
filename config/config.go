package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	Token     string
	BotPrefix string
	config    *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	log.Println("config: reading config file...")
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalln("\u001b[31mERROR:\u001b[0m", err.Error())
		return err
	}

	// log.Println(string(file))
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln("\u001b[31mERROR:\u001b[0m", err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}
