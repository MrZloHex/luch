package programmes

/*


import (
	"strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Achtung struct {
	cmd string
}

func (ach *Achtung) callback(m Messanger, upd tgbotapi.Update) error {
	ach.cmd = upd.CallbackData()
	msg := tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID, "")

	switch upd.CallbackData() {
	case "NEW:TIMER":
		msg.Text = "Send please name and duration of timer `oven 20m`"
	case "GET:TIMER":
		fallthrough
	default:
		msg.Text = "NOT IMPL"
	}

	_, err := m.SendBot(msg)
	m.RequestBot(tgbotapi.NewCallback(upd.CallbackQuery.ID, ""))
	return err
}

func achtungMakeKB() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("New Timer", "NEW:TIMER"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Get Timers", "GET:TIMER"),
		),
	)
}


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
