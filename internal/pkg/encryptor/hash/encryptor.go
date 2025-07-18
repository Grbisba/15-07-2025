package hash

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/Grbisba/15-07-2025/internal/pkg/encryptor"
)

var _ encryptor.Interface = (*Encryptor)(nil)

type Option func(encryptor *Encryptor)

type Encryptor struct {
	log                        *zap.Logger
	cost                       int
	generateFromPasswordFunc   func([]byte, int) ([]byte, error)
	compareHashAndPasswordFunc func([]byte, []byte) error
}

func New(log *zap.Logger, opts ...Option) *Encryptor {
	if log == nil {
		log = zap.NewNop()
	}
	e := &Encryptor{
		log:                        log.Named("encryptor"),
		cost:                       bcrypt.DefaultCost,
		generateFromPasswordFunc:   bcrypt.GenerateFromPassword,
		compareHashAndPasswordFunc: bcrypt.CompareHashAndPassword,
	}

	e.log.Debug("created default struct for encryptor")

	for _, opt := range opts {
		opt(e)
		e.log.Debug("created option for encryptor", zap.Any("option", opt))
	}
	e.log.Info("encryptor initialized successfully")
	return e
}

func (e *Encryptor) EncryptPassword(password string) (string, error) {
	e.log.Debug("encrypting password", zap.String("password", password))
	hash, err := e.generateFromPasswordFunc([]byte(password), e.cost)
	if err != nil {
		e.log.Debug("failed to generate password", zap.Error(err))
		return "", encryptor.ErrorEncryptPassword
	}
	e.log.Debug("generated password", zap.String("hash", string(hash)))
	return string(hash), nil
}

func (e *Encryptor) CompareHashAndPassword(hash, password string) error {
	e.log.Debug("comparing hash and password", zap.String("hash", hash), zap.String("password", password))
	err := e.compareHashAndPasswordFunc([]byte(hash), []byte(password))
	if err != nil {
		e.log.Debug("failed to compare hash and password", zap.String("hash", hash), zap.String("password", password))
		return encryptor.ErrorDecryptPassword
	}
	e.log.Debug("detected password", zap.String("hash", hash), zap.String("password", password))
	return nil
}
