package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyboard tgbotapi.ReplyKeyboardMarkup

func (bot *Bot) setupKeyboard() {
	bot.kb = Keyboard(tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Turn On"),
			tgbotapi.NewKeyboardButton("Turn Off"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Toggle Amp"),
			tgbotapi.NewKeyboardButton("Switch Effect"),
		),
	))
	bot.kb.ResizeKeyboard = true // nicer sizing
	bot.kb.OneTimeKeyboard = false
}
