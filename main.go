package main

import (
	"ava-go/bot"
	"log"
)

func main() {
	if err := bot.Start(); err != nil {
		log.Fatalln("\u001b[31mERROR:\u001b[0m", err.Error())
	}

	<-make(chan struct{})
}
