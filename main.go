package main

import (
	"ava-go/bot"
	"ava-go/config"
	"log"
	"os"
	"runtime"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// func main() {

// 	if err != nil {

// 	}
// 	defer f.Close()
// 	log.SetOutput(f)

// 	err = config.ReadConfig()
// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}

// 	bot.Start()

// 	<-make(chan struct{})
// 	return
// }

func main() {
	ava := NewAva()

	err := ava.Start()

	if err != nil {
		panic(err)
	}

	defer ava.Close()
	runtime.Goexit()
}

type Ava struct {
	session *discordgo.Session
}

func NewAva() Ava {
	return Ava{
		session: nil,
	}
}

func (a *Ava) Start() error {
	f, err := os.OpenFile("bot.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return errors.Wrapf(err, "Failed to create log file")
	}

	defer f.Close()
	log.SetOutput(f)

	if err = config.ReadConfig(); err != nil {
		return errors.Wrapf(err, "Failed to read ava config")
	}

	if a.session, err = bot.Start(); err != nil {
		return errors.Wrapf(err, "Bot failure")
	}

	return nil
}

func (a *Ava) Close() (err error) {
	if err = a.session.Close(); err != nil {
		return errors.Wrapf(err, "Failed to close session")
	}

	return
}
