package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/zeusito/gochat/config"
	"github.com/zeusito/gochat/internal/services/chat"
	"go.uber.org/zap"
)

// WebSocketServer Holder for the server configurations
type WebSocketServer struct {
	Logger     *zap.SugaredLogger
	sc         config.ServerConfigurations
	sessionSvc chat.IService
}

// NewWebSocketServer Initializes a new server
func NewWebSocketServer(logger *zap.SugaredLogger, serverConf config.ServerConfigurations, ss chat.IService) *WebSocketServer {
	return &WebSocketServer{
		Logger:     logger,
		sc:         serverConf,
		sessionSvc: ss,
	}
}

// Start Fires the http server
func (s *WebSocketServer) Start() {
	listeningAddr := ":" + strconv.Itoa(s.sc.Port)
	s.Logger.Infof("Server listening on port %s", listeningAddr)

	// Registering the websocket handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.sessionSvc.HandleNewConnection(w, r)
	})

	// Customizing the server
	server := &http.Server{
		Addr:         listeningAddr,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	// Start the server
	err := server.ListenAndServe()
	if err != nil {
		s.Logger.Fatalf("Failed to start http server. %v", err)
	}
}
