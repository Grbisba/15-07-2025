package encryptor

type Interface interface {
	EncryptPassword(password string) (string, error)
	CompareHashAndPassword(hash, password string) error
}
