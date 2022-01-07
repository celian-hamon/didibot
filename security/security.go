package security

import (
	"discordbot/config"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func IsAdmin(userID string, adminList []string) bool {
	config.ReadConfig()
	for v := range adminList {
		if userID == adminList[v] {
			return true
		}
	}
	return false
}

func Log(cmd string, arg []string, erro string, reply string, m *discordgo.MessageCreate) {
	f, err := os.OpenFile("logs"+time.Now().Format("20060102")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	config.Check(err)

	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println("AUTHORS: " + m.Author.Username + " " + m.Author.ID)
	logger.Println("COMMAND: " + cmd)
	logger.Println("ARGUMENTS: " + strings.Join(arg[0:], " "))
	if erro != "" {
		logger.Println("Erreur : " + erro)
	}
	if reply != "" {
		logger.Println("Reply : " + reply)
	}
	logger.Println("----------------------------------------")

}
