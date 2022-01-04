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
	goBot.AddHandler(ready)
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
			err, reply := timer(arg, m.ChannelID, m, s)
			security.Log(cmd, arg, err, reply, m)
		}
		if cmd == "sondage" {
			err, reply := sondage(arg, m.ChannelID, m, s)
			security.Log(cmd, arg, err, reply, m)
		}
		if cmd == "echo" {
			err, reply := echo(arg, m.ChannelID, m.Message.ID, m, s)
			security.Log(cmd, arg, err, reply, m)
		}
		if cmd == "spam" {
			if !security.IsAdmin(m.Author.ID, config.AdminList) {
				_, err := send(m.ChannelID, "Tu as pas les droits chacal")
				if err != nil {
					security.Log(cmd, arg, err.Error(), "", m)
				}
			} else {
				err, reply := spam(arg, m.ChannelID, m, s)
				security.Log(cmd, arg[1:], err, reply, m)
			}

		}
		if cmd == "help" {
			err, reply := help(arg, m.ChannelID, m.Message.Author.ID, m, s)
			security.Log(cmd, arg, err, reply, m)
		}
		if cmd == "ratio" {
			err, reply := ratio(arg, m.ChannelID, m.Message.ID, m.MessageReference.MessageID, m, s)
			security.Log(cmd, arg, err, reply, m)
		}
	}

}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	_ = s.UpdateGameStatus(0, "sucer des musulmans")
}
