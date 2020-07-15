package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load("secrets.env")
}

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("TOKEN"))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = discord.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")

	<-make(chan struct{})
	return
}
