package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/zeusito/gochat/config"
	"github.com/zeusito/gochat/internal/controllers"
	"github.com/zeusito/gochat/internal/services/chat"
	"github.com/zeusito/gochat/internal/services/jose"
)

func main() {
	// Logger
	logger := config.NewLogger()
	defer config.CloseLogger(logger)

	// Configs
	configs := config.LoadConfig(logger)

	// Validator
	validate := validator.New()

	// Jose
	joseService, err := jose.NewService(logger, configs.Keys.Public)
	if err != nil {
		logger.Fatalf("Failed to create jose service. %v", err)
		return
	}

	// -- Dependency Injection --
	chatService := chat.NewDefaultService(logger, validate, joseService)

	server := controllers.NewWebSocketServer(logger, configs.Server, chatService)
	// -- End of Dependency Injection --

	// Fun begins
	server.Start()
}
