package chat

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gobwas/ws"
	"github.com/zeusito/gochat/internal/models"
	"github.com/zeusito/gochat/internal/services/jose"
	"go.uber.org/zap"
)

type DefaultService struct {
	log      *zap.SugaredLogger
	validate *validator.Validate
	joseSvc  jose.IService
}

func NewDefaultService(l *zap.SugaredLogger, v *validator.Validate, js jose.IService) *DefaultService {
	return &DefaultService{
		log:      l,
		validate: v,
		joseSvc:  js,
	}
}

func (s *DefaultService) HandleNewConnection(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)

	if err != nil {
		s.log.Errorf("Failed to upgrade the http connection to websocket. %v", err)
		return
	}

	s.log.Infof("New connection from %s", conn.RemoteAddr().String())

	// Wrap the connection in a session
	sess := &Session{
		Id:               "anon",
		Request:          r,
		conn:             conn,
		validator:        s.validate,
		authenticated:    false,
		onMessageHandler: s.HandleMessage,
		onErrorHandler:   s.HandleError,
		onClose:          s.HandleClose,
	}

	// Request the client to authenticate
	err = sess.WriteMessage(models.PredefinedUnAuthMessage)
	if err != nil {
		s.log.Errorf("Failed to write message to client. %v", err)
		sess.Close()
		return
	}

	// Start the read loop
	go sess.ReadLoop()
}

func (s *DefaultService) HandleMessage(session *Session, msg models.MyMessage) {
	if msg.Type == models.AuthorizationMessageType {
		// Ignore if already authenticated
		if session.authenticated {
			return
		}

		// Process JWT token

		session.authenticated = true
		return
	}

	// Check if session is authenticated, if not, request authentication
	if !session.authenticated {
		err := session.WriteMessage(models.PredefinedUnAuthMessage)
		if err != nil {
			s.log.Errorf("Failed to write message to client. %v", err)
			session.Close()
			return
		}
	}

}

func (s *DefaultService) HandleError(session *Session, err error) {
	s.log.Errorf("Error occurred in session %s. %v", session.Id, err)
}

func (s *DefaultService) HandleClose(session *Session) {
	s.log.Infof("Session %s closed.", session.Id)
}
