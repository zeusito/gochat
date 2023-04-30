package chat

import (
	"bytes"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/goccy/go-json"
	"github.com/zeusito/gochat/internal/models"
)

type Session struct {
	Id               string
	Request          *http.Request
	conn             net.Conn
	validator        *validator.Validate
	authenticated    bool
	onMessageHandler func(session *Session, msg models.MyMessage)
	onErrorHandler   func(session *Session, err error)
	onClose          func(session *Session)
}

func (s *Session) WriteMessage(msg models.MyMessage) error {
	js, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = wsutil.WriteServerText(s.conn, js)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) ReadLoop() {
	for {
		msg, op, err := wsutil.ReadClientData(s.conn)
		if err != nil {
			s.onErrorHandler(s, err)
			break
		}

		switch op {
		case ws.OpText:
			var payload models.MyMessage
			err := json.NewDecoder(bytes.NewReader(msg)).Decode(&payload)
			if err != nil {
				s.onErrorHandler(s, err)
				break
			}

			err = s.validator.Struct(payload)
			if err != nil {
				s.onErrorHandler(s, err)
				break
			}

			s.onMessageHandler(s, payload)
		case ws.OpClose:
			s.onClose(s)
		}
	}
}

func (s *Session) Close() error {
	return s.conn.Close()
}
