package session

import (
	"sync"

	"github.com/zeusito/gochat/internal/models"
	"go.uber.org/zap"
)

type InMemoryRepository struct {
	log      *zap.SugaredLogger
	sessions map[string]*models.Session
	rwMutex  *sync.RWMutex
}

func NewInMemoryRepository(l *zap.SugaredLogger) *InMemoryRepository {
	return &InMemoryRepository{
		log:      l,
		sessions: make(map[string]*models.Session),
		rwMutex:  &sync.RWMutex{},
	}
}

func (r *InMemoryRepository) Store(session *models.Session) error {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	r.sessions[session.Id] = session

	return nil
}

func (r *InMemoryRepository) FindOneByID(id string) (*models.Session, bool) {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	if sess, ok := r.sessions[id]; ok {
		return sess, true
	}

	return nil, false
}

func (r *InMemoryRepository) GetSessionCount() int {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	return len(r.sessions)
}

func (r *InMemoryRepository) FindAll() []*models.Session {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	sessions := make([]*models.Session, 0, len(r.sessions))

	for _, sess := range r.sessions {
		sessions = append(sessions, sess)
	}

	return sessions
}
