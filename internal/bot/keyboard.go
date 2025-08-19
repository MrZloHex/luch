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
			tgbotapi.NewKeyboardButton("/vertex"),
		),
	)
	bot.kb.kb.ResizeKeyboard = true
	bot.kb.kb.OneTimeKeyboard = false
}

func (bot *Bot) proccessInlineKeyboard(upd tgbotapi.Update) error {
	log.Debug("CALLBACK")

	msg := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "")
	msg.Text = bot.SendReq(upd.CallbackData())
	_, err := bot.api.Send(msg)
	bot.api.Request(tgbotapi.NewCallback(upd.CallbackQuery.ID, ""))
	return err
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
	}

	return kb
}
