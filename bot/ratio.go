package bot

import (
	"discordbot/config"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ratio(arg []string, channelID string, messageID string, reply string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {

	send := s.ChannelMessageSend
	if len(arg) > 1 {
		_, err := send(channelID, m.Author.Username+", veuillez pas mettre d'arguements.")
		config.Check(err)

	} else {
		_ = s.ChannelMessageDelete(channelID, messageID)
		ratio, err := s.ChannelMessageSendReply(channelID, "ratio", m.MessageReference)
		_ = s.MessageReactionAdd(m.ChannelID, ratio.ID, "👍")
		_ = s.MessageReactionAdd(m.ChannelID, ratio.ID, "👎")
		config.Check(err)
		reply, err := s.ChannelMessage(channelID, m.MessageReference.MessageID)
		config.Check(err)

		timer := 1 * time.Minute
		t := time.Now().Add(timer)
		for {
			if time.Now().Before(t) {
				test, err := s.MessageReactions(channelID, ratio.ID, "👍", 99, "", "1")
				config.Check(err)
				fmt.Println(len(test))
				if len(test) > 4 {
					_, err := s.ChannelMessageSendReply(channelID, "Tu t'es fais ratio boloss", m.MessageReference)
					config.Check(err)
					prv, err := s.UserChannelCreate(reply.Author.ID)
					config.Check(err)
					_, err = send(prv.ID, config.Invit)
					config.Check(err)
					err = s.GuildMemberDeleteWithReason(m.GuildID, reply.Author.ID, "tu t'es fais ratio mdr")
					config.Check(err)
					break
				}
				continue
			}
			_, err = send(channelID, "fin du ratio t'es trop nul <@!"+m.Author.ID+">")
			if err != nil {
				return err.Error(), reply.Content
			}
			break
		}
		return "", reply.Content
	}
	return "", ""
}
