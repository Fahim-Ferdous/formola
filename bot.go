package main

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func botReady(_ *discordgo.Session, m *discordgo.Ready) {
	fmt.Println("Logged in as:", m.User)
}

func runBot(discordAuthToken string) chan<- struct{} {
	dg, err := discordgo.New("Bot " + discordAuthToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil
	}

	dg.AddHandler(botReady)

	dg.Identify.Intents = discordgo.IntentsDirectMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return nil
	}

	go func() {
		for mq := range msgqueue {
			// TODO: Get channel ID from forms table in DB.
			// TODO: Send embeds, rather than JSON :/
			// TODO: Warn about jumping to links sent through the message.
			b, _ := json.MarshalIndent(mq.Obj, "", "  ")
			dg.ChannelMessageSend(testChannelID, fmt.Sprintf("```json\n%s\n```", string(b)))
		}
	}()
	shutdownBot := make(chan struct{})
	go func() {
		<-shutdownBot
		if err := dg.Close(); err != nil {
			fmt.Println("Err shutdownBot:", err)
		}
	}()
	return shutdownBot
}
