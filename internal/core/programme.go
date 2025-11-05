package core

type PrgKind uint

const (
	PRG_IDLE PrgKind = iota
	PRG_VERTEX
	PRG_SCRIPT
	PRG_ACHTUNG
	PRG_CONTROL
)

func (prg PrgKind) String() string {
	switch prg {
	case PRG_IDLE:
		return "IDLE"
	case PRG_VERTEX:
		return "VERTEX"
	case PRG_SCRIPT:
		return "SCRIPT"
	case PRG_ACHTUNG:
		return "ACHTUNG"
	case PRG_CONTROL:
		return "CONTROL"
	}

	return "UNREACHABLE"
}
