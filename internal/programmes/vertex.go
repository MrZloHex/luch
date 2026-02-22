package programmes

import (
	"fmt"
	"strings"

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
			tgbotapi.NewInlineKeyboardButtonData("Lamp Off", "VERTEX:OFF:LAMP"),
			tgbotapi.NewInlineKeyboardButtonData("Lamp On", "VERTEX:ON:LAMP"),
			tgbotapi.NewInlineKeyboardButtonData("Lamp Toggle", "VERTEX:TOGGLE:LAMP"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Led Off", "VERTEX:OFF:LED"),
			tgbotapi.NewInlineKeyboardButtonData("Led On", "VERTEX:ON:LED"),
			tgbotapi.NewInlineKeyboardButtonData("Led Toggle", "VERTEX:TOGGLE:LED"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Led Solid", "VERTEX:SET:LED:MODE:SOLID"),
			tgbotapi.NewInlineKeyboardButtonData("Led Fade", "VERTEX:SET:LED:MODE:FADE"),
			tgbotapi.NewInlineKeyboardButtonData("Led Blink", "VERTEX:SET:LED:MODE:BLINK"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Led Bright", "BRIGHT"),
			tgbotapi.NewInlineKeyboardButtonData("Buzzer On", "VERTEX:ON:BUZZ"),
			tgbotapi.NewInlineKeyboardButtonData("Buzzer Off", "VERTEX:OFF:BUZZ"),
			tgbotapi.NewInlineKeyboardButtonData("Buzzer Toggle", "VERTEX:TOGGLE:BUZZ"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Get State", "VERTEX:GET:STATE"),
			tgbotapi.NewInlineKeyboardButtonData("Get Uptime", "VERTEX:GET:UPTIME"),
		),
	)
}

func (v *Vertex) fetchState(conn Conn) string {
	gets := []string{
		"VERTEX:GET:LAMP:STATE",
		"VERTEX:GET:LED:STATE",
		"VERTEX:GET:LED:MODE",
		"VERTEX:GET:LED:BRIGHT",
	}
	var b strings.Builder
	for _, cmd := range gets {
		rx, err := conn.Ptcl.TransmitReceive(cmd)
		if err != nil {
			fmt.Fprintf(&b, "%s: failed\n", cmd)
			continue
		}
		fmt.Fprintf(&b, "%s\n", rx.String())
	}
	return strings.TrimSpace(b.String())
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
		rx, err := conn.Ptcl.TransmitReceive([]string{"VERTEX:SET:LED:BRIGHT", text})
		var resp string
		if err != nil {
			log.Error("vertex: transmit failed", "err", err)
			resp = "failed to transmit brightness"
		} else {
			resp = rx.String()
		}
		msg.Text = resp
		v.waitingBrightness = false

	case text == "VERTEX:GET:STATE":
		msg.Text = v.fetchState(conn)

	case text == "VERTEX:GET:UPTIME":
		rx, err := conn.Ptcl.TransmitReceive("VERTEX:GET:UPTIME")
		if err != nil {
			log.Error("vertex: transmit failed", "err", err)
			msg.Text = "failed to get uptime"
		} else {
			msg.Text = rx.String()
		}

	default:
		rx, err := conn.Ptcl.TransmitReceive(text)
		var resp string
		if err != nil {
			log.Error("vertex: transmit failed", "err", err)
			resp = "failed to transmit command"
		} else {
			resp = rx.String()
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
