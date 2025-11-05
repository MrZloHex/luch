package main

import (
	"os"
	"time"

	"log"
	"log/slog"

	"github.com/lmittmann/tint"

	cli "github.com/spf13/pflag"

	"github.com/joho/godotenv"

	"luch/internal/bot"
	"luch/internal/core"
	"luch/internal/luch"
	"luch/pkg/protocol"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func main() {
	envFile := cli.StringP("env", "e", ".env", "Env file path")
	url := cli.StringP("url", "u", "ws://localhost:8092", "Url of hub")
	logLevel := cli.StringP("log", "l", "info", "Log level")
	notifier := cli.StringP("json", "j", "notify.json", "Path to JSON where locate chat id for notification")
	botDebug := cli.Bool("bot-debug", false, "Enable debug output for bot")
	cli.Parse()

	log_handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level: logLevelMap[*logLevel],
	})
	slog.SetDefault(slog.New(log_handler))

	stdToSlog := slog.NewLogLogger(log_handler, slog.LevelDebug)
	log.SetFlags(0)
	log.SetOutput(stdToSlog.Writer())

	godotenv.Load(*envFile)

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		slog.Error("Failed to get TOKEN")
		os.Exit(1)
	}
	whitelist := os.Getenv("WHITELIST")
	if whitelist == "" {
		slog.Error("Failed to get WhiteList")
		os.Exit(1)
	}

	ptcl_cfg := protocol.PtclConfig{
		Shard:   "LUCH",
		Url:     *url,
		Reconn:  5,
		Timeout: 3 * time.Second,
	}

	ptcl, err := protocol.NewProtocol(ptcl_cfg)
	if err != nil {
		slog.Error("Failed to init protocol")
		os.Exit(1)
	}

	cfg := bot.BotConfig{
		Token:     token,
		Debug:     *botDebug,
		Logger:    stdToSlog,
		Notify:    *notifier,
		WhiteList: whitelist,
	}

	bot, err := bot.NewBot(cfg)
	if err != nil {
		slog.Error("Failed to init bot", "err", err)
	}
	bot.Setup()

	slog.Info("BOOTING UP", "bot", bot.GetName(), "url", ptcl_cfg.Url)

	luch, err := luch.Init(bot, ptcl)
	bot.SetEvent(luch.GetEventChan())
	proto := core.NewProtoEvHandler(luch.GetEventChan())
	ptcl.EmitOut(proto.EmitOut)

	go ptcl.Run()
	go bot.Run()

	for {
		luch.Run()
	}
}
