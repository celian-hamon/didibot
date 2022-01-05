package main

import (
	"discordbot/config"

	"discordbot/bot"
)

func main() {

	err := config.ReadConfig()
	if err != nil {
		return
	}
	bot.Start()

	<-make(chan struct{})
	return
}
