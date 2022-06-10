package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
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
		return errors.Wrapf(err, "Failed to open config file")
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse JSON format")
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}
