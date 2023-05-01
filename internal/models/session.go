package models

import (
	"bytes"
	"errors"
	"io"
	"net"

	"github.com/go-playground/validator/v10"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

type ISession interface {
	WriteMessage(msg ChatResponse) error
	ReadLoop()
	Close() error
}

type Session struct {
	Id               string
	Conn             net.Conn
	Validator        *validator.Validate
	Log              *zap.SugaredLogger
	Authenticated    bool
	Claims           MyClaims
	OnMessageHandler func(session *Session, msg ChatRequest)
	OnErrorHandler   func(session *Session, err error)
	OnClose          func(session *Session)
}

func (s *Session) WriteMessage(msg ChatResponse) error {
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
	defer s.Conn.Close()

	// Infinite loop
	for {
		msg, op, err := wsutil.ReadClientData(s.Conn)
		if err != nil {
			s.Log.Errorf("Error reading from session %s. %v", s.Id, err)
			if errors.Is(err, io.EOF) {
				return
			}
			continue
		}

		// For simplicity, we only handle text messages, which, are actually expected to be JSON encoded
		if op == ws.OpText {
			var payload ChatRequest
			err := json.NewDecoder(bytes.NewReader(msg)).Decode(&payload)

			// on error, let the client know
			if err != nil {
				s.Log.Errorf("Error reading from session %s. %v", s.Id, err)
				_ = s.WriteMessage(ChatResponse{
					Type: ErrorMessageType,
					Data: "Invalid message format",
				})
				continue
			}

			err = s.Validator.Struct(payload)

			// on error, let the client know
			if err != nil {
				s.Log.Errorf("Error validating message from session %s. %v", s.Id, err)
				_ = s.WriteMessage(ChatResponse{
					Type: ErrorMessageType,
					Data: "Invalid message format",
				})
				continue
			}

			// Dispatch the message to the handler
			s.OnMessageHandler(s, payload)
		}
	}
}

func (s *Session) Close() error {
	return s.Conn.Close()
}
