package session

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeusito/gochat/config"
	"github.com/zeusito/gochat/internal/models"
)

var logger = config.NewLogger()

func TestInMemoryRepository_Store(t *testing.T) {
	repo := NewInMemoryRepository(logger)

	for i := 0; i < 10; i++ {
		_ = repo.Store(&models.Session{
			Id: strconv.Itoa(i),
		})
	}

	assert.Equal(t, 10, repo.GetSessionCount(), "should be equal")

	// Pick one
	numberSeven, ok := repo.FindOneByID("7")
	assert.True(t, ok, "should be true")
	assert.Equal(t, "7", numberSeven.Id, "should be equal")

	notHere, ok := repo.FindOneByID("235")
	assert.False(t, ok, "should be false")
	assert.Nil(t, notHere, "should be nil")
}
