package programmes

import (
	_ "strings"
	log "log/slog"
	"luch/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Achtung struct {
	cmd string
}

func achtungKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("New Timer", "NEW:TIMER"),
			tgbotapi.NewInlineKeyboardButtonData("New Alarm", "NEW:ALARM"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Get Timers", "GET:TIMER"),
			tgbotapi.NewInlineKeyboardButtonData("Get Alarms", "GET:ALARM"),
		),
	)
}

func (ach *Achtung) Start(conn Conn, upd tgbotapi.Update) error {
	chatID, _, ok := bot.PickIDnTXT(upd)
	if !ok {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, "Commands for ACHTUNG:")
	msg.ReplyMarkup = achtungKeyboard()

	if _, err := conn.Bot.Send(msg); err != nil {
		log.Error("achtung: failed to send start message", "err", err)
		return err
	}
	return nil
}

func (ach *Achtung) UpdateBot(conn Conn, upd tgbotapi.Update) error {
	chatID, text, ok := bot.PickIDnTXT(upd)
	if !ok {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, "")

	switch text {
	case "NEW:TIMER":
		msg.Text = "Send please name and duration of timer `oven 20m`"
	case "GET:TIMER":
		fallthrough
	default:
		msg.Text = "NOT IMPL"
	}

	if _, err := conn.Bot.Send(msg); err != nil {
		log.Error("vertex: failed to send response", "err", err)
	}

	if upd.CallbackQuery != nil {
		if _, err := conn.Bot.Request(tgbotapi.NewCallback(upd.CallbackQuery.ID, "")); err != nil {
			log.Warn("vertex: failed to ack callback", "err", err)
		}
	}

	return nil
}

/*
func (ach *Achtung) newTimer(m Messanger, upd tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
	tim := strings.Split(upd.Message.Text, " ")
	if len(tim) < 2 {
		msg.Text = "Please send with correct arguments"
		m.SendBot(msg)
		return nil
	}

	msg.Text = m.SendWS("ACHTUNG", "SET:TIMER", tim[0], tim[1])
	_, err := m.SendBot(msg)
	return err
}

func (ach *Achtung) Execute(m Messanger, upd tgbotapi.Update) error {
	if upd.CallbackQuery != nil {
		return ach.callback(m, upd)
	}

	switch ach.cmd {
	case "NEW:TIMER":
		return ach.newTimer(m, upd)
	default:
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
		msg.Text = "Commands for ACHTUNG:"
		msg.ReplyMarkup = achtungMakeKB()
		_, err := m.SendBot(msg)
		return err
	}
}
*/
