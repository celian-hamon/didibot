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
			return "", ""
		case "next":
			argm, _ := time.Parse("January", arg[2])
			mounths := argm.Month()
			fmt.Println(mounths)
			if _, err := parseEdt(channelID, messageID, m, s, mounths); err != nil {
				send(channelID, "erreur")
				return "", ""
			} else {
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
		case "creneau":
			send(channelID, "creneau")
			_, err := creneau()
			if err != nil {
				send(channelID, "error")
				return "", ""
			} else {
				send(channelID, "ok 😔")
				return "", ""
			}

		}
	}
	if len(arg) == 1 {
		if _, err := parseEdt(channelID, messageID, m, s, 0); err != nil {
			send(channelID, "erreur")
			return "", ""
		} else {
			fmt.Println("ok2")
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
					file, _ := json.MarshalIndent(data, "", " ")
					err = ioutil.WriteFile("./bot/edt/edtglobal.json", file, 0644)
					creneau()
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
func creneau() (string, error) {

	file, _ := ioutil.ReadFile("./bot/edt/edtglobal.json")
	var data []EDTGlobal
	var result []byte

	var idents []EDTSem2
	if err := json.Unmarshal(file, &data); err != nil {
		log.Println(err)
		return "", err
	}
	i := 0
	for _, v := range data {
		i++
		if i == 1 {
			idents = append(idents, EDTSem2{
				DebutSemaine: v.DebutSeance,
			})
		}
		if i == 10 {
			i = 0
		}

	}

	result, err := json.Marshal(idents)
	security.Check("creneau", err)
	f, err := os.OpenFile("./bot/edt/creneau.json", os.O_CREATE|os.O_WRONLY, 0644)
	security.Check("creneau", err)
	defer f.Close()
	_, err = io.WriteString(f, string(result))
	security.Check("creneau", err)

	return "", nil
}

//parse the global edt to make a semaine edt
func parseEdt(channelID string, messageID string, m *discordgo.MessageCreate, s *discordgo.Session, mounth time.Month) (string, error) {
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
		lundi, _ := date(0)
		datee, _ := time.Parse("2006-01-02", lundi)
		vend := datee.AddDate(0, 0, 4)
		fmt.Println(vend.Day() < time.Now().Day())
		if mounth != 0 {
			fmt.Println("fierjiofjerio", mounth)
			lundi, _ = date(mounth - time.Now().Month())
			fmt.Println("cc2", lundi)
			datee, _ = time.Parse("2006-01-02", lundi)
			vend = datee.AddDate(0, 0, 4)
			fmt.Println(datee)
			fmt.Println("lundi", lundi)
			fmt.Println("vend", vend.Format("2006-01-02"))
		} else {
			mounth = time.Now().Month()
		}

		for i := 0; i < 5; i++ {
			jour := datee.AddDate(0, 0, i)
			annee, mois, day := jour.Date()
			journee = strconv.Itoa(annee) + "-0" + strconv.Itoa(int(mois)) + "-0" + strconv.Itoa(day)
			fmt.Println("journee", journee)
			for _, v := range data {
				if strings.Contains(v.DebutSeance, journee) {
					idents = append(idents, EDTSem{NomMatiere: v.NomMatiere, NomSession: v.NomSession, NomSalle: v.NomSalle, IntervenantNom: v.IntervenantNom, IntervenantPrenom: v.IntervenantPrenom, DebutSeance: v.DebutSeance, FinSeance: v.FinSeance})
					security.Check("edt", err)
				}

			}

		}
		result, err = json.Marshal(idents)
		security.Check("edt", err)
		fmt.Println("gros caca", mounth)
		f, err := os.OpenFile("./bot/edt/edtsemaine"+mounth.String()+".json", os.O_CREATE|os.O_WRONLY, 0644)
		security.Check("edt", err)
		defer f.Close()
		_, err = io.WriteString(f, string(result))
		security.Check("edt", err)

	}
	parseEdtsem(channelID, messageID, m, s, mounth)
	return "", nil
}

//get date fonction
func date(add time.Month) (string, error) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	if add != 0 {
		currentMonth = currentMonth + add
		fmt.Println(currentMonth)
	}
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	fmt.Println(firstOfMonth)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	for i := 0; i <= lastOfMonth.Day(); i++ {
		day := time.Date(currentYear, currentMonth, i, 0, 0, 0, 0, currentLocation).Weekday()
		if day == time.Monday {
			lundi := strconv.Itoa(currentYear) + "-0" + strconv.Itoa(int(currentMonth)) + "-0" + strconv.Itoa(i)
			return lundi, nil
		}
	}
	return "", nil
}

//parse the semaine edt to send msg to the channel
func parseEdtsem(channelID string, messageID string, m *discordgo.MessageCreate, s *discordgo.Session, mounth time.Month) {
	fmt.Println("cc", mounth)
	jsonFile, err := os.Open("./bot/edt/edtsemaine" + mounth.String() + ".json")
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
			debut, _ := time.Parse("2006-01-02T15:04:05Z", day["DebutSeance"])
			fin, _ := time.Parse("2006-01-02T15:04:05Z", day["FinSeance"])
			matiere := day["NomSession"] + " " + day["NomMatiere"]
			prof := day["IntervenantNom"] + " " + day["IntervenantPrenom"]
			jour := debut.Format("2 Jan 2006")
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
