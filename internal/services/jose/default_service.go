package jose

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/zeusito/gochat/internal/models"
	"go.uber.org/zap"
)

// DefaultService Implements a default JWT Service (Jose)
type DefaultService struct {
	log       *zap.SugaredLogger
	publicKey *rsa.PublicKey
}

// NewService all keys must be base64 encoded
func NewService(logger *zap.SugaredLogger, publicKey string) (*DefaultService, error) {

	if len(publicKey) == 0 {
		logger.Errorf("Public key is required")
		return nil, errors.New("public key is required")
	}

	// Decodes and parses the public key
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		logger.Errorf("Failed to decode public key. %v", err)
		return nil, err
	}

	block, _ := pem.Decode(decodedPublicKey)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		logger.Errorf("Failed to parse public key. %v", err)
		return nil, err
	}

	return &DefaultService{
		publicKey: key.(*rsa.PublicKey),
		log:       logger,
	}, nil
}

// Parse Validates an existing JWS and returns its claims
func (s *DefaultService) Parse(jwsToken string) (*models.MyClaims, error) {

	parsed, err := jwt.Parse([]byte(jwsToken), jwt.WithValidate(true), jwt.WithKey(jwa.RS256, s.publicKey),
		jwt.WithIssuer("zeusito"), jwt.WithAudience("zeusito.me"), jwt.WithRequiredClaim("userId"),
		jwt.WithRequiredClaim("userName"))

	if err != nil {
		return nil, err
	}

	var userID, userName string

	if v, ok := parsed.Get("userId"); ok {
		userID = v.(string)
	}

	if v, ok := parsed.Get("userName"); ok {
		userName = v.(string)
	}

	return &models.MyClaims{
		UserID:   userID,
		UserName: userName,
	}, nil
}
