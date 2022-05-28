package config

import (
	"os"
)

var (
	Token     string
	BotPrefix string
	AdminList []string
	Invit     string
	config    *configStruct
)

type configStruct struct {
	Token     string   `json:"Token"`
	BotPrefix string   `json:"BotPrefix"`
	AdminList []string `json:"AdminList"`
	Invit     string   `json:"Invit"`
}

func ReadConfig() error {
	AdminList = []string{""}
	Invit = ""
	Token = os.Getenv("TOKEN")
	BotPrefix = os.Getenv("PREFIX")

	return nil
}
