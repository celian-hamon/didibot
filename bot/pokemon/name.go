package pokemon

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

var numberMap = map[int]string{
	1: ":one:",
	2: ":two:",
	3: ":three:",
	4: ":four:",
	5: ":five:",
	6: ":six:",
	7: ":seven:",
	8: ":eight:",
	9: ":nine:",
}

func Name(arg []string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	if !isAlreadyStarter(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Tu n'as aucun pokemon ! Tu peux en avoir un avec `o!p-starter`")
		return "", ""
	}
	if len(arg) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Tu dois préciser un surnom de pokémon et le pokemon !")
		return "", ""
	} else {
		pokemonName := arg[1]
		pokemonSurname := arg[2]

		db, _ := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/pokemon")

		defer db.Close()

		rows, err := db.Query("SELECT id FROM `pokemon_table` WHERE `name` = ? AND `surname` is null", pokemonName)

		if err != nil {
			panic(err)
		}
		var pokemonId []string
		for rows.Next() {
			var id string
			rows.Scan(&id)
			fmt.Println(id)
			pokemonId = append(pokemonId, id)
		}

		if len(pokemonId) == 1 {
			_, err := db.Query("UPDATE `pokemon_table` SET `surname` = ? WHERE `id` = ?", pokemonSurname, pokemonId[0])
			if err != nil {
				fmt.Println(err.Error())
			}
			s.MessageReactionAdd(m.ChannelID, m.ID, ":ok_hand:")
			return "", ""
		} else if len(pokemonId) > 1 {
			reponse, _ := s.ChannelMessageSend(m.ChannelID, "Quel pokémon souhaites-tu renommer ?")
			for i := range pokemonId {
				s.MessageReactionAdd(m.ChannelID, reponse.ID, numberMap[i+1])
			}
			for {
				var reactionArray [][]*discordgo.User
				for count := range pokemonId {
					test, _ := s.MessageReactions(m.ChannelID, reponse.ID, numberMap[count], 99, "", "1")
					reactionArray = append(reactionArray, test)
					for _, user := range reactionArray {
						if len(user) > 1 {
							for _, j := range user {
								if j.ID == m.Author.ID {
									s.ChannelMessageSend(m.ChannelID, "Tu demande le pokémon numéro "+numberMap[count]+" !")
									return "", ""
								}
							}
						}
					}
				}
			}
		}
	}
	return "", ""
}
