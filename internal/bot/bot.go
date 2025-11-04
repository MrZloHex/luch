package bot

import (
	stdlog "log"
	log "log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"luch/pkg/stt"

	"time"

	"luch/internal/core"
)

type BotConfig struct {
	Token  string
	Debug  bool
	Logger *stdlog.Logger
	Notify string
}

type Bot struct {
	api *tgbotapi.BotAPI

	tr *stt.Transcriber

	cmds Commands
	kb   Keyboard
	not  Notifier

	out chan<- core.Event
}

func NewBot(cfg BotConfig) (*Bot, error) {
	log.Debug("init telebot")

	bot := Bot{
		not: Notifier{
			notifyFile: cfg.Notify,
		},
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

func (bot *Bot) SetEvent(out chan<- core.Event) {
	bot.out = out
}

func (bot *Bot) Setup() {
	bot.setupNotifier()

	err := bot.fetchCommands()
	if err != nil {
		log.Error("Failed to retrive commads", "err", err)
	}
	log.Debug("Commands from Telegram", "cmd", bot.cmds)

	bot.setupKeyboard()

	bot.tr, err = stt.NewTranscriber("third_party/whisper.cpp/models/ggml-medium.bin")
	if err != nil {
		log.Error("load model: %v", err)
	}

	bot.NotifyAll("Bot started")
}

func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.api.GetUpdatesChan(u)

	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			bot.processCmd(update)
		} else {
			bot.out <- core.Event{
				Kind: core.EV_BOT,
				Bot:  update,
			}
		}
	}
}

func (bot *Bot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return bot.api.Send(c)
}

func (bot *Bot) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	return bot.api.Request(c)
}

func (bot *Bot) GetFile(c tgbotapi.FileConfig) (tgbotapi.File, error) {
	return bot.api.GetFile(c)
}

func (bot *Bot) GetName() string {
	return bot.api.Self.UserName
}

func (bot *Bot) GetToken() string {
	return bot.api.Token
}

func PickIDnTXT(upd tgbotapi.Update) (chatID int64, text string, ok bool) {
	switch {
	case upd.CallbackQuery != nil:
		return upd.CallbackQuery.Message.Chat.ID, upd.CallbackData(), true
	case upd.Message != nil:
		return upd.Message.Chat.ID, upd.Message.Text, true
	default:
		return 0, "", false
	}
}
