package chat

import (
	"net"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zeusito/gochat/internal/models"
	"github.com/zeusito/gochat/internal/repositories/session"
	"go.uber.org/zap"
)

type DefaultService struct {
	log          *zap.SugaredLogger
	validate     *validator.Validate
	sessionStore session.IRepository
}

func NewDefaultService(l *zap.SugaredLogger, v *validator.Validate, sr session.IRepository) *DefaultService {
	return &DefaultService{
		log:          l,
		validate:     v,
		sessionStore: sr,
	}
}

func (s *DefaultService) MemberJoin(claims models.MyClaims, conn net.Conn) {
	// Create a new session
	sess := &models.Session{
		Id:               claims.UserID,
		Conn:             conn,
		Claims:           claims,
		Validator:        s.validate,
		Log:              s.log,
		Authenticated:    true,
		OnMessageHandler: s.onNewMessage,
		OnErrorHandler:   s.onError,
		OnClose:          s.onClose,
	}

	// Welcome the user
	_ = sess.WriteMessage(models.ChatResponse{
		Type: models.ChatMessageType,
		Data: "Welcome " + claims.UserName + "! You are now connected to the chat server.",
	})

	// Store the session
	_ = s.sessionStore.Store(sess)

	// Start the read loop
	go sess.ReadLoop()
}

func (s *DefaultService) onNewMessage(session *models.Session, msg models.ChatRequest) {
	s.log.Infof("New message received from session %s. %v", session.Id, msg)

	if msg.Type == models.ChatMessageType {
		if len(msg.Destination) == 0 {
			// broadcast msg
			all := s.sessionStore.FindAll()

			// All but yourself
			for _, sess := range all {
				if sess.Id != session.Id {
					_ = sess.WriteMessage(models.ChatResponse{
						Type:      models.ChatMessageType,
						Data:      msg.Data,
						From:      session.Claims.UserID,
						Timestamp: time.Now().UTC().Format(time.RFC3339),
					})
				}
			}
			return
		}

		// direct msg, avoid sending to yourself
		if msg.Destination == session.Claims.UserID {
			_ = session.WriteMessage(models.ChatResponse{
				Type: models.ErrorMessageType,
				Data: "You cannot send a direct message to yourself",
			})
			return
		}
		target, ok := s.sessionStore.FindOneByID(msg.Destination)
		if ok {
			_ = target.WriteMessage(models.ChatResponse{
				Type:      models.ChatMessageType,
				Data:      msg.Data,
				From:      session.Claims.UserID,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})
		} else {
			_ = session.WriteMessage(models.ChatResponse{
				Type: models.ErrorMessageType,
				Data: "User " + msg.Destination + " not found",
			})
		}
	}

	if msg.Type == models.SubscribeMessageType {
		_ = session.WriteMessage(models.ChatResponse{
			Type: models.ChatMessageType,
			Data: "You are now subscribed to " + msg.Destination,
		})
	}
}

func (s *DefaultService) onError(session *models.Session, err error) {
	s.log.Errorf("Error occurred in session %s. %v", session.Id, err)
}

func (s *DefaultService) onClose(session *models.Session) {
	s.log.Infof("Session %s closed.", session.Id)
}
