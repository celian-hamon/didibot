package main

import (
	"discordbot/config"
	"net/http"
	"os"

	"discordbot/bot"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go http.ListenAndServe(":"+port, nil)

	err := config.ReadConfig()
	if err != nil {
		return
	}
	bot.Start()

	<-make(chan struct{})
}
