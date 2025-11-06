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
	waitName  bool
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
			tgbotapi.NewInlineKeyboardButtonData("Get List", "GET:LIST"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Get Info", "GET:INFO"),
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
	msg := fmt.Sprintf("Fired %s __%s__", upd.Noun, upd.Args[0])
	cmd := fmt.Sprintf("STOP:%s", upd.Noun)
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Stop", cmd),
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
	} else if ach.waitAlarm {
		msg.Text, err = ach.newAlarm(conn, text)
		ach.waitAlarm = false
	} else if ach.waitName {
		msg.Text, err = ach.getJob(conn, text)
		ach.waitName = false
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
		case "STOP:ALARM":
			log.Info("Stopping alarm", "name", ach.int_name)
			msg.Text = "Stoping alarm"
			ach.interrupt = false
			// TODO: maybe check if resp is ok
			_, err := conn.Ptcl.TransmitReceive([]string{"ACHTUNG", text, ach.int_name})
			if err != nil {
				log.Warn("Failed to stop alarm", "err", err)
				msg.Text = "Failed to stop alarm"
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
		case "NEW:ALARM":
			msg.Text = "Send please name and datetime `meetup 2025.11.21 10.00`"
			ach.waitAlarm = true
		case "GET:INFO":
			msg.Text = "Send please name of timer or alarm"
			ach.waitName = true
		case "GET:LIST":
			rx, err := conn.Ptcl.TransmitReceive([]string{"ACHTUNG", text})
			if err != nil {
				log.Error("Failed to get list", "err", err)
			}
			for i, arg := range rx.Args {
				if i % 2 == 0 {
					msg.Text += fmt.Sprintf("%s: ", arg)
				} else {
					msg.Text += fmt.Sprintf("%s\n", arg)
				}
			}
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

func (ach *Achtung) newAlarm(conn Conn, text string) (string, error) {
	tim := strings.Split(text, " ")
	if len(tim) < 3 {
		return "Please send with correct arguments", nil
	}

	rx, err := conn.Ptcl.TransmitReceive([]string{"ACHTUNG:NEW:ALARM", tim[0], tim[1], tim[2]})
	if err != nil {
		log.Warn("Failed to TxRx to achtung", "err", err)
		return "Failed to txrx to achtung", err
	}

	log.Info("Set up new alarm", "name", tim[0], "date", tim[1], "time", tim[2])
	return rx.String(), nil
}

func (ach *Achtung) getJob(conn Conn, text string) (string, error) {
	rx, err := conn.Ptcl.TransmitReceive([]string{"ACHTUNG:GET:JOB", text})
	if err != nil {
		log.Warn("Failed to TxRx to achtung", "err", err)
		return "Failed to txrx to achtung", err
	}

	log.Info("JOB: ", "name", text, "time", rx.Args[0])
	return rx.String(), nil


}

