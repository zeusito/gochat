package jose

import "github.com/zeusito/gochat/internal/models"

type IService interface {
	Parse(jwsToken string) (*models.MyClaims, error)
}
