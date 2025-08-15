package protocol

import (
	"fmt"
	log "log/slog"
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
}

func NewProtocol(cfg PtclConfig) (*Protocol, error) {
	log.Debug("init websocket protocol", "url", cfg.Url)

	ptcl := Protocol{
		cfg: cfg,
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

func (ptcl *Protocol) Write(to string, payload string) error {
	packet := []byte(fmt.Sprintf("%s:%s:%s", to, payload, ptcl.cfg.Shard))
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
			}
		}

		log.Info("Got msg", "msg", string(msg))
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

	log.Info("Succefully reconnected")
}
