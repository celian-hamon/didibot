package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
type EDTSem struct {
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
				parseEdtsem(channelID, messageID, m, s)
				return "", ""
			}
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
					file, _ := json.MarshalIndent(data, "", " ")
					err = ioutil.WriteFile("./bot/edt/edtglobal.json", file, 0644)
					if err != nil {
						log.Println(err)
						return "", err
					}
				}
			}
		}
	}

	return "", nil
}

//parse the global edt to make a semaine edt
func parseEdt() (string, error) {
	var result []byte
	var idents []EDTSem
	var journee string
	file, err := ioutil.ReadFile("./bot/edt/edtglobal.json")
	if err != nil {
		log.Println(err)
		return "", err
	}
	var data []EDTGlobal
	if err := json.Unmarshal(file, &data); err != nil {
		log.Println(err)
		return "", err
	} else {
		lundi, _ := date()
		datee, _ := time.Parse("2006-01-02", lundi)
		for i := 0; i < 5; i++ {
			jour := datee.AddDate(0, 0, i)
			annee, mois, day := jour.Date()
			journee = strconv.Itoa(annee) + "-0" + strconv.Itoa(int(mois)) + "-0" + strconv.Itoa(day)
			for _, v := range data {
				if strings.Contains(v.DebutSeance, journee) {
					idents = append(idents, EDTSem{NomMatiere: v.NomMatiere, NomSession: v.NomSession, NomSalle: v.NomSalle, IntervenantNom: v.IntervenantNom, IntervenantPrenom: v.IntervenantPrenom, DebutSeance: v.DebutSeance, FinSeance: v.FinSeance})
					if err != nil {
						log.Println(err)
						return "", err
					}
				}

			}

		}
		result, err = json.Marshal(idents)
		if err != nil {
			log.Println(err)
			return "", err
		}

		f, err := os.OpenFile("./bot/edt/edtsemaine.json", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			return "", err
		}

		_, err = io.WriteString(f, string(result))
		if err != nil {
			log.Println(err)
			return "", err
		}
		f.Close()
		return "", nil
	}

}

//get date fonction
func date() (string, error) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	for i := 0; i <= lastOfMonth.Day(); i++ {
		day := time.Date(currentYear, currentMonth, i, 0, 0, 0, 0, currentLocation).Weekday()
		if day == time.Monday {
			lundi := strconv.Itoa(currentYear) + "-0" + strconv.Itoa(int(now.Month())) + "-0" + strconv.Itoa(i)
			return lundi, nil
		}
	}
	return "", nil
}

//parse the semaine edt to send msg to the channel
func parseEdtsem(channelID string, messageID string, m *discordgo.MessageCreate, s *discordgo.Session) {
	send := s.ChannelMessageSend
	file, err := ioutil.ReadFile("./bot/edt/edtsemaine.json")
	if err != nil {
		log.Println(err)
	}
	var data []EDTSem
	if err := json.Unmarshal(file, &data); err != nil {
		log.Println(err)
	} else {
		for _, v := range data {

			send(channelID, fmt.Sprintf("%s %s %s %s %s %s %s", v.NomMatiere, v.NomSession, v.NomSalle, v.IntervenantNom, v.IntervenantPrenom, v.DebutSeance, v.FinSeance))
		}
	}
}
