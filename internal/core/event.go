package core

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
import "luch/pkg/protocol"

type EventKind uint

const (
	EV_BOT EventKind = iota
	EV_CTRL
	EV_WS
)

type CtrlEventKind uint

const (
	SET_PRG CtrlEventKind = iota
)

type CtrlEvent struct {
	Kind CtrlEventKind
	Prg  PrgKind
}

type Event struct {
	Kind EventKind
	Bot  tgbotapi.Update
	Ctrl CtrlEvent
	WS   protocol.Message
}
