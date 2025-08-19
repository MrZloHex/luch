package bot

import (
	"fmt"
	stdlog "log"
	log "log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	prgs "luch/internal/programmes"
	"luch/pkg/protocol"

	"io"
	"net/http"
	"os"
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
}

func (bot *Bot) SendWS(parts ...string) string {
	resp, err := bot.ptcl.Send(parts...)
	if err != nil {
		return fmt.Sprintf("Failed to send request: %s", err.Error())
	} else {
		return string(resp)
	}
}

func downloadFile(url string, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
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

func (bot *Bot) GetName() string {
	return bot.api.Self.UserName
}
