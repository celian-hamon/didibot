package bot

import (
	"discordbot/config"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ratio(arg []string, channelID string, messageID string, reply string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	type data struct {
		MaxAge    int
		MaxUses   int
		Temporary bool
		Unique    bool
	}
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

		timer := 1 * time.Minute
		t := time.Now().Add(timer)
		for {
			if time.Now().Before(t) {
				test, _ := s.MessageReactions(channelID, ratio.ID, "👍", 99, "", "1")

				fmt.Println(len(test))
				if len(test) > 1 {
					_, err := s.ChannelMessageSendReply(channelID, "Tu t'es fais ratio boloss", m.MessageReference)
					if err != nil {
						return err.Error(), ""
					}
					prv, err := s.UserChannelCreate(reply.Author.ID)
					if err != nil {
						return err.Error(), ""
					}
					_, err = send(prv.ID, config.Invit)
					if err != nil {
						return err.Error(), ""
					}
					err = s.GuildMemberDeleteWithReason(m.GuildID, reply.Author.ID, "tu t'es fais ratio mdr")
					if err != nil {
						return err.Error(), ""
					}
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
