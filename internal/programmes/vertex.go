package programmes

import (
	log "log/slog"
	"luch/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Vertex struct {
	waitingBrightness bool
}

func vertexKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Lamp Off", "VERTEX:LAMP:OFF"),
			tgbotapi.NewInlineKeyboardButtonData("Lamp On", "VERTEX:LAMP:ON"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Led Off", "VERTEX:LED:OFF"),
			tgbotapi.NewInlineKeyboardButtonData("Led Blink", "VERTEX:LED:BLINK"),
			tgbotapi.NewInlineKeyboardButtonData("Led Fade", "VERTEX:LED:FADE"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Led Solid", "VERTEX:LED:SOLID"),
			tgbotapi.NewInlineKeyboardButtonData("Led Bright", "BRIGHT"),
		),
	)
}

func (v *Vertex) Start(conn Conn, upd tgbotapi.Update) error {
	chatID, _, ok := bot.PickIDnTXT(upd)
	if !ok {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, "Commands for VERTEX:")
	msg.ReplyMarkup = vertexKeyboard()

	if _, err := conn.Bot.Send(msg); err != nil {
		log.Error("vertex: failed to send start message", "err", err)
		return err
	}
	return nil
}

func (v *Vertex) UpdateBot(conn Conn, upd tgbotapi.Update) error {
	chatID, text, ok := bot.PickIDnTXT(upd)
	if !ok {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, "")

	switch {
	case text == "BRIGHT":
		msg.Text = "Please send 0..=255 LED brightness"
		v.waitingBrightness = true

	case v.waitingBrightness:
		log.Info("vertex: received brightness value", "value", text)
		rx, err := conn.Ptcl.TransmitReceive([]string{"VERTEX:LED:BRIGHT", text})
		resp := rx.String()
		if err != nil {
			log.Error("vertex: transmit failed", "err", err)
			resp = "failed to transmit brightness"
		}
		msg.Text = resp
		v.waitingBrightness = false

	default:
		rx, err := conn.Ptcl.TransmitReceive(text)
		resp := rx.String()
		if err != nil {
			log.Error("vertex: transmit failed", "err", err)
			resp = "failed to transmit command"
		}
		msg.Text = resp
	}

	if _, err := conn.Bot.Send(msg); err != nil {
		log.Error("vertex: failed to send response", "err", err)
	}

	if upd.CallbackQuery != nil {
		if _, err := conn.Bot.Request(tgbotapi.NewCallback(upd.CallbackQuery.ID, "")); err != nil {
			log.Warn("vertex: failed to ack callback", "err", err)
		}
	}

	return nil
}
