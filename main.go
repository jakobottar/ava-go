package main

import (
	"ava-go/bot"
)

func main() {
	bot.Start()

	<-make(chan struct{})
}
