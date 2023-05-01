package models

const (
	ErrorMessageType     = "error"
	ChatMessageType      = "chat"
	SubscribeMessageType = "subscribe"
)

type ChatRequest struct {
	Type        string `json:"type" validate:"required,oneof=chat subscribe"`
	Data        string `json:"data" validate:"required,max=500"`
	Destination string `json:"destination" validate:"max=50"`
}

type ChatResponse struct {
	Type      string `json:"type"`
	Data      string `json:"data"`
	From      string `json:"from"`
	Timestamp string `json:"timestamp"`
}
