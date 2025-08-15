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

	switch upd.Message.Command() {
	case "Lamp On":
	case "Lamp Off":
	case "Led Off":
	case "Next Effect":
	default:
		msg.Text = "Unknown msg"
		log.Warn("Unknown msg", "msg", upd.Message.Command())
		return nil
	}

	_, err := bot.api.Send(msg)
	return err
}
