package DiscoBot

import (
	"Gomain3/controllers"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/subosito/gotenv"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	owner                    = "204309394966446081"
	wordsChannel             = "470385587518701568"
	postJpWord               = false
	JPWordDelayMinutes       = 5
	lastJPWordTime     int64 = 0
)

type (
	DiscordBot struct {
		Dg  *discordgo.Session
		Api *controllers.ApiController
	}
)

func NewDiscordBot() *DiscordBot {
	gotenv.Load()
	dg, err := discordgo.New("Bot " + os.Getenv("DG_TOKEN"))
	if err != nil {
		fmt.Printf("Could Not Create Discord Bot: %s\n", err.Error())
		return nil
	}

	err = dg.Open()

	dg.AddHandler(MessageHandler)
	if err != nil {
		fmt.Printf("Error Connecting Bot: %s\n", err.Error())
		return nil
	}
	fmt.Println("Bot is now running.")

	return &DiscordBot{Dg: dg}
}
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == "!giveMeMyFuckingId" {
		s.ChannelMessageSend(m.ChannelID, "Fine: "+m.Author.ID)
	}
	if m.Content == "!channelId" {
		s.ChannelMessageSend(m.ChannelID, m.ChannelID)
	}
	if m.Content[0] != '!' {
		return
	}
	cmd := strings.Split(m.Content, " ")
	if cmd == nil {
		return
	}

	if m.Author.ID == owner {
		if cmd[0] == "!on" {
			postJpWord = true
		}
		if cmd[0] == "!off" {
			postJpWord = false
		}
		if cmd[0] == "!delay" {
			if len(cmd) < 2 {
				s.ChannelMessageSend(m.ChannelID, "Invalid Delay")
				return
			}
			t, err := strconv.Atoi(cmd[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid Delay")
			}
			JPWordDelayMinutes = t

		}

	}
}
func (bot *DiscordBot) SetAPIController(api *controllers.ApiController) {
	bot.Api = api
}
func (bot *DiscordBot) Run() {
	for {
		if postJpWord {
			if time.Now().Unix() > int64(lastJPWordTime+int64(JPWordDelayMinutes*60)) {
				lastJPWordTime = time.Now().Unix()
				bot.Dg.ChannelMessageSend(wordsChannel, "@spades "+bot.Api.RandomWord())
			}
		}
		time.Sleep(10 * time.Second)
	}
}
