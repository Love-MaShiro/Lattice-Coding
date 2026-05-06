package crypto

type Encryptor interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type NoopEncryptor struct{}

func (e *NoopEncryptor) Encrypt(plaintext string) (string, error) {
	return plaintext, nil
}

func (e *NoopEncryptor) Decrypt(ciphertext string) (string, error) {
	return ciphertext, nil
}

func NewNoopEncryptor() Encryptor {
	return &NoopEncryptor{}
}
