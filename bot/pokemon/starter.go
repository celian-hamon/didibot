package pokemon

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mtslzr/pokeapi-go"
)

var starterPokemon = []string{
	"bulbasaur",
	"charmander",
	"squirtle",
}

func Starter(arg []string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSendEmbed
	var messageList []string
	for _, i := range starterPokemon {
		droppedPokemon, err := pokeapi.Pokemon(i)
		if err != nil {
			return err.Error(), ""
		}
		var typeString string
		for i := 0; i < len(droppedPokemon.Types); i++ {
			typeString += " " + droppedPokemon.Types[i].Type.Name
		}
		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       typeColor[droppedPokemon.Types[0].Type.Name], // Green
			Description: strconv.Itoa(droppedPokemon.BaseExperience) + " XP",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Types : ",
					Value:  typeString,
					Inline: false,
				},
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: droppedPokemon.Sprites.FrontDefault,
			},
			Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
			Title:     droppedPokemon.Name,
		}
		message, err := send(m.ChannelID, embed)
		messageList = append(messageList, message.ID)
		s.MessageReactionAdd(m.ChannelID, message.ID, "💰")
		if err != nil {
			return err.Error(), ""
		}
	}
	for {
		for it, i := range messageList {
			test, _ := s.MessageReactions(m.ChannelID, i, "💰", 99, "", "1")

			if len(test) > 1 {
				for _, j := range test {
					if j.ID == m.Author.ID {
						pokemon := starterPokemon[it]
						_, _ = s.ChannelMessageSend(m.ChannelID, "You selected **"+pokemon+"** as your starter pokemon!")
						sprite, _ := pokeapi.Pokemon(pokemon)
						_, _ = s.ChannelMessageSend(m.ChannelID, sprite.Sprites.FrontDefault)
						for _, j := range messageList {
							if j != i {
								_ = s.ChannelMessageDelete(m.ChannelID, j)
							}
						}
						return "", ""
					}
				}
			}
		}
	}

}
