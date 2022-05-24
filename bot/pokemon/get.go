package pokemon

import (
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mtslzr/pokeapi-go"
)

func Get(arg []string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	if !isAlreadyStarter(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Tu n'as pas de starter ! Tu peux en avoir un avec `o!p-starter`")
		return "", ""
	}
	send := s.ChannelMessageSendEmbed
	droppedPokemon, err := pokeapi.Pokemon(strings.ToLower(arg[1]))

	var typeString string
	for i := 0; i < len(droppedPokemon.Types); i++ {
		typeString += " " + droppedPokemon.Types[i].Type.Name
	}

	if err != nil {
		return err.Error(), ""
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  typeColor[droppedPokemon.Types[0].Type.Name], // Green
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Height : ",
				Value:  strconv.Itoa(droppedPokemon.Height) + "0 cm",
				Inline: true,
			},
			{
				Name:   "Weigth : ",
				Value:  strconv.Itoa(droppedPokemon.Weight) + "0 grammes",
				Inline: true,
			},
			{
				Name:   "Hp : ",
				Value:  strconv.Itoa(droppedPokemon.Stats[0].BaseStat) + " hp",
				Inline: true,
			},
			{
				Name:   "Xp : ",
				Value:  strconv.Itoa(droppedPokemon.BaseExperience) + " xp",
				Inline: true,
			},
			{
				Name:   "Types : ",
				Value:  typeString,
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: droppedPokemon.Sprites.FrontDefault,
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     droppedPokemon.Name,
	}
	_, err = send(m.ChannelID, embed)
	if err != nil {
		return err.Error(), ""
	}
	return "", ""
}
