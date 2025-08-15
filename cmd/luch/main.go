package main

import (
	"os"

	"github.com/lmittmann/tint"
	"log"
	"log/slog"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"luch/pkg/protocol"
)

func main() {
	log_handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(log_handler))

	stdToSlog := slog.NewLogLogger(log_handler, slog.LevelDebug)
	log.SetFlags(0)
	log.SetOutput(stdToSlog.Writer())

	godotenv.Load()

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		slog.Error("Failed to get TOKEN")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error("Failed to init bot", "err", err)
	}
	tgbotapi.SetLogger(stdToSlog)

	// bot.Debug = true

	ptcl, err := protocol.NewProtocol("LUCH", "ws://localhost:8092")
	if err != nil {
		slog.Error("Failed to init protocol")
		os.Exit(1);
	}
	_ = ptcl;

	slog.Info("BOOTING UP ON", "bot", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			slog.Info("Got msg", "from", update.Message.From.UserName, "text", update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
