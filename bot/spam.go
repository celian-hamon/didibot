package bot

import (
	"discordbot/config"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func spam(arg []string, channelID string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSend

	nombreSpam, stringErr := strconv.Atoi(arg[1])
	if stringErr != nil {
		_, _ = send(channelID, m.Author.Username+", veuillez renseigner le nombre de spam a éffectuer comme ceci : "+arg[0]+" <nombre de spam> <texte a spam>")
	}
	arg[1] = ""
	if len(arg) > 1 {
		for i := 0; i < nombreSpam; i++ {
			_, err := send(channelID, strings.Join(arg[1:], " "))
			config.Check(err)
		}
	} else {
		_, err := send(channelID, m.Author.Username+", veuillez mettre ce que vous voulez spam.")
		config.Check(err)
	}
	return "", ""
}
