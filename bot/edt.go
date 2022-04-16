package bot

import (
	"discordbot/security"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type EDTGlobal struct {
	CodeSeance        int    `json:"CodeSeance"`
	NomSession        string `json:"NomSession"`
	NomMatiere        string `json:"NomMatiere"`
	IntervenantNom    string `json:"IntervenantNom"`
	IntervenantPrenom string `json:"IntervenantPrenom"`
	DebutSeance       string `json:"DebutSeance"`
	FinSeance         string `json:"FinSeance"`
	NomSalle          string `json:"NomSalle"`
}

type EDTSend struct {
	array []map[string]string
}
type EDTSem struct {
	NomMatiere        string
	IntervenantNom    string
	IntervenantPrenom string
	DebutSeance       string
	FinSeance         string
	NomSalle          string
	NomSession        string
}
type EDTSem2 struct {
	DebutSemaine string
}

func edt(arg []string, channelID string, messageID string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSend
	if len(arg) == 2 || len(arg) == 3 {
		switch arg[1] {
		case "help":
			send(channelID, "Salut si tu veux savoir le prochain EDT, tu peux faire m!edt")
			send(channelID, "Si tu veux savoir les mois prochains, tu peux faire m!edt next 'le nombre de semaine en plus'")
			return "", ""
		case "next":
			if semaine, err := strconv.Atoi(arg[2]); err != nil {
				security.Check("edt", err)
				send(channelID, "Tu dois mettre un nombre après m!edt next")
				return "", err.Error()
			} else {
				fmt.Println(semaine)
				if _, err := parseEdtsem(channelID, messageID, m, s, semaine); err != nil {
					security.Check("edt reload", err)
					send(channelID, "erreur")
				} else {
					return "", ""
				}
			}
		case "reload":
			send(channelID, "reloading")
			if _, err := reload(); err != nil {
				security.Check("edt reload", err)
			} else {
				send(channelID, "reloaded")
				return "", ""
			}
		default:
			send(channelID, "Veuillez mettre un argument valable, pour obtenir de l'aide 'm!edt help'")
		}
	}
	if len(arg) == 1 {
		if _, err := parseEdtsem(channelID, messageID, m, s, 0); err != nil {
			security.Check("edt reload", err)
			send(channelID, "erreur")
		} else {
			return "", ""
		}
	}
	return "", ""
}

// reload the edt from api.
func reload() (string, error) {
	url := "https://api.alternancerouen.fr/planification/session/2290160.json"
	request := http.Client{
		Timeout: time.Second * 10,
	}
	if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
		log.Fatal(err)
	} else {
		if resp, err := request.Do(req); err != nil {
			log.Println(err)
			return "", err
		} else {
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				log.Println(err)
				return "", err
			} else {
				var data []EDTGlobal
				if err := json.Unmarshal(body, &data); err != nil {
					log.Println(err)
					return "", err
				} else {
					os.Remove("./bot/edt/edtglobal.json")

					file, _ := json.MarshalIndent(data, "", " ")
					if err != nil {
						log.Println(err)
						return "", err
					}
					err = ioutil.WriteFile("./bot/edt/edtglobal.json", file, 0644)
					if err != nil {
						log.Println(err)
						return "", err

					}

				}
			}
		}
	}
	parseEdt()
	return "", nil
}

//parse the global edt to make a semaine edt
func parseEdt() (string, error) {
	var result []byte
	var idents []EDTSem
	var data []EDTGlobal
	compteur, semaine := 0, 0
	if file, err := ioutil.ReadFile("./bot/edt/edtglobal.json"); err != nil {
		log.Println(err)
		return "", err
	} else {
		if err := json.Unmarshal(file, &data); err != nil {
			log.Println(err)
			return "", err
		} else {
			for _, v := range data {
				compteur++
				idents = append(idents, EDTSem{NomMatiere: v.NomMatiere, NomSession: v.NomSession, NomSalle: v.NomSalle, IntervenantNom: v.IntervenantNom, IntervenantPrenom: v.IntervenantPrenom, DebutSeance: v.DebutSeance, FinSeance: v.FinSeance})
				if compteur == 10 {
					if result, err = json.Marshal(idents); err != nil {
						security.Check("parseEdt", err)
						return "", err
					} else {
						os.Remove("./bot/edt/edtsemaine" + strconv.Itoa(semaine) + ".json")

						if f, err := os.OpenFile("./bot/edt/edtsemaine"+strconv.Itoa(semaine)+".json", os.O_CREATE|os.O_WRONLY, 0644); err != nil {
							security.Check("parseEdt", err)
							return "", err
						} else {
							defer f.Close()
							if _, err = io.WriteString(f, string(result)); err != nil {
								security.Check("parseEdt", err)
								return "", err
							} else {
								idents = nil
								compteur = 0
								semaine++
								continue
							}

						}

					}

				}

			}
		}
	}
	return "", nil

}

//get date fonction

//parse the semaine edt to send msg to the channel
func parseEdtsem(channelID string, messageID string, m *discordgo.MessageCreate, s *discordgo.Session, semaine int) (string, error) {
	jsonFile, err := os.Open("./bot/edt/edtsemaine" + strconv.Itoa(semaine) + ".json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var data EDTSend
	if err := json.Unmarshal(byteValue, &data.array); err != nil {
		log.Println(err)
	} else {

		embed := &discordgo.MessageEmbed{}
		days := [...]string{"Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi"}

		for index, day := range data.array {

			today := time.Now().AddDate(0, 0, 6).Format("2006-01-02")
			if today > day["FinSeance"] {
				semaine++
				parseEdtsem(channelID, messageID, m, s, semaine)
				break
			} else {
				var debut time.Time
				var fin time.Time
				var prof string
				var jour string
				var matiere string
				if day["DebutSeance"] == "" {
					debut, _ = time.Parse("15:04", "00:00")
					jour = "Non défini"
				} else {
					debut, _ = time.Parse("2006-01-02T15:04:05Z", day["DebutSeance"])
					jour = debut.Format("2 Jan 2006")
				}
				if day["FinSeance"] == "" {
					
					fin, _ = time.Parse("15:04", "00:00")

				} else {
					fin, _ = time.Parse("2006-01-02T15:04:05Z", day["FinSeance"])
				}
				if day["IntervenantNom"] == "" || day["IntervenantPrenom"] == "" {
					prof = "Aucun prof défini"
				} else {
					prof = day["IntervenantPrenom"] + " " + day["IntervenantNom"]
				}
				if day["NomSalle"] == "" {
					day["NomSalle"] = "Aucune salle définie"
				}
				if day["NomSession"] == "" || day["NomMatiere"] == "" {
					matiere = "Aucune matière définie"
				} else {
					matiere = day["NomSession"] + " " + day["NomMatiere"]
				}
				fmt.Println(debut, " ", fin, " ", prof, " ", jour, " ", matiere)
				fields := []*discordgo.MessageEmbedField{
					{
						Name:   "Heure",
						Value:  debut.Format("15:04") + " - " + fin.Format("15:04"),
						Inline: false,
					},
					{
						Name:   "Matiere",
						Value:  matiere,
						Inline: true,
					},
					{
						Name:   "Professeur",
						Value:  prof,
						Inline: true,
					},
					{
						Name:   "Salle",
						Value:  day["NomSalle"],
						Inline: true,
					},
				}
				if index%2 == 0 {
					embed = &discordgo.MessageEmbed{
						Footer: &discordgo.MessageEmbedFooter{
							Text: "Made by Paco and Cece",
						},
						Author:    &discordgo.MessageEmbedAuthor{},
						Color:     0x2596be,
						Thumbnail: &discordgo.MessageEmbedThumbnail{},
						Timestamp: time.Now().Format(time.RFC3339),
						Title:     days[index/2] + " " + jour,
					}

					embed.Fields = append(embed.Fields, fields...)

				}
				if index%2 != 0 {
					// Format Reset and Send the embed
					embed.Fields = append(embed.Fields, fields...)
					s.ChannelMessageSendEmbed(channelID, embed)

					embed = &discordgo.MessageEmbed{}
				}
			}
		}

	}
	return "", nil
}
