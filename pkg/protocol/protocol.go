package protocol

import (
	"fmt"
	log "log/slog"
	"time"

	ws "github.com/gorilla/websocket"
)

type Protocol struct {
	shard string
	url   string
	conn  *ws.Conn
}

func NewProtocol(name string, url string) (*Protocol, error) {
	log.Debug("init websocket protocol", "url", url)

	ptcl := Protocol{
		shard: name,
		url:   url,
	}

	conn, _, err := ws.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Error("Failed to dial url", "err", err)
		return nil, err
	}
	ptcl.conn = conn

	return &ptcl, nil
}

func (ptcl *Protocol) Run() {
	ptcl.read()
}

func (ptcl *Protocol) Write(to string, payload string) error {
	packet := []byte(fmt.Sprintf("%s:%s:%s", to, payload, ptcl.shard))
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
			if ws.IsCloseError(err,
				ws.CloseNormalClosure,
				ws.CloseGoingAway,
				ws.CloseAbnormalClosure) {
				log.Warn("Connection closed", "err", err)
				break
			} else {
				log.Error("Failed to read", "err", err)
			}
		}

		log.Info("Got msg", "msg", string(msg))
	}

	ptcl.tryReconn()
}

func (ptcl *Protocol) tryReconn() {
	log.Warn("Trying to reconnect on", "url", ptcl.url)

	for {
		conn, _, err := ws.DefaultDialer.Dial(ptcl.url, nil)
		if err == nil {
			ptcl.conn = conn
			break
		}

		log.Warn("Failed to dial url", "err", err)
		time.Sleep(time.Second * 5)
	}

	log.Info("Succefully reconnected")
	go ptcl.read()
}
