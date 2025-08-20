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
	)
	bot.kb.kb.ResizeKeyboard = true
	bot.kb.kb.OneTimeKeyboard = false
}

func (bot *Bot) makeInlineKeyboard(to string) tgbotapi.InlineKeyboardMarkup {
	var kb tgbotapi.InlineKeyboardMarkup
	switch to {
	case "vertex":
		kb = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Next Effect", "VERTEX:NEXT:EFFECT"),
				tgbotapi.NewInlineKeyboardButtonData("Lamp On", "VERTEX:LAMP:ON"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Led Off", "VERTEX:LED:OFF"),
				tgbotapi.NewInlineKeyboardButtonData("Lamp Off", "VERTEX:LAMP:OFF"),
			),
		)
	case "scriptorium":
		kb = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("New Note", "VERTEX:NEXT:EFFECT"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Edit Note", "VERTEX:LED:OFF"),
				tgbotapi.NewInlineKeyboardButtonData("Get Notes", "VERTEX:LAMP:OFF"),
			),
		)
	}

	return kb
}
