package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "log/slog"
)

type Keyboard struct {
	kb     tgbotapi.ReplyKeyboardMarkup
	labels map[string]struct{}
}

func (bot *Bot) setupKeyboard() {
	bot.kb.kb = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Next Effect"),
			tgbotapi.NewKeyboardButton("Lamp On"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Led Off"),
			tgbotapi.NewKeyboardButton("Lamp Off"),
		),
	)
	bot.kb.labels = map[string]struct{}{
		"Lamp On":     {},
		"Lamp Off":    {},
		"Led Off":     {},
		"Next Effect": {},
	}
	bot.kb.kb.ResizeKeyboard = true
	bot.kb.kb.OneTimeKeyboard = false
}

func (bot *Bot) isKeyboard(upd tgbotapi.Update) bool {
	_, ok := bot.kb.labels[upd.Message.Text]
	return ok
}

func (bot *Bot) proccessKeyboard(upd tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")

	switch upd.Message.Text {
	case "Lamp On":
		bot.ptcl.Send("VERTEX", "LAMP:ON")
	case "Lamp Off":
		bot.ptcl.Send("VERTEX", "LAMP:OFF")
	case "Led Off":
		bot.ptcl.Send("VERTEX", "LED:OFF")
	case "Next Effect":
		bot.ptcl.Send("VERTEX", "LED:NEXT")
	default:
		msg.Text = "Unknown msg"
		log.Warn("Unknown msg", "msg", upd.Message.Command())
		return nil
	}

	_, err := bot.api.Send(msg)
	return err
}
