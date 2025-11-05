package programmes

import (
	log "log/slog"
	"luch/internal/bot"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Control struct {
	waitCmd bool
}

func controlKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Transmit", "TRANSMIT"),
		),
	)
}

func (c *Control) Start(conn Conn, upd tgbotapi.Update) error {
	chatID, _, ok := bot.PickIDnTXT(upd)
	if !ok {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, "Commands for CONTROL:")
	msg.ReplyMarkup = controlKeyboard()

	if _, err := conn.Bot.Send(msg); err != nil {
		log.Error("failed to send start message", "err", err)
		return err
	}
	return nil
}

func (c *Control) UpdateBot(conn Conn, upd tgbotapi.Update) error {
	chatID, text, ok := bot.PickIDnTXT(upd)
	if !ok {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, "")

	switch {
	case text == "TRANSMIT":
		msg.Text = "Send 2 arguments: `TO[space]PAYLOAD`"
		c.waitCmd = true

	case c.waitCmd:
		args := strings.Split(upd.Message.Text, " ")
		rx, err := conn.Ptcl.TransmitReceive(args)
		if err != nil {
			log.Warn("Failed to sent command", "err", err, "args", args)
		}
		msg.Text = rx.String()
		c.waitCmd = false

	default:
		log.Warn("Nothing to do")
	}

	if _, err := conn.Bot.Send(msg); err != nil {
		log.Error("failed to send response", "err", err)
	}

	if upd.CallbackQuery != nil {
		if _, err := conn.Bot.Request(tgbotapi.NewCallback(upd.CallbackQuery.ID, "")); err != nil {
			log.Warn("failed to ack callback", "err", err)
		}
	}

	return nil
}
