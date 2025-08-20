package programmes

import (
	"fmt"
	_ "log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Vertex struct {
	brightness bool
}

func (v *Vertex) callback(m Messanger, upd tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "")

	if upd.CallbackData() == "BRIGHT" {
		msg.Text = "Please send 0..=255 led brightness"
		v.brightness = true
	} else {
		msg.Text = m.SendWS(upd.CallbackData())
	}

	_, err := m.SendBot(msg)
	m.RequestBot(tgbotapi.NewCallback(upd.CallbackQuery.ID, ""))
	return err
}

func vertexMakeKB() tgbotapi.InlineKeyboardMarkup {
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

func (v *Vertex) Execute(m Messanger, upd tgbotapi.Update) error {
	if upd.CallbackQuery != nil {
		return v.callback(m, upd)
	}

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
	if v.brightness {
		msg.Text = m.SendWS(fmt.Sprintf("VERTEX:LED:BRIGHT:%s", upd.Message.Text))
		v.brightness = false
	} else {
		msg.Text = "Commands for VERTEX:"
		msg.ReplyMarkup = vertexMakeKB()
	}

	_, err := m.SendBot(msg)
	return err
}
