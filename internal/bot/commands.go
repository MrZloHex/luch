package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "log/slog"
	"luch/pkg/util"
	"strings"
)

type Commands []tgbotapi.BotCommand

var def_cmds = Commands{
	{Command: "start", Description: "Start the bot"},
	{Command: "help", Description: "Show help"},
	{Command: "menu", Description: "Show reply keyboard"},
	{Command: "hide", Description: "Hide reply keyboard"},
	{Command: "notify", Description: "Toggle notification for me"},
}

func (bot *Bot) fetchCommands() error {
	cmds, err := bot.api.GetMyCommands()
	if err != nil {
		return err
	}

	same := util.EqualSlices(cmds, def_cmds, func(x, y tgbotapi.BotCommand) bool {
		return x.Command == y.Command && x.Description == y.Description
	}, true)

	if same {
		bot.cmds = cmds
		return nil
	}
	log.Warn("Commands are not set properly, should be", "cmd", def_cmds)

	_, err = bot.api.Request(tgbotapi.NewSetMyCommands(def_cmds...))
	if err != nil {
		log.Error("Failed to set default commands", "err", err)
		return err
	}

	bot.cmds = def_cmds
	return nil
}

func (bot *Bot) processCmd(upd tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")

	switch upd.Message.Command() {
	case "start":
		msg.Text = "Welcome to *LUCH*! Luch is bot for monlith system.\nUse /help to see commands."
		msg.ReplyMarkup = bot.kb.kb
	case "help":
		msg.Text = bot.buildHelp()
	case "menu":
		msg.Text = "Here’s the menu:"
		msg.ReplyMarkup = bot.kb.kb
	case "hide":
		msg.Text = "Keyboard hidden."
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	case "notify":
		msg.Text = bot.toggleNotify(upd.Message.Chat.ID)
		bot.saveNotifiers()
	default:
		msg.Text = "Unknown command"
		log.Warn("Unknown command", "cmd", upd.Message.Command())
		return nil
	}

	_, err := bot.api.Send(msg)
	return err
}

func (bot *Bot) buildHelp() string {
	if len(def_cmds) == 0 {
		return "No commands are configured yet."
	}
	var b strings.Builder
	b.WriteString("*Available commands:*\n")
	for _, c := range def_cmds {
		fmt.Fprintf(&b, "/%s — %s\n", c.Command, c.Description)
	}
	return b.String()
}
