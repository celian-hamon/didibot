package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	AdminList = []string{}
	fmt.Println("Reading from config file...")
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		Token := os.Getenv("TOKEN")
		BotPrefix = os.Getenv("BOTPREFIX")
		Invit = ""
		AdminList = []string{}
		if Token == "" || BotPrefix == "" {
			fmt.Println("Error reading from env :", err)
			return err
		}
	}
	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	Token = config.Token
	BotPrefix = config.BotPrefix
	AdminList = config.AdminList
	Invit = config.Invit

	return nil
}
