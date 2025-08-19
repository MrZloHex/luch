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
	log.Debug("CALLBACK")

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
	msg.Text = bot.SendReq(upd.CallbackData())
	_, err := bot.api.Send(msg)
	return err
}

func (bot *Bot) makeKeyboard(to string) tgbotapi.InlineKeyboardMarkup {
	var kb tgbotapi.InlineKeyboardMarkup
	switch to {
	case "vertex":
		kb = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Next Effect", "VERTEX:NEXT:EFFECT"),
				tgbotapi.NewInlineKeyboardButtonData("Led Off", "VERTEX:LED:OFF"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Lamp On", "VERTEX:LAMP:ON"),
				tgbotapi.NewInlineKeyboardButtonData("Lamp Off", "VERTEX:LAMP:OFF"),
			),
		)
	}

	return kb
}

