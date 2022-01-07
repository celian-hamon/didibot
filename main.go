package main

import (
	"discordbot/config"

	"discordbot/bot"
)

func main() {

	err := config.ReadConfig()
	config.Check(err)
	bot.Start()

	<-make(chan struct{})
	return
}
