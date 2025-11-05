package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyboard struct {
	kb     tgbotapi.ReplyKeyboardMarkup
	labels map[string]struct{}
}

func (bot *Bot) setupKeyboard() {
	bot.kb.kb = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/vertex"),
			tgbotapi.NewKeyboardButton("/notes"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/achtung"),
			tgbotapi.NewKeyboardButton("/control"),
		),
	)
	bot.kb.kb.ResizeKeyboard = true
	bot.kb.kb.OneTimeKeyboard = false
}
