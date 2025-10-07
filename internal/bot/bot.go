package bot

import (
	"fmt"
	stdlog "log"
	log "log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	prgs "luch/internal/programmes"
	"luch/pkg/protocol"
	"luch/pkg/stt"

	"time"
)

type BotConfig struct {
	Token  string
	Debug  bool
	Logger *stdlog.Logger
	Notify string
}

type Bot struct {
	api *tgbotapi.BotAPI

	ptcl *protocol.Protocol
	tr   *stt.Transcriber


	prg prgs.Programme

	cmds Commands
	kb   Keyboard
	not  Notifier
}

func NewBot(cfg BotConfig, ptcl *protocol.Protocol) (*Bot, error) {
	log.Debug("init telebot")

	bot := Bot{
		ptcl: ptcl,
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

	bot.prg.Msg = &bot

	return &bot, nil
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
}

func (bot *Bot) SendWS(parts ...string) string {
	resp, err := bot.ptcl.TransmitReceive(parts...)
	if err != nil {
		return fmt.Sprintf("Failed to send request: %s", err.Error())
	} else {
		return string(resp)
	}
}

func (bot *Bot) listenWs() {
	for {

	}
}

func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.api.GetUpdatesChan(u)

	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	go bot.listenWs()

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			bot.processCmd(update)
		}

		bot.prg.Execute(update)
	}
}

func (bot *Bot) SendBot(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return bot.api.Send(c)
}

func (bot *Bot) RequestBot(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	return bot.api.Request(c)
}

func (bot *Bot) GetFileBot(c tgbotapi.FileConfig) (tgbotapi.File, error) {
	return bot.api.GetFile(c)
}

func (bot *Bot) GetName() string {
	return bot.api.Self.UserName
}

func (bot *Bot) GetToken() string {
	return bot.api.Token
}
