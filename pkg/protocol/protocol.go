package protocol

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Kind string

const (
	KindCMD Kind = "CMD"
	KindOK  Kind = "OK"
	KindERR Kind = "ERR"
)

type ErrorCode string

const (
	ErrEARG    ErrorCode = "EARG"
	ErrEUNK    ErrorCode = "EUNK"
	ErrEBUSY   ErrorCode = "EBUSY"
	ErrEPERM   ErrorCode = "EPERM"
	ErrENOTSUP ErrorCode = "ENOTSUP"
	ErrETIME   ErrorCode = "ETIME"
	ErrEINTR   ErrorCode = "EINTR"
)

type Message struct {
	To   string
	ID   uint64
	Kind Kind
	Noun string
	Verb string
	Args []string
	From string
}

func (m *Message) String() string {
	parts := make([]string, 0, 6+len(m.Args))
	parts = append(parts, m.To)
	parts = append(parts, strconv.FormatUint(m.ID, 10))
	parts = append(parts, strings.ToUpper(string(m.Kind)))
	parts = append(parts, strings.ToUpper(m.Noun))
	parts = append(parts, strings.ToUpper(m.Verb))
	parts = append(parts, m.Args...)
	parts = append(parts, m.From)
	return strings.Join(parts, ":")
}

func (m *Message) Clone() *Message {
	out := *m
	if len(m.Args) > 0 {
		out.Args = append([]string(nil), m.Args...)
	}
	return &out
}

func Parse(line string) (*Message, error) {
	s := strings.TrimSpace(line)
	if s == "" {
		return nil, errors.New("empty message")
	}
	if strings.ContainsAny(s, " \t\r\n") {
		// spaces not allowed (UART/WebSocket frames are single-line)
		return nil, fmt.Errorf("invalid whitespace present")
	}
	parts := strings.Split(s, ":")
	if len(parts) < 6 {
		return nil, fmt.Errorf("too few fields: got %d, want >= 6", len(parts))
	}

	to := parts[0]
	idStr := parts[1]
	kindStr := parts[2]
	noun := parts[3]
	verb := parts[4]
	from := parts[len(parts)-1]
	args := append([]string(nil), parts[5:len(parts)-1]...)

	if !isToken(to) && !isHexID(to) && to != "ALL" {
		return nil, fmt.Errorf("invalid TO token: %q", to)
	}
	if !isToken(from) && !isHexID(from) {
		return nil, fmt.Errorf("invalid FROM token: %q", from)
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ID (decimal expected): %w", err)
	}

	kind := Kind(strings.ToUpper(kindStr))
	switch kind {
	case KindCMD, KindOK, KindERR:
	default:
		return nil, fmt.Errorf("invalid KIND: %q", kindStr)
	}

	if !isToken(noun) || !isToken(verb) {
		return nil, fmt.Errorf("invalid NOUN/VERB: %q %q", noun, verb)
	}
	for i, a := range args {
		if !isToken(a) {
			return nil, fmt.Errorf("invalid ARG[%d]: %q", i, a)
		}
	}

	msg := &Message{
		To:   to,
		ID:   id,
		Kind: kind,
		Noun: strings.ToUpper(noun),
		Verb: strings.ToUpper(verb),
		Args: args,
		From: from,
	}
	return msg, nil
}

func (m *Message) Validate() error {
	if m.To == "CONCENTRATOR" {
		return fmt.Errorf("CONCENTRATOR is non-addressable")
	}
	if m.Kind == KindCMD {
		if spec, ok := lookupSpec(m.Noun, m.Verb); ok {
			if err := spec.check(m.Args); err != nil {
				return err
			}
		}
	}
	return nil
}

func NextID() uint64 {
	return atomic.AddUint64(&idCounter, 1)
}

var idCounter uint64 = uint64(time.Now().Unix() % 1_000_000)

func BuildCMD(to, from, noun, verb string, args ...string) *Message {
	return &Message{
		To:   to,
		ID:   NextID(),
		Kind: KindCMD,
		Noun: strings.ToUpper(noun),
		Verb: strings.ToUpper(verb),
		Args: args,
		From: from,
	}
}

func BuildResp(kind Kind, req *Message, args ...string) *Message {
	if kind != KindOK && kind != KindERR {
		panic("BuildResp: kind must be OK or ERR")
	}
	return &Message{
		To:   req.From,
		ID:   req.ID,
		Kind: kind,
		Noun: req.Noun,
		Verb: req.Verb,
		Args: args,
		From: req.To,
	}
}

// Helpers & validation

var (
	tokenRe   = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)
	hexIDRe   = regexp.MustCompile(`^[0-9A-F]{2}$`)
	digitsRe  = regexp.MustCompile(`^[0-9]+$`)
	isoBasicR = regexp.MustCompile(`^\d{8}T\d{6}([+-]\d{4})$`)
)

func isToken(s string) bool {
	return tokenRe.MatchString(s)
}

func isHexID(s string) bool {
	return hexIDRe.MatchString(strings.ToUpper(s))
}

// Time helpers (for ACHTUNG ALARM:NEW <when>)

// ParseWhen parses either epoch seconds or ISO Basic (YYYYMMDDThhmmss±HHMM).
func ParseWhen(tok string) (time.Time, error) {
	if digitsRe.MatchString(tok) {
		sec, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(sec, 0).UTC(), nil
	}
	if !isoBasicR.MatchString(tok) {
		return time.Time{}, fmt.Errorf("invalid when token: %q", tok)
	}
	// Go layout for 20060102T150405-0700
	const layout = "20060102T150405-0700"
	t, err := time.Parse(layout, tok)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// FormatISOBasic formats a time as YYYYMMDDThhmmss±HHMM.
func FormatISOBasic(t time.Time) string {
	return t.Format("20060102T150405-0700")
}

// === Command specs (strict validation for known shards) ===

type opSpec struct {
	min int
	max int
	v   []argValidator
}

type argValidator func(string) error

func (s opSpec) check(args []string) error {
	if len(args) < s.min || len(args) > s.max {
		return fmt.Errorf("invalid arity: got %d, want %d..%d", len(args), s.min, s.max)
	}
	for i := 0; i < len(args) && i < len(s.v); i++ {
		if s.v[i] == nil {
			continue
		}
		if err := s.v[i](args[i]); err != nil {
			return fmt.Errorf("arg[%d]: %w", i, err)
		}
	}
	return nil
}

func lookupSpec(noun, verb string) (opSpec, bool) {
	key := strings.ToUpper(noun) + ":" + strings.ToUpper(verb)
	s, ok := specs[key]
	return s, ok
}

var specs = map[string]opSpec{
	// VERTEX
	"LAMP:ON":        {min: 0, max: 0},
	"LAMP:OFF":       {min: 0, max: 0},
	"LED:ON":         {min: 0, max: 0},
	"LED:OFF":        {min: 0, max: 0},
	"LED:BRIGHTNESS": {min: 1, max: 1, v: []argValidator{isIntRange(0, 100)}},
	"LED:EFFECT":     {min: 1, max: 3, v: []argValidator{oneOf("BLINK", "FADE", "SOLID"), isPosInt(), isIntRange(0, 100)}},
	"BUZZER:ON":      {min: 0, max: 3, v: []argValidator{isPosInt(), isPosInt(), isIntRange(0, 100)}},
	"BUZZER:OFF":     {min: 0, max: 0},

	// ACHTUNG — timers
	"TIMER:NEW":    {min: 1, max: 2, v: []argValidator{isDurationToken(), nil}},
	"TIMER:DELETE": {min: 1, max: 1},
	"TIMER:LIST":   {min: 0, max: 0},
	"TIMER:PAUSE":  {min: 1, max: 1},
	"TIMER:RESUME": {min: 1, max: 1},

	// ACHTUNG — alarms
	"ALARM:NEW":    {min: 1, max: 2, v: []argValidator{isWhenToken(), nil}},
	"ALARM:DELETE": {min: 1, max: 1},
	"ALARM:LIST":   {min: 0, max: 0},
	"ALARM:PAUSE":  {min: 1, max: 1},
	"ALARM:RESUME": {min: 1, max: 1},
}

// === ID scoping helpers ======================================================

// ComposeGlobalID creates a bus-wide unique decimal ID by embedding the
// caller's shard hex ID (e.g., "A1" => 161) into the high digits and mixing a
// local counter into the low digits:
//
//	ID = shardDec * 1e12 + (counter % 1e12)
//
// This keeps IDs decimal-only while virtually eliminating cross-device
// collisions.
func ComposeGlobalID(shardHex string, counter uint64) (uint64, error) {
	v, err := strconv.ParseUint(strings.ToUpper(shardHex), 16, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid shard hex id: %q", shardHex)
	}
	const mod = 1_000_000_000_000 // 1e12
	return v*mod + (counter % mod), nil
}

var scopedCtr uint64

// NextGlobalID returns a composed ID using a process-local counter. Pass your
// own shard hex (e.g., "A1"). Safe for concurrent use.
func NextGlobalID(shardHex string) (uint64, error) {
	c := atomic.AddUint64(&scopedCtr, 1)
	return ComposeGlobalID(shardHex, uint64(time.Now().UnixNano())+c)
}

// NextEpochScopedID builds ID = epoch_ms*100 + (shardDec % 100).
// Requires synchronized clocks but is simple and compact. Risk of collision is
// negligible if devices don't emit >100 commands in the same millisecond.
func NextEpochScopedID(shardHex string) (uint64, error) {
	v, err := strconv.ParseUint(strings.ToUpper(shardHex), 16, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid shard hex id: %q", shardHex)
	}
	ms := uint64(time.Now().UnixMilli())
	return ms*100 + (v % 100), nil
}

// Validators

func isIntRange(min, max int) argValidator {
	return func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("not int: %v", err)
		}
		if n < min || n > max {
			return fmt.Errorf("out of range %d..%d", min, max)
		}
		return nil
	}
}

func isPosInt() argValidator {
	return func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("not int: %v", err)
		}
		if n <= 0 {
			return fmt.Errorf("must be > 0")
		}
		return nil
	}
}

func oneOf(opts ...string) argValidator {
	set := make(map[string]struct{}, len(opts))
	for _, o := range opts {
		set[o] = struct{}{}
	}
	return func(s string) error {
		s = strings.ToUpper(s)
		if _, ok := set[s]; !ok {
			return fmt.Errorf("must be one of %v", opts)
		}
		return nil
	}
}

var durRe = regexp.MustCompile(`^(?:\d+ms|\d+s|\d+m|\d+h|\d+d|\d+h\d+m)$`)

func isDurationToken() argValidator {
	return func(s string) error {
		if !durRe.MatchString(s) {
			return fmt.Errorf("invalid duration token")
		}
		return nil
	}
}

func isWhenToken() argValidator {
	return func(s string) error {
		if digitsRe.MatchString(s) {
			return nil // epoch seconds
		}
		if isoBasicR.MatchString(s) {
			return nil
		}
		return fmt.Errorf("invalid when token (epoch seconds or 20060102T150405±HHMM)")
	}
}

// Example:
//
//  msg := BuildCMD("VERTEX", "LUCH", "LED", "BRIGHTNESS", "70")
//  line := msg.String() // "VERTEX:<id>:CMD:LED:BRIGHTNESS:70:LUCH"
//
//  parsed, err := Parse(line)
//  if err == nil {
//      _ = parsed.Validate(true)
//  }
