package models

import (
	"bytes"
	"net"

	"github.com/go-playground/validator/v10"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/goccy/go-json"
)

type Session struct {
	Id               string
	Conn             net.Conn
	Validator        *validator.Validate
	Authenticated    bool
	Claims           MyClaims
	OnMessageHandler func(session *Session, msg MyMessage)
	OnErrorHandler   func(session *Session, err error)
	OnClose          func(session *Session)
}

func (s *Session) WriteMessage(msg MyMessage) error {
	js, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = wsutil.WriteServerText(s.Conn, js)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) ReadLoop() {
	for {
		msg, op, err := wsutil.ReadClientData(s.Conn)
		if err != nil {
			s.OnErrorHandler(s, err)
			break
		}

		switch op {
		case ws.OpText:
			var payload MyMessage
			err := json.NewDecoder(bytes.NewReader(msg)).Decode(&payload)
			if err != nil {
				s.OnErrorHandler(s, err)
				break
			}

			err = s.Validator.Struct(payload)
			if err != nil {
				s.OnErrorHandler(s, err)
				break
			}

			s.OnMessageHandler(s, payload)
		case ws.OpClose:
			s.OnClose(s)
		}
	}
}

func (s *Session) Close() error {
	return s.Conn.Close()
}
