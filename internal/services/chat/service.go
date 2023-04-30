package chat

import "net/http"

type IService interface {
	HandleNewConnection(w http.ResponseWriter, r *http.Request)
}
