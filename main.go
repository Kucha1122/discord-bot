package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load("secrets.env")
}

var BotID string
var discord *discordgo.Session
var char *CharacterResponse

type CharacterResponse struct {
	Characters struct {
		Data struct {
			Name              string `json:"name"`
			Title             string `json:"title"`
			Sex               string `json:"sex"`
			Vocation          string `json:"vocation"`
			Level             int    `json:"level"`
			AchievementPoints int    `json:"achievement_points"`
			World             string `json:"world"`
			Residence         string `json:"residence"`
			LastLogin         []struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"last_login"`
			AccountStatus string `json:"account_status"`
			Status        string `json:"status"`
		} `json:"data"`
		Achievements []interface{} `json:"achievements"`
		Deaths       []struct {
			Date struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"date"`
			Level    int           `json:"level"`
			Reason   string        `json:"reason"`
			Involved []interface{} `json:"involved"`
		} `json:"deaths"`
		AccountInformation struct {
			LoyaltyTitle string `json:"loyalty_title"`
			Created      struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"created"`
		} `json:"account_information"`
		OtherCharacters []struct {
			Name   string `json:"name"`
			World  string `json:"world"`
			Status string `json:"status"`
		} `json:"other_characters"`
	} `json:"characters"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("TOKEN"))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := discord.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID
	discord.AddHandler(CommandHandler)
	err = discord.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")

	<-make(chan struct{})
	return
}

func CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	}

	if m.Content == "Siema" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "elo")
	}

	if m.Content == "!char" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Musisz podac nazwe postaci np. !char Uther Morlenfra")
	}

	if strings.Contains(m.Content, "!char") && len(m.Content) > len("!char") {
		GetCharacterInfo(After(m.Content, "!char"))
		_, _ = s.ChannelMessageSend(m.ChannelID, PrintCharacterInfo())
	}
}

func GetCharacterInfo(CharName string) {
	response, err := http.Get("https://api.tibiadata.com/v2/characters/" + CharName + ".json")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	json.Unmarshal([]byte(body), &char)
}

func PrintCharacterInfo() string {

	BasicCharInfo :=
		"```apache" +
			"\nName: " + char.Characters.Data.Name +
			" " + strings.ToUpper(char.Characters.Data.Status) +
			"\nTitle:" + char.Characters.Data.Title +
			"\nSex:" + char.Characters.Data.Sex +
			"\nVocation:" + char.Characters.Data.Vocation +
			"\nLevel:" + strconv.Itoa(char.Characters.Data.Level) +
			"\nAchievement Points:" + strconv.Itoa(char.Characters.Data.AchievementPoints) +
			"\nWorld:" + char.Characters.Data.World +
			"\nResidence:" + char.Characters.Data.Residence +
			"\nAccount Status:" + char.Characters.Data.AccountStatus +
			"```"
	if char.Characters.Data.Name == "" {
		return "Character does not exist."
	}

	if len(char.Characters.Deaths) != 0 {
		CharacterDeaths := "\n"

		for _, death := range char.Characters.Deaths {
			CharacterDeaths += string(death.Date.Date) + ", " + string(death.Date.Timezone) + " " + death.Reason + " at Level " + strconv.Itoa(death.Level) + ".\n"
		}

		CharacterDeaths = "```cs" + "\n" + CharacterDeaths + "\n```"

		char = nil
		return BasicCharInfo + CharacterDeaths
	}

	char = nil

	return BasicCharInfo
}

func After(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}

	adjustedPos := pos + len(a)

	if adjustedPos >= len(value) {
		return ""
	}

	return value[adjustedPos:len(value)]
}
