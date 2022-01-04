package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func echo(arg []string, channelID string, messageID string, m *discordgo.MessageCreate, s *discordgo.Session) string {
	send := s.ChannelMessageSend
	if len(arg) > 1 {
		_, err := send(m.ChannelID, strings.Join(arg[1:], " "))
		if err != nil {
			return err.Error()
		}
	} else {
		_, err := send(m.ChannelID, m.Author.Username+", veuillez mettre un argument.")
		if err != nil {
			return err.Error()
		} 
	}
	return ""
}
