package bot

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func timer(arg []string, channelID string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSend
	if len(arg) > 1 {
		_, err := send(m.ChannelID, m.Author.Username+", veuillez pas mettre d'arguments.")
		if err != nil {
			return err.Error(), ""
		}
	} else {
		currentTime := time.Now()
		_, err := send(m.ChannelID, "Time : "+currentTime.Format("02-01-2006 15:04:05"))
		if err != nil {
			return err.Error(), ""
		}
	}
	return "", ""
}
