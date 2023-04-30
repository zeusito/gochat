package models

const (
	AuthorizationMessageType = "authorization"
	ErrorMessageType         = "error"
	UnAuthorizedMessageType  = "unauthorized"
)

type MyMessage struct {
	Type string `json:"type" validate:"required, oneof=authorization unauthorized error"`
	Data string `json:"data" validate:"required"`
}

var PredefinedUnAuthMessage = MyMessage{
	Type: UnAuthorizedMessageType,
	Data: "please authenticate",
}
