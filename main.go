package main

import (
	"ava-go/bot"
	"ava-go/config"
	"log"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		log.Fatalln("\u001b[31mERROR:\u001b[0m", err.Error())
	}

	bot.Start()

	<-make(chan struct{})
}
