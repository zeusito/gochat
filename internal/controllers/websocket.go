package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gobwas/ws"
	"github.com/zeusito/gochat/config"
	"github.com/zeusito/gochat/internal/models"
	"github.com/zeusito/gochat/internal/services/chat"
	"github.com/zeusito/gochat/internal/services/jose"
	"go.uber.org/zap"
)

// WebSocketServer Holder for the server configurations
type WebSocketServer struct {
	log        *zap.SugaredLogger
	sc         config.ServerConfigurations
	sessionSvc chat.IService
	joseSvc    jose.IService
}

// NewWebSocketServer Initializes a new server
func NewWebSocketServer(logger *zap.SugaredLogger, serverConf config.ServerConfigurations, ss chat.IService, js jose.IService) *WebSocketServer {
	return &WebSocketServer{
		log:        logger,
		sc:         serverConf,
		sessionSvc: ss,
		joseSvc:    js,
	}
}

// Start Fires the http server
func (c *WebSocketServer) Start() {
	listeningAddr := ":" + strconv.Itoa(c.sc.Port)
	c.log.Infof("Server listening on port %s", listeningAddr)

	// Registering the websocket handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Auth filter
		claims, ok := c.authFilter(w, r)

		// If authentication failed then forbid the connection
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// If authentication succeeded then upgrade the connection
		c.upgradeConnection(w, r, claims)
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
		c.log.Fatalf("Failed to start http server. %v", err)
	}
}

func (c *WebSocketServer) upgradeConnection(w http.ResponseWriter, r *http.Request, claims *models.MyClaims) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)

	if err != nil {
		c.log.Errorf("Failed to upgrade the http connection to websocket. %v", err)
		return
	}

	c.log.Infof("New connection from %s belonging to user %s", conn.RemoteAddr().String(), claims.UserID)

	// Join the user to the chat
	c.sessionSvc.MemberJoin(*claims, conn)
}
