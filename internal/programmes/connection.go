package programmes

import (
	"luch/internal/bot"
	"luch/pkg/protocol"
)

type Conn struct {
	Bot  *bot.Bot
	Ptcl *protocol.Protocol
}
