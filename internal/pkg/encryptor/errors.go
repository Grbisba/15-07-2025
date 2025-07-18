package encryptor

import "errors"

var (
	ErrorEncryptPassword = errors.New("failed to encrypt password")
	ErrorDecryptPassword = errors.New("failed to decrypt password")
)
