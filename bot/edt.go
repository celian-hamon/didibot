package bot

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type EDT struct {
	CodeSeance        int    `json:"CodeSeance"`
	NomSession        string `json:"NomSession"`
	NomMatiere        string `json:"NomMatiere"`
	IntervenantNom    string `json:"IntervenantNom"`
	IntervenantPrenom string `json:"IntervenantPrenom"`
	DebutSeance       string `json:"DebutSeance"`
	FinSeance         string `json:"FinSeance"`
	NomSalle          string `json:"NomSalle"`
}
type envoie struct {
	NomMatiere        string
	NomSession        string
	NomSalle          string
	IntervenantNom    string
	IntervenantPrenom string
	DebutSeance       string
	FinSeance         string
}

func edt(arg []string, channelID string, messageID string, m *discordgo.MessageCreate, s *discordgo.Session) (string, string) {
	send := s.ChannelMessageSend
	if len(arg) == 2 {
		switch arg[1] {
		case "help":
			send(channelID, "salut")
			return "", ""
		case "next":
			if _, err := parseEdt(); err != nil {
				send(channelID, "erreur")
				return "", ""
			} else {

				return "", ""
			}
			send(channelID, "next")
			return "", ""
		case "reload":
			send(channelID, "reloading")
			if _, err := reload(); err != nil {
				send(channelID, "error")
				return "", ""
			} else {
				send(channelID, "reloaded")
				return "", ""
			}
		}
	} else {
		send(channelID, "veuillez mettre un argument valide")
	}
	return "", ""
}

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
				var data []EDT
				if err := json.Unmarshal(body, &data); err != nil {
					log.Println(err)
					return "", err
				} else {
					file, _ := json.MarshalIndent(data, "", " ")
					err = ioutil.WriteFile("./bot/edt/edt.json", file, 0644)
					if err != nil {
						log.Println(err)
						return "", err
					}
					fmt.Println(data)

				}
			}
		}
	}

	return "", nil
}

func parseEdt() (envoie, error) {
	file, err := ioutil.ReadFile("./bot/edt/edt.json")
	if err != nil {
		log.Println(err)
		return envoie{}, err
	}
	var data []EDT
	if err := json.Unmarshal(file, &data); err != nil {
		log.Println(err)
		return envoie{}, err
	} else {
		lundi, _ := date()
		for _, v := range data {
			if strings.Contains(v.DebutSeance, lundi) {

				//fmt.Println(v.NomMatiere, v.NomSession, v.NomSalle, v.IntervenantNom, v.IntervenantPrenom, v.DebutSeance, v.FinSeance)
				pp := envoie{
					NomMatiere:        v.NomMatiere,
					NomSession:        v.NomSession,
					NomSalle:          v.NomSalle,
					IntervenantNom:    v.IntervenantNom,
					IntervenantPrenom: v.IntervenantPrenom,
					DebutSeance:       v.DebutSeance,
					FinSeance:         v.FinSeance,
				}

				datee, _ := time.Parse("2006-01-02", lundi)
				for i := 0; i < 5; i++ {
					jour := datee.AddDate(0, 0, i)
					annee, mois, day := jour.Date()
					journee := strconv.Itoa(annee) + "-0" + strconv.Itoa(int(mois)) + "-0" + strconv.Itoa(day)
					for _, v := range data {
						if strings.Contains(v.DebutSeance, journee) {
							//fmt.Println(v.NomMatiere, v.NomSession, v.NomSalle, v.IntervenantNom, v.IntervenantPrenom, v.DebutSeance, v.FinSeance)
							pp = envoie{
								NomMatiere:        v.NomMatiere,
								NomSession:        v.NomSession,
								NomSalle:          v.NomSalle,
								IntervenantNom:    v.IntervenantNom,
								IntervenantPrenom: v.IntervenantPrenom,
								DebutSeance:       v.DebutSeance,
								FinSeance:         v.FinSeance,
							}
							return pp, nil
						}

					}

				}
			}
		}
	}

	return envoie{}, nil
}
func date() (string, error) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	for i := 0; i <= lastOfMonth.Day(); i++ {
		day := time.Date(currentYear, currentMonth, i, 0, 0, 0, 0, currentLocation).Weekday()
		if day == time.Monday {
			fmt.Println(i)
			fmt.Println(currentYear, "-0", int(now.Month()), "-0", i)
			lundi := strconv.Itoa(currentYear) + "-0" + strconv.Itoa(int(now.Month())) + "-0" + strconv.Itoa(i)
			fmt.Println(lundi)
			return lundi, nil
		}
	}
	return "", nil
}
