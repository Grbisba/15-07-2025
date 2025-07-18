package jwt

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/Grbisba/15-07-2025/internal/config"
	"github.com/Grbisba/15-07-2025/internal/pkg/token"
)

var _ token.Provider = (*Provider)(nil)

type JWT struct {
	jwt.RegisteredClaims
	IsAccess bool `json:"is_access"`
}

type Provider struct {
	log             *zap.Logger
	publicKey       *rsa.PublicKey
	privateKey      *rsa.PrivateKey
	accessLifeTime  int
	refreshLifeTime int
}

func NewProvider(cfg *config.Config, log *zap.Logger) (*Provider, error) {
	if log == nil {
		log = zap.NewNop()
	}
	log.Named("token")

	publicKeyRaw, err := os.ReadFile(cfg.JWT.PublicKeyPath)
	if err != nil {
		log.Debug("failed to read jwt public key", zap.Error(err))
		return nil, token.ErrorReadPublicKey
	}

	log.Debug("successful read public key path")

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyRaw)
	if err != nil {
		log.Debug("failed to parse jwt public key", zap.Error(err))
		return nil, token.ErrorParsedPublicKey
	}
	log.Debug("successful parse public key")

	privateKeyRaw, err := os.ReadFile(cfg.JWT.PrivateKeyPath)
	if err != nil {
		log.Debug("failed to read jwt private key", zap.Error(err))
		return nil, token.ErrorReadPrivateKey
	}

	log.Debug("successful read private key path")

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyRaw)
	if err != nil {
		log.Debug("failed to parse jwt private key", zap.Error(err))
		return nil, token.ErrorParsedPrivateKey
	}

	log.Debug("successful parse private key")
	provider := &Provider{
		log:             log,
		publicKey:       publicKey,
		privateKey:      privateKey,
		accessLifeTime:  cfg.JWT.AccessTokenLifeTime,
		refreshLifeTime: cfg.JWT.RefreshTokenLifeTime,
	}
	log.Debug("successful create new jwt provider")
	return provider, nil
}
