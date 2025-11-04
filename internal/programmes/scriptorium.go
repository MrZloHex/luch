package programmes

/*

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Note struct {
	Text string
}

type Scriptorium struct {
	cmd string
}

func (s *Scriptorium) callback(m Messanger, upd tgbotapi.Update) error {
	s.cmd = upd.CallbackData()
	msg := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "")

	switch upd.CallbackData() {
	case "NEW":
		msg.Text = "Send please audio file or text of note"
	case "OKAY":
		msg.Text = "Here should be websocket send command"
	case "NEW:EDIT":
		msg.Text = "Please send new text for note"
		s.cmd = "NEW"
	case "CANCEL":
		s.cmd = ""
	case "EDIT":
		fallthrough
	case "GET":
		msg.Text = "NOT IMPL"
	}

	_, err := m.SendBot(msg)
	m.RequestBot(tgbotapi.NewCallback(upd.CallbackQuery.ID, ""))
	return err
}

func scriptMakeKB() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("New Note", "NEW"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Edit Note", "EDIT"),
			tgbotapi.NewInlineKeyboardButtonData("Get Note", "GET"),
		),
	)
}

func (s *Scriptorium) newNote(m Messanger, upd tgbotapi.Update) error {
	txt, err := m.GetTextOrVoice(upd)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
	msg.Text = "You want to add a note with:\n`\n"+txt+"\n`"
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Okay", "OKAY"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Edit", "NEW:EDIT"),
			tgbotapi.NewInlineKeyboardButtonData("Cancel", "CANCEL"),
		),
	)
	m.SendBot(msg)

	return nil
}

func (s *Scriptorium) Execute(m Messanger, upd tgbotapi.Update) error {
	if upd.CallbackQuery != nil {
		return s.callback(m, upd)
	}

	switch s.cmd {
	case "NEW":
		s.newNote(m, upd)
	case "EDIT":
		fallthrough
	case "GET":
		fallthrough
	default:
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
		msg.Text = "Commands for SCRIPTORIUM:"
		msg.ReplyMarkup = scriptMakeKB()
		_, err := m.SendBot(msg)
		return err
	}

	return nil
}
*/
