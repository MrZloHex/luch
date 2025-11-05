package core

import "luch/pkg/protocol"

type ProtoEvents struct {
	out chan<- Event
}

func NewProtoEvHandler(out chan<- Event) *ProtoEvents {
	return &ProtoEvents{out: out}
}

func (ev *ProtoEvents) EmitOut(msg *protocol.Message) {
	ev.out <- Event{
		Kind: EV_WS,
		WS: *msg,
	}
}
