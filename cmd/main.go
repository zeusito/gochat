package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/zeusito/gochat/config"
	"github.com/zeusito/gochat/internal/controllers"
	"github.com/zeusito/gochat/internal/services/chat"
)

func main() {
	// Logger
	logger := config.NewLogger()
	defer config.CloseLogger(logger)

	// Configs
	configs := config.LoadConfig(logger)

	// Validator
	validate := validator.New()

	// -- Dependency Injection --
	chatService := chat.NewDefaultService(logger, validate)

	server := controllers.NewWebSocketServer(logger, configs.Server, chatService)
	// -- End of Dependency Injection --

	// Fun begins
	server.Start()
}
