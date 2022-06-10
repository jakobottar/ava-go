package main

import (
	"ava-go/bot"
	"ava-go/config"
	"runtime"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

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
	err := config.ReadConfig()

	if err != nil {
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
