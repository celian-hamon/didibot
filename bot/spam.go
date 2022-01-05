package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func spam(arg []string, channelID string, m *discordgo.MessageCreate, s *discordgo.Session) (string,string) {
	send := s.ChannelMessageSend
	if len(arg) > 1 {
		for i := 0; i < 4; i++ {
			_, err := send(channelID, strings.Join(arg[1:], " "))
			if err != nil {
				return err.Error(),""
			}
		}
	} else {
		_, err := send(channelID, m.Author.Username+", veuillez mettre ce que vous voulez spam.")
		if err != nil {
			return err.Error(), ""
		}
	}
	return "",""
}
