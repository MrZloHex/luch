package programmes

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PrgKind uint

const (
	PRG_IDLE PrgKind = iota
	PRG_VERTEX
	PRG_SCRIPT
)

type Programme struct {
	kind PrgKind
	Msg  Messanger

	vertex Vertex
	script Scriptorium
}

func (prg *Programme) Execute(upd tgbotapi.Update) error {
	switch prg.Which() {
	case PRG_IDLE:
		Idle()
	case PRG_VERTEX:
		return prg.vertex.Execute(prg.Msg, upd)
	case PRG_SCRIPT:
		return prg.script.Execute(prg.Msg, upd)
	}

	return nil
}

func (prg Programme) Which() PrgKind {
	return prg.kind
}

func (prg *Programme) Set(kind PrgKind) {
	prg.kind = kind
}

type Messanger interface {
	SendBot(tgbotapi.Chattable) (tgbotapi.Message, error)
	RequestBot(tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	GetFileBot(tgbotapi.FileConfig) (tgbotapi.File, error)
	GetTextOrVoice(tgbotapi.Update) (string, error)
	SendWS(...string) string
	GetToken() string
}
