package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func sondage(arg []string, channelID string, m *discordgo.MessageCreate, s *discordgo.Session) string {
	send := s.ChannelMessageSend
	if len(arg) > 1 {
		msg, err := send(m.ChannelID, "Le sondage de "+m.Author.Username+" est : "+strings.Join(arg[1:], " "))
		_ = s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		_ = s.MessageReactionAdd(m.ChannelID, msg.ID, "👍")
		_ = s.MessageReactionAdd(m.ChannelID, msg.ID, "👎")
		if err != nil {
			return err.Error()
		}
	} else {
		_, err := send(m.ChannelID, m.Author.Username+", veuillez mettre l'intitulé du sondage.")
		if err != nil {
			return err.Error()
		}
	}
	return ""
}
