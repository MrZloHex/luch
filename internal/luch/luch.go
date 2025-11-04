package luch

import (
	"luch/internal/bot"
	"luch/pkg/protocol"
)

type Luch struct {
	bot *bot.Bot
	ptcl *protocol.Protocol
}

func Init(bot *bot.Bot, ptcl *protocol.Protocol) (*Luch, error) {
	luch := Luch {
		ptcl: ptcl,
		bot: bot,
	}

	return &luch, nil
}

func (luch *Luch) Run() {
}
