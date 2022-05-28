package main

import (
	"discordbot/config"
	"net/http"
	"os"

	"discordbot/bot"
)

func handler(w http.ResponseWriter, r *http.Request) {
	err := config.ReadConfig()
	if err != nil {
		return
	}
	go bot.Start()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)

	<-make(chan struct{})
}
