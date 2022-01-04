package bot

import (
	"github.com/bwmarrin/discordgo"
)

func ratio(arg []string, channelID string, messageID string, reply string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSend
	if len(arg) > 1 {
		_, err := send(channelID, m.Author.Username+", veuillez pas mettre d'arguements.")
		if err != nil {
			return err.Error(), ""
		}
	} else {
		_ = s.ChannelMessageDelete(channelID, messageID)
		ratio, err := s.ChannelMessageSendReply(channelID, "ratio", m.MessageReference)

		_ = s.MessageReactionAdd(m.ChannelID, ratio.ID, "👍")
		_ = s.MessageReactionAdd(m.ChannelID, ratio.ID, "👎")

		if err != nil {
			return err.Error(), ""
		}
		reply, err := s.ChannelMessage(channelID, m.MessageReference.MessageID)
		if err != nil {
			return err.Error(), ""
		}
		return "", reply.Content
	}

	return "", ""
}
