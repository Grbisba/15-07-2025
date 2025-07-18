package token

import "errors"

var (
	ErrorSignedToken      = errors.New("failed to sign token")
	ErrorParsedToken      = errors.New("failed to parse token")
	ErrorParsedClaims     = errors.New("failed to parse claims")
	ErrorParsedID         = errors.New("failed to parse ID")
	ErrorMethod           = errors.New("unexpected signing method")
	ErrorReadPublicKey    = errors.New("failed to read public key")
	ErrorParsedPublicKey  = errors.New("failed to parse public key")
	ErrorReadPrivateKey   = errors.New("failed to read private key")
	ErrorParsedPrivateKey = errors.New("failed to parse private key")
	ErrorTimeExpired      = errors.New("token expired")
)
