package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load("secrets.env")
}

var BotID string
var discord *discordgo.Session
var char CharacterResponse

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
		_, _ = s.ChannelMessageSend(m.ChannelID, GetCharacterInfo("Atan Sarbeth"))
	}
}

func GetCharacterInfo(CharName string) string {
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

	return char.Characters.Data.Name
}
