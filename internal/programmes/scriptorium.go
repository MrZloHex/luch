package programmes

import (
	log "log/slog"
	"luch/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Scriptorium struct {
	waitAudio bool
}

func scriptoriumKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Transcribe", "TRANSCRIPTION"),
		),
	)
}

func (s *Scriptorium) Start(conn Conn, upd tgbotapi.Update) error {
	chatID, _, ok := bot.PickIDnTXT(upd)
	if !ok {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, "Commands for SCRIPTORIUM:")
	msg.ReplyMarkup = scriptoriumKeyboard()

	if _, err := conn.Bot.Send(msg); err != nil {
		log.Error("failed to send start message", "err", err)
		return err
	}
	return nil
}

func (s *Scriptorium) UpdateBot(conn Conn, upd tgbotapi.Update) error {
	chatID, text, ok := bot.PickIDnTXT(upd)
	if !ok {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, "")

	switch {
	case text == "TRANSCRIPTION":
		msg.Text = "Please send audio file to triscibe"
		s.waitAudio = true

	case s.waitAudio:
		scr, err := conn.Bot.GetTextOrVoice(upd)
		if err != nil {
			log.Error("Failed to get text or voice", "err", err)
			msg.Text = "Failed"
		} else {
			s.waitAudio = false
			msg.Text = scr
		}

	default:
		log.Warn("Nothing to do")
	}

	if _, err := conn.Bot.Send(msg); err != nil {
		log.Error("failed to send response", "err", err)
	}

	if upd.CallbackQuery != nil {
		if _, err := conn.Bot.Request(tgbotapi.NewCallback(upd.CallbackQuery.ID, "")); err != nil {
			log.Warn("failed to ack callback", "err", err)
		}
	}

	return nil
}
