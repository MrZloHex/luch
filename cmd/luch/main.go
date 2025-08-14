package main

import (
	"os"

	"github.com/lmittmann/tint"
	log "log/slog"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	log.SetDefault(log.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level: log.LevelDebug,
		}),
	))

	godotenv.Load()
	
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Error("Failed to get TOKEN");
		os.Exit(1);
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Error("Failed to init bot", "err", err)
	}

	// bot.Debug = true

	log.Info("BOOTING UP ON", "bot", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Info("Got msg", "from", update.Message.From.UserName, "text", update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
