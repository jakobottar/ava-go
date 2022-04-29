package main

import (
	"ava-go/bot"
	"ava-go/config"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("bot.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	err = config.ReadConfig()
	if err != nil {
		log.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
