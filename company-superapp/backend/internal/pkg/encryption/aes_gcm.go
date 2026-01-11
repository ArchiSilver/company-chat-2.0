package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
)

var (
	ErrInvalidKey        = errors.New("encryption key must be 32 bytes for AES-256")
	ErrCiphertextTooShort = errors.New("ciphertext too short")
)

type EncryptionService struct {
	key []byte
}

func NewEncryptionService() (*EncryptionService, error) {
	keyHex := os.Getenv("ENCRYPTION_KEY")
	if keyHex == "" {
		return nil, errors.New("ENCRYPTION_KEY environment variable not set")
	}

	key := []byte(keyHex)
	if len(key) != 32 {
		return nil, ErrInvalidKey
	}

	return &EncryptionService{key: key}, nil
}

func NewEncryptionServiceWithKey(key []byte) (*EncryptionService, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKey
	}
	return &EncryptionService{key: key}, nil
}

func (s *EncryptionService) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (s *EncryptionService) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrCiphertextTooShort
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
