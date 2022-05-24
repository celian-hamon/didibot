package pokemon

import (
	"database/sql"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

func Pokedex(arg []string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSendEmbed
	if !isAlreadyStarter(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Tu n'as aucun pokemon ! Tu peux en avoir un avec `o!p-starter`")
		return "", ""
	}
	listPokemon := strings.Join(getAllPokemon(m.Author.ID), "\n")
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Fields: []*discordgo.MessageEmbedField{},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://thumbor.sd-cdn.fr/AFsYPznCBrZr2NQsxjxAbWCa4bE=/1200x630/cdn.sd-cdn.fr/comiga/2019/02/application-pokedex-pokemon-card-dex.png",
		},
		Description: listPokemon,
		Timestamp:   time.Now().Format(time.RFC3339),
		Title:       "Pokédex",
	}

	_, err := send(m.ChannelID, embed)
	if err != nil {
		return err.Error(), ""
	}
	return "", ""
}

func getAllPokemon(s string) []string {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/pokemon")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	row, err := db.Query("SELECT id FROM user_table WHERE discordId = ?", s)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	var user_id string
	for row.Next() {
		err = row.Scan(&user_id)
		if err != nil {
			panic(err)
		}
	}

	rows, err := db.Query("SELECT pokemon_id FROM user_pokemon_join WHERE user_id = ?", user_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var listPokemonId []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		listPokemonId = append(listPokemonId, name)
	}

	var listPokemonName []string
	for i := 0; i < len(listPokemonId); i++ {
		row := db.QueryRow("SELECT name FROM pokemon_table WHERE id = ?", listPokemonId[i])
		var name string
		row.Scan(&name)
		listPokemonName = append(listPokemonName, name)
	}

	return listPokemonName
}
