package main

import (
	"github.com/zeusito/gochat/config"
	"github.com/zeusito/gochat/internal/controllers"
)

func main() {
	// Logger
	logger := config.NewLogger()
	defer config.CloseLogger(logger)

	// Configs
	configs := config.LoadConfig(logger)

	// Server Router
	server := controllers.NewWebSocketServer(logger, configs.Server)

	// -- Dependency Injection --
	// -- End of Dependency Injection --

	// Fun begins
	server.Start()
}
