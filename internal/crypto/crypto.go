package cryptoutil

import (
	"fmt"
	"crypto/sha256"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
)

func HashSha256(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	// Convert from []byte to hex string
	return fmt.Sprintf("%x", hash.Sum(nil))
} 

func HashScrypt(s string) (string, string, error) {
	salt, err := GenerateSalt()
	if err != nil  {
		panic(err)
	}
	hash, err := scrypt.Key([]byte(s), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(hash), base64.StdEncoding.EncodeToString(salt), nil
} 

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)	
	if err != nil {
		return []byte{}, err
	}
	return salt, nil
}

func HashScryptSalt(s string, salt string) (string, error) {
	decodedSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}
	
	newHash, err := scrypt.Key([]byte(s), decodedSalt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(newHash), nil
}

func EncryptString(plaintext, password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	out := append(salt, append(nonce, cipherText...)...)
	return base64.StdEncoding.EncodeToString(out), nil
}

func DecryptString(encrypted, password string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < 16 {
		return "", fmt.Errorf("ciphertext too short")
	}

	salt := data[:16]
	data = data[16:]

	key, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short for nonce")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
