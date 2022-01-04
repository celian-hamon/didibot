package bot

import (
	"fmt"
	"strings"

	"discordbot/config"
	"discordbot/security"

	"github.com/bwmarrin/discordgo"
)

var BotID string

func Start() {

	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())

	}

	BotID = u.ID
	goBot.AddHandler(messageHandler)
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Bot is running")

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	send := s.ChannelMessageSend
	if strings.HasPrefix(m.Content, config.BotPrefix) {
		m.Content = strings.Replace(m.Content, config.BotPrefix, "", -1)
		arg := strings.Split(m.Content, " ")
		cmd := arg[0]
		if cmd == "time" {
			err := timer(arg, m.ChannelID, m, s)
			security.Log(cmd, arg, err, m)
		}
		if cmd == "sondage" {
			err := sondage(arg, m.ChannelID, m, s)
			security.Log(cmd, arg, err, m)
		}
		if cmd == "echo" {
			err := echo(arg, m.ChannelID, m.Message.ID, m, s)
			security.Log(cmd, arg, err, m)
		}
		if cmd == "spam" {
			if !security.IsAdmin(m.Author.ID, config.AdminList) {
				_, err := send(m.ChannelID, "Tu as pas les droits chacal")
				if err != nil {
					security.Log(cmd, arg, err.Error(), m)
				}
			} else {
				err := spam(arg, m.ChannelID, m, s)
				security.Log(cmd, arg[1:], err, m)
			}

		}
		if cmd == "help" {
			err := help(arg, m.ChannelID, m.Message.Author.ID, m, s)
			security.Log(cmd, arg, err, m)
		}

	}

}
