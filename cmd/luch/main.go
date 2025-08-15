package main

import (
	"os"

	"github.com/lmittmann/tint"
	"log"
	"log/slog"

	"github.com/joho/godotenv"

	"luch/pkg/protocol"

	"luch/internal/bot"
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

	ptcl_cfg := protocol.PtclConfig{
		Shard:  "LUCH",
		Url:    "ws://localhost:8092",
		Reconn: 5,
	}

	ptcl, err := protocol.NewProtocol(ptcl_cfg)
	if err != nil {
		slog.Error("Failed to init protocol")
		os.Exit(1)
	}

	cfg := bot.BotConfig{
		Token:  token,
		Debug:  false,
		Logger: stdToSlog,
	}

	bot, err := bot.NewBot(cfg, ptcl)
	if err != nil {
		slog.Error("Failed to init bot", "err", err)
	}

	slog.Info("BOOTING UP", "bot", bot.GetName(), "url", ptcl_cfg.Url)

	bot.Setup()

	go ptcl.Run()
	go bot.Run()

	for {
	}
}
