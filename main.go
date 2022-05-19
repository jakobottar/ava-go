package main

import (
	"ava-go/bot"
	"ava-go/config"
	"log"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		log.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
