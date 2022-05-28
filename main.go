package main

import (
	"discordbot/config"
	"net/http"
	"os"

	"discordbot/bot"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		return
	}
	go bot.Start()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)

	<-make(chan struct{})
}
