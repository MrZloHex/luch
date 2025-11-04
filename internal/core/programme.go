package core

type PrgKind uint

const (
	PRG_IDLE PrgKind = iota
	PRG_VERTEX
	PRG_SCRIPT
	PRG_ACHTUNG
)
