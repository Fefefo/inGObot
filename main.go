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
var myid string

func main() {
	cfg, err := ini.Load("my.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
	}

	discapi := cfg.Section("").Key("disc_api").String()
	myid = cfg.Section("").Key("myid").String()

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
		splittedemojiID := strings.Split(m.Content, " ")

		key := splittedemojiID[0]
		query := strings.Join(splittedemojiID[1:], " ")
		if key == prefix+"refreshanime" && m.Author.ID == myid {
			animeList = scraper.GetAnimeList()
			s.ChannelMessageSend(m.ChannelID, "Anime loaded : "+strconv.Itoa(len(animeList)))
		}
		if key == prefix+"anime" {

			if len(query) > 2 {
				if len(animeList.SelectByBothNames(query)) > 0 {
					animeList := animeList.SelectByBothNames(query)

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
					msg, _ := s.ChannelMessageSendEmbed(m.ChannelID, &emb)
					for i := range emb.Fields {
						var emojiID string
						switch i + 1 {
						case 1:
							emojiID = "1️⃣"
						case 2:
							emojiID = "2️⃣"
						case 3:
							emojiID = "3️⃣"
						case 4:
							emojiID = "4️⃣"
						case 5:
							emojiID = "5️⃣"
						case 6:
							emojiID = "6️⃣"
						case 7:
							emojiID = "7️⃣"
						case 8:
							emojiID = "8️⃣"
						case 9:
							emojiID = "9️⃣"
						case 10:
							emojiID = "🔟"
						}
						err := s.MessageReactionAdd(m.ChannelID, msg.ID, emojiID)
						if err != nil {
							fmt.Println(err)
						}
					}
					s.AddHandler(reactionAddForTheTheme)
				} else {
					s.ChannelMessageSend(m.ChannelID, "Anime non trovato")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Almeno 3 lettere")
			}
		}
	}
}

func reactionAddForTheTheme(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if s.State.User.ID == m.UserID {
		return
	}
	msg, _ := s.ChannelMessage(m.ChannelID, m.MessageID)
	var num int
	switch m.Emoji.Name {
	case "1️⃣":
		num = 0
	case "2️⃣":
		num = 1
	case "3️⃣":
		num = 2
	case "4️⃣":
		num = 3
	case "5️⃣":
		num = 4
	case "6️⃣":
		num = 5
	case "7️⃣":
		num = 6
	case "8️⃣":
		num = 7
	case "9️⃣":
		num = 8
	case "🔟":
		num = 9
	}
	if len(msg.Embeds) > 0 && msg.Embeds[0].Title == "Anime themes" {
		if num < len(msg.Embeds[0].Fields) {
			s.MessageReactionsRemoveAll(m.ChannelID, msg.ID)
			// s.ChannelMessageDelete(m.ChannelID, msg.ID)
			s.ChannelMessageSend(m.ChannelID, "**"+msg.Embeds[0].Fields[num].Name+"**\n"+msg.Embeds[0].Fields[num].Value)
		}
	}
}
