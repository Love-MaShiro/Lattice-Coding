package persistence

import (
	"encoding/base64"
	"lattice-coding/internal/modules/provider/domain"
)

// SimpleCryptoService 简单加密服务（预留，目前仅做 base64）
type SimpleCryptoService struct {
}

func NewCryptoService() domain.CryptoService {
	return &SimpleCryptoService{}
}

func (s *SimpleCryptoService) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString([]byte(plaintext)), nil
}

func (s *SimpleCryptoService) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}
	bytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
