package bot

import (
	"github.com/bwmarrin/discordgo"
)

func help(arg []string, channelID string, userID string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSend
	if len(arg) > 1 {
		_, err := send(channelID, m.Author.Username+", veuillez pas mettre d'arguements.")
		if err != nil {
			return err.Error(),""
		}
	} else {
		prv, err := s.UserChannelCreate(userID)
		if err != nil {
			return err.Error(),""
		}
		_, err = send(channelID, "Go check tes mp batard.")
		if err != nil {
			return err.Error(),""
		}
		_, err = s.ChannelMessageSend(prv.ID, "Help en cours de dev, ça arrive chakal !")
		if err != nil {
			return err.Error(),""
		}

	}
	return "",""
}
