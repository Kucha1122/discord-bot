package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const token string = "NzMyOTQyNDgzODc2NzQxMTIw.Xw79SQ.A6Ja8r6tsK_3VbWfn6u_TeNDxUI"

func main() {
	discord, err := discordgo.New("Bot " + token)

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
