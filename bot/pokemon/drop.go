package pokemon

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mtslzr/pokeapi-go"
)

var typeColor = map[string]int{
	"normal":   0xA8A878,
	"fire":     0xF08030,
	"water":    0x6890F0,
	"electric": 0xF8D030,
	"grass":    0x78C850,
	"ice":      0x98D8D8,
	"fighting": 0xC03028,
	"poison":   0xA040A0,
	"ground":   0xE0C068,
	"flying":   0xA890F0,
	"psychic":  0xF85888,
	"steel":    0xB8B8D0,
	"ghost":    0x705898,
	"dragon":   0x7038F8,
	"dark":     0x705848,
	"fairy":    0xEE99AC,
	"rock":     0xB8A038,
}

func Drop(arg []string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSendEmbed
	droppedPokemon, err := pokeapi.Pokemon(trueRandom(1, 1118))

	var typeString string
	for i := 0; i < len(droppedPokemon.Types); i++ {
		typeString += " " + droppedPokemon.Types[i].Type.Name
	}

	if err != nil {
		return err.Error(), ""
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       typeColor[droppedPokemon.Types[0].Type.Name], // Green
		Description: strconv.Itoa(droppedPokemon.BaseExperience) + " XP",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Height : ",
				Value:  strconv.Itoa(droppedPokemon.Height) + "0 cm",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Weigth : ",
				Value:  strconv.Itoa(droppedPokemon.Weight) + "0 grammes",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
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
	_, err = send(m.ChannelID, embed)
	if err != nil {
		return err.Error(), ""
	}
	return "", ""
}

func trueRandom(min, max int) string {
	rand.Seed(time.Now().UnixNano())
	result := strconv.Itoa(rand.Intn(max-min+1) + min)
	return result
}
