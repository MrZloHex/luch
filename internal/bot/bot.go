package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	stdlog "log"
	log "log/slog"

	"luch/pkg/protocol"
)

type BotConfig struct {
	Token  string
	Debug  bool
	Logger *stdlog.Logger
}

type Bot struct {
	api *tgbotapi.BotAPI

	ptcl *protocol.Protocol
}

func NewBot(cfg BotConfig, ptcl *protocol.Protocol) (*Bot, error) {
	log.Debug("init telebot")

	bot := Bot{
		ptcl: ptcl,
	}

	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Error("Failed to init api", "err", err)
		return nil, err
	}
	tgbotapi.SetLogger(cfg.Logger)
	api.Debug = cfg.Debug

	bot.api = api

	log.Debug("authorised telebot")

	return &bot, nil
}

func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Info("Got msg", "from", update.Message.From.UserName, "text", update.Message.Text)

			bot.ptcl.Write("VERTEX", update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.api.Send(msg)
		}
	}
}

func (bot *Bot) GetName() string {
	return bot.api.Self.UserName
}
