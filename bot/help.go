package bot

import (
	"discordbot/config"

	"github.com/bwmarrin/discordgo"
)

func help(arg []string, channelID string, userID string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSend
	if len(arg) > 1 {
		_, err := send(channelID, m.Author.Username+", veuillez pas mettre d'arguements.")
		config.Check(err)
	} else {
		prv, err := s.UserChannelCreate(userID)
		config.Check(err)
		_, err = send(channelID, "Go check tes mp.")
		config.Check(err)
		_, err = s.ChannelMessageSend(prv.ID, "Help en cours de dev, ça arrive chakal !")
		config.Check(err)

	}
	return "", ""
}
