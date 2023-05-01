package session

import (
	"github.com/zeusito/gochat/internal/models"
)

type IRepository interface {
	Store(session *models.Session) error
	FindOneByID(id string) (*models.Session, bool)
	GetSessionCount() int
}
