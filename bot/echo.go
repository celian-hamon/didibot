package bot

import (
	"discordbot/config"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func echo(arg []string, channelID string, messageID string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSend
	if len(arg) > 1 {
		_, err := send(m.ChannelID, strings.Join(arg[1:], " "))
		config.Check(err)
		err = s.ChannelMessageDelete(channelID, messageID)
		config.Check(err)
	} else {
		_, err := send(m.ChannelID, m.Author.Username+", veuillez mettre un argument.")
		config.Check(err)
	}
	return "", ""
}
