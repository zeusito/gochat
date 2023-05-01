package models

const (
	ErrorMessageType = "error"
	ChatMessageType  = "chat"
)

type MyMessage struct {
	Type string `json:"type" validate:"required, oneof=authorization unauthorized error"`
	Data string `json:"data" validate:"required"`
}
