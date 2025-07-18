package jwt

import (
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/Grbisba/15-07-2025/internal/model/dto"
	"github.com/Grbisba/15-07-2025/internal/pkg/token"
)

func (prv *Provider) CreateAccessAndRefreshTokenForUser(userID uuid.UUID) (string, string, error) {
	access, err := prv.CreateAccessTokenForUser(userID)
	if err != nil {
		return "", "", err
	}

	refresh, err := prv.CreateRefreshTokenForUser(userID)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (prv *Provider) CreateAccessTokenForUser(userID uuid.UUID) (string, error) {
	prv.log.Debug("start creating access token")

	now := time.Now()

	add := time.Duration(prv.accessLifeTime) * time.Hour

	claims := &JWT{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(add)),
		},
		IsAccess: true,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	prv.log.Debug("created jwt with claims", zap.Any("claims", claims))

	tokenString, err := t.SignedString(prv.privateKey)

	if err != nil {
		prv.log.Debug("failed to sign access token", zap.Error(err))
		return "", token.ErrorSignedToken
	}

	prv.log.Debug("created access token", zap.String("access token", tokenString))

	return tokenString, nil
}

func (prv *Provider) CreateRefreshTokenForUser(userID uuid.UUID) (string, error) {
	prv.log.Debug("start creating refresh token")

	now := time.Now()
	add := time.Duration(prv.refreshLifeTime) * time.Hour

	claims := &JWT{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(add)),
		},
		IsAccess: false,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	prv.log.Debug("created jwt with claims", zap.Any("claims", claims))

	tokenString, err := t.SignedString(prv.privateKey)

	if err != nil {
		prv.log.Debug("failed to sign refresh token", zap.Error(err))
		return "", token.ErrorSignedToken
	}

	prv.log.Debug("created refresh token", zap.String("refresh token", tokenString))

	return tokenString, nil
}

func (prv *Provider) GetDataFromToken(raw string) (*dto.UserDataInToken, error) {
	splitToken := strings.Split(raw, "Bearer")
	if len(splitToken) != 2 {
		return nil, token.ErrorParsedToken
	}

	parsedToken, err := jwt.ParseWithClaims(strings.TrimSpace(splitToken[1]), &JWT{}, prv.readKeyFunc)
	if err != nil {
		prv.log.Debug("failed to parse jwt token", zap.Error(err))
		return nil, token.ErrorParsedToken
	}

	claims, ok := parsedToken.Claims.(*JWT)
	if !ok {
		prv.log.Debug("failed to parse jwt token claims")
		return nil, token.ErrorParsedClaims
	}

	if claims.ExpiresAt.Before(time.Now()) {
		prv.log.Debug("expired jwt token", zap.Any("claims", claims))
		return nil, token.ErrorTimeExpired
	}

	parsedID, err := uuid.Parse(claims.RegisteredClaims.Issuer)
	if err != nil {
		prv.log.Debug("failed to parse jwt token claim", zap.Error(err))
		return nil, token.ErrorParsedID
	}

	data := &dto.UserDataInToken{
		ID:        parsedID,
		IsAccess:  claims.IsAccess,
		ExpiresAt: claims.ExpiresAt.Time,
		NotBefore: claims.NotBefore.Time,
		IssuedAt:  claims.IssuedAt.Time,
	}

	prv.log.Debug("successfully parsed data", zap.Any("data", data))
	return data, nil
}

func (prv *Provider) readKeyFunc(t *jwt.Token) (interface{}, error) {
	switch t.Method.(type) {
	case *jwt.SigningMethodRSA:
		prv.log.Debug("successful read key")
		return prv.publicKey, nil
	default:
		prv.log.Debug("unsupported signing method")
		return nil, token.ErrorMethod
	}
}
