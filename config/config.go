package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Check(err)

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)
	Check(err)
	Token = config.Token
	BotPrefix = config.BotPrefix
	AdminList = config.AdminList
	Invit = config.Invit

	return nil
}
func Check(err error) (string, string) {
	if err != nil {
		fmt.Println(err.Error())
		return err.Error(), ""
	}
	return "", ""
}
