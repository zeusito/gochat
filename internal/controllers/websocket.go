package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/zeusito/gochat/config"
	"go.uber.org/zap"
)

// WebSocketServer Holder for the server configurations
type WebSocketServer struct {
	Logger *zap.SugaredLogger
	sc     config.ServerConfigurations
}

// NewWebSocketServer Initializes a new server
func NewWebSocketServer(logger *zap.SugaredLogger, serverConf config.ServerConfigurations) *WebSocketServer {
	return &WebSocketServer{
		Logger: logger,
		sc:     serverConf,
	}
}

// Start Fires the http server
func (s *WebSocketServer) Start() {
	listeningAddr := ":" + strconv.Itoa(s.sc.Port)
	s.Logger.Infof("Server listening on port %s", listeningAddr)

	// Registering the websocket handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)

		if err != nil {
			s.Logger.Errorf("Failed to upgrade the http connection to websocket. %v", err)
			return
		}

		// background thread to read messages from the websocket
		go func() {
			defer conn.Close()

			s.Logger.Info("New websocket connection established")

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					s.Logger.Warnf("Failed to read message from websocket. %v", err)
					return
				}

				s.Logger.Infof("Message: %s, Op: %d", string(msg), op)
			}
		}()
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
