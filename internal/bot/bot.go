package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	stdlog "log"
	log "log/slog"

	"luch/pkg/protocol"

	"time"
)

type BotConfig struct {
	Token  string
	Debug  bool
	Logger *stdlog.Logger
}

type Bot struct {
	api *tgbotapi.BotAPI

	ptcl *protocol.Protocol

	cmds Commands
	kb   Keyboard
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

func (bot *Bot) Setup() {
	err := bot.fetchCommands()
	if err != nil {
		log.Error("Failed to retrive commads", "err", err)
	}
	log.Debug("Commands from Telegram", "cmd", bot.cmds)

	bot.setupKeyboard()
}

func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.api.GetUpdatesChan(u)

	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Debug("Got smth", "from", update.Message.From.UserName, "text", update.Message.Text)

		switch {
		case update.Message.IsCommand():
			bot.processCmd(update)
		case bot.isKeyboard(update):
			bot.proccessKeyboard(update)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No such thingy, sorry\nIf you implement it or contact developer\nSee /help")
			bot.api.Send(msg)

		}

	}
}

func (bot *Bot) GetName() string {
	return bot.api.Self.UserName
}
