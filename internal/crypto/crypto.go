package cryptoutil

import (
	"fmt"
	"io"
	"crypto/sha256"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

func HashSha256(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	// Convert from []byte to hex string
	return fmt.Sprintf("%x", hash.Sum(nil))
} 

func EncryptString(s string, key string) string {
	hashedKey := sha256.New()
	hashedKey.Write([]byte(key))
	bytesKey := hashedKey.Sum(nil)
	block, err := aes.NewCipher(bytesKey)
	if err != nil {
		panic(err)
	}
	aesGCM, err := cipher.NewGCM(block) // Galois Counter Mode
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, aesGCM.NonceSize()) // must be unique per encryption
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	cipherText := aesGCM.Seal(nil, nonce, []byte(s), nil)
	return hex.EncodeToString(append(nonce, cipherText...)) 
}
