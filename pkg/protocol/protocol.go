package protocol

import (
	"fmt"
	log "log/slog"
	"strings"
	"time"

	ws "github.com/gorilla/websocket"
)

type PtclConfig struct {
	Shard  string
	Url    string
	Reconn uint
}

type Protocol struct {
	cfg  PtclConfig
	conn *ws.Conn

	resp chan []byte

	onDisconnect func()
	onConnect    func()
}

func NewProtocol(cfg PtclConfig) (*Protocol, error) {
	log.Debug("init websocket protocol", "url", cfg.Url)

	ptcl := Protocol{
		cfg:  cfg,
		resp: make(chan []byte),
	}

	conn, _, err := ws.DefaultDialer.Dial(cfg.Url, nil)
	if err != nil {
		log.Error("Failed to dial url", "err", err)
		return nil, err
	}
	ptcl.conn = conn

	return &ptcl, nil
}

func (ptcl *Protocol) Run() {
	for {
		ptcl.read()
		ptcl.tryReconn()
	}
}

func (ptcl *Protocol) OnDisconnect(f func()) { ptcl.onDisconnect = f }
func (ptcl *Protocol) OnConnect(f func())    { ptcl.onConnect = f }

func (ptcl *Protocol) Send(parts ...string) ([]byte, error) {
	pay := strings.Join(parts, ":")
	err := ptcl.write(fmt.Sprintf("%s:%s", pay, ptcl.cfg.Shard))
	if err != nil {
		return nil, err
	}
	return <-ptcl.resp, nil
}

func (ptcl *Protocol) write(payload string) error {
	packet := []byte(payload)
	err := ptcl.conn.WriteMessage(ws.TextMessage, packet)
	if err != nil {
		log.Error("Failed to write", "err", err)
	}
	return err
}

func (ptcl *Protocol) read() {
	for {
		_, msg, err := ptcl.conn.ReadMessage()
		if err != nil {
			if WsIsClosed(err) {
				log.Warn("Connection closed", "err", err)
				break
			} else {
				log.Error("Failed to read", "err", err)
				continue
			}
		}

		log.Info("Got msg", "msg", string(msg))
		ptcl.resp <- msg
	}

	if ptcl.onDisconnect != nil {
		ptcl.onDisconnect()
	}
}

func (ptcl *Protocol) tryReconn() {
	log.Warn("Trying to reconnect on", "url", ptcl.cfg.Url)

	for {
		conn, _, err := ws.DefaultDialer.Dial(ptcl.cfg.Url, nil)
		if err == nil {
			ptcl.conn = conn
			break
		}

		log.Warn("Failed to dial url", "err", err)
		time.Sleep(time.Second * time.Duration(ptcl.cfg.Reconn))
	}

	if ptcl.onConnect != nil {
		ptcl.onConnect()
	}
	log.Info("Succefully reconnected")
}
