package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Fefefo/moeScraper/scraper"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/ini.v1"
)

var animeList scraper.List

func main() {
	cfg, err := ini.Load("my.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
	}

	discapi := cfg.Section("").Key("disc_api").String()

	bot, err := discordgo.New("Bot " + discapi)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	animeList = scraper.GetAnimeList()

	bot.AddHandler(messageCreate)
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening Discord session,", err)
		return
	}

	fmt.Println("Bot online!")

	bot.UpdateListeningStatus("「AHEGAO」🌸 DO THE AHEGAO 🌸")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := "^"

	if string([]rune(m.Content)[0:1]) == prefix {
		splittedText := strings.Split(m.Content, " ")

		key := splittedText[0]
		query := strings.Join(splittedText[1:], " ")

		if key == prefix+"anime" {

			if len(query) > 2 {
				if len(animeList.SelectByBothNames(query)) > 0 {
					animeList = animeList.SelectByBothNames(query)

					var fields []*discordgo.MessageEmbedField

					for i, k := 0, 0; k < 10 && i < len(animeList); i++ {
						for j := 0; j < len(animeList[i].Songs) && k < 10; j, k = j+1, k+1 {
							fields = append(fields, &discordgo.MessageEmbedField{
								Name:  strconv.Itoa(k+1) + ") " + animeList[i].NameJap + " - " + animeList[i].Songs[j].Version + " - " + animeList[i].Songs[j].Title,
								Value: animeList[i].Songs[j].Link,
							})
						}
					}
					emb := discordgo.MessageEmbed{
						Title:  "Anime themes",
						Fields: fields,
					}

					s.ChannelMessageSendEmbed(m.ChannelID, &emb)
				} else {
					s.ChannelMessageSend(m.ChannelID, "Anime non trovato")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Almeno 3 lettere")
			}
		}
	}
}
