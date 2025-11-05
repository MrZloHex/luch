package programmes

import (
	"fmt"
	log "log/slog"
	"luch/internal/bot"
	"luch/pkg/protocol"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Achtung struct {
	cmd       string
	waitTimer bool
	waitAlarm bool
	interrupt bool
	int_name  string
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

func (ach *Achtung) StartIT(conn Conn, upd protocol.Message) error {
	msg := fmt.Sprintf("Fired %s", upd.Noun)
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Stop", "STOP:TIMER"),
		),
	)

	conn.Bot.NotifyAllWithMarkup(msg, kb)
	ach.interrupt = true
	ach.int_name = upd.Args[0]

	return nil
}

func (ach *Achtung) UpdateBot(conn Conn, upd tgbotapi.Update) error {
	chatID, text, ok := bot.PickIDnTXT(upd)
	if !ok {
		return nil
	}

	msg := tgbotapi.NewMessage(chatID, "")
	var err error

	if ach.waitTimer {
		msg.Text, err = ach.newTimer(conn, text)
		ach.waitTimer = false
	} else if ach.interrupt {
		switch text {
		case "STOP:TIMER":
			log.Info("Stopping timer", "name", ach.int_name)
			msg.Text = "Stoping timer"
			ach.interrupt = false
			// TODO: maybe check if resp is ok
			_, err := conn.Ptcl.TransmitReceive([]string{"ACHTUNG", text, ach.int_name})
			if err != nil {
				log.Warn("Failed to stop timer", "err", err)
				msg.Text = "Failed to stop timer"
				ach.interrupt = true
			}
		default:
			log.Error("NOT POSSIBLE")
		}
	} else {
		switch text {
		case "NEW:TIMER":
			msg.Text = "Send please name and duration of timer `oven 20m`"
			ach.waitTimer = true
		default:
			msg.Text = "NOT IMPL"
		}
	}

	if err != nil {
		log.Error("Failed to process request")
		msg.Text = "Failed"
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

func (ach *Achtung) newTimer(conn Conn, text string) (string, error) {
	tim := strings.Split(text, " ")
	if len(tim) < 2 {
		return "Please send with correct arguments", nil
	}

	rx, err := conn.Ptcl.TransmitReceive([]string{"ACHTUNG:NEW:TIMER", tim[0], tim[1]})
	if err != nil {
		log.Warn("Failed to TxRx to achtung", "err", err)
		return "Failed to txrx to achtung", err
	}

	log.Info("Set up new timer", "name", tim[0], "dur", tim[1])
	return rx.String(), nil
}

/*

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
