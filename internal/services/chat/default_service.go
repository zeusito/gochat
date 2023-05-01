package chat

import (
	"net"

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
		Authenticated:    true,
		OnMessageHandler: s.onNewMessage,
		OnErrorHandler:   s.onError,
		OnClose:          s.onClose,
	}

	// Welcome the user
	_ = sess.WriteMessage(models.MyMessage{
		Type: models.ChatMessageType,
		Data: "Welcome " + claims.UserName + "! You are now connected to the chat server.",
	})

	// Store the session
	_ = s.sessionStore.Store(sess)

	// Start the read loop
	go sess.ReadLoop()
}

func (s *DefaultService) onNewMessage(session *models.Session, msg models.MyMessage) {
	s.log.Infof("New message received from session %s. %v", session.Id, msg)
}

func (s *DefaultService) onError(session *models.Session, err error) {
	s.log.Errorf("Error occurred in session %s. %v", session.Id, err)
}

func (s *DefaultService) onClose(session *models.Session) {
	s.log.Infof("Session %s closed.", session.Id)
}
