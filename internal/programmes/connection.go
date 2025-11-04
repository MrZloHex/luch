package programmes

import (
	"luch/pkg/protocol"
	"luch/internal/bot"
)

type Conn struct {
	Bot  *bot.Bot
	Ptcl *protocol.Protocol
}
