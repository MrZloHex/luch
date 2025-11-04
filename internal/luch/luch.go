package luch

import (
	"luch/internal/bot"
	"luch/internal/core"
	"luch/internal/programmes"
	"luch/pkg/protocol"

	log "log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type Luch struct {
	events chan core.Event

	conn programmes.Conn

	currPrg core.PrgKind
	vertex programmes.Vertex
	// script Scriptorium
	// achtung Achtung
}

func Init(bot *bot.Bot, ptcl *protocol.Protocol) (*Luch, error) {
	luch := Luch{
		conn: programmes.Conn{
			Ptcl:   ptcl,
			Bot:    bot,
		},
		events: make(chan core.Event, 1024),
	}

	return &luch, nil
}

func (luch *Luch) GetEventChan() chan core.Event {
	return luch.events
}

func (luch *Luch) Run() {
	for ev := range luch.events {
		switch ev.Kind {
		case core.EV_CTRL:
			luch.handleCtrlEvent(ev)
		case core.EV_BOT:
			luch.updateBotPrg(ev.Bot)
		case core.EV_WS:
			log.Warn("Event EV_WS")
		}
	}
}

func (luch *Luch) handleCtrlEvent(ev core.Event) {
	switch ev.Ctrl.Kind {
	case core.SET_PRG:
		luch.currPrg = ev.Ctrl.Prg
		luch.startPrg(ev.Bot)
	}
}

func (luch *Luch) startPrg(upd tgbotapi.Update) {
	switch luch.currPrg {
	case core.PRG_VERTEX:
		luch.vertex.Start(luch.conn, upd)
	}
}

func (luch *Luch) updateBotPrg(upd tgbotapi.Update) {
	switch luch.currPrg {
	case core.PRG_VERTEX:
		luch.vertex.UpdateBot(luch.conn, upd)
	}
}

