package main

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/chat"
)

func main() {

	bot, err := tgbotapi.NewBotAPI("BOT_TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		ctx := context.Background()
		s := openai.NewSession("OPEN_AI_TOKEN")

		client := chat.NewClient(s, "gpt-3.5-turbo")
		resp, err := client.CreateCompletion(ctx, &chat.CreateCompletionParams{
			Messages: []*chat.Message{
				{Role: "user", Content: update.Message.Text},
			},
		})

		if err != nil {
			log.Fatalf("Failed to complete: %v", err)
		}
		test := ""
		for _, choice := range resp.Choices {
			msg := choice.Message
			test = msg.Content
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, test)
		bot.Send(msg)
	}
}
