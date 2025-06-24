package cryptoutil

import (
	"fmt"
	"io"
	"crypto/sha256"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
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
	salt := make([]byte, 16)
	_, err := rand.Read(salt)	
	if err != nil {
		return "", "", err
	}
	hash, err := scrypt.Key([]byte(s), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(hash), base64.StdEncoding.EncodeToString(salt), nil
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

func EncryptString(s string, key string) string {
	hashedKey, salt, err := HashScrypt(s)
	byteHashedKey, _ := base64.StdEncoding.DecodeString(hashedKey)
	byteSalt, _ := base64.StdEncoding.DecodeString(salt)
	block, err := aes.NewCipher(byteHashedKey)
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
	return hex.EncodeToString(append(byteSalt, append(nonce, cipherText...)...)) 
}

func DecryptString(s string, key string) string {
	encryptedText, err := hex.DecodeString(s)
	saltSize := 16 // Assume that all encryption is done within the program so all salts are 16 bytes
	if len(encryptedText) < saltSize {
		panic(fmt.Errorf("ciphertext shorter than salt?"))
	}
	salt := base64.StdEncoding.EncodeToString(encryptedText[:saltSize])
	hashedKey, err := HashScryptSalt(key, salt)
	if err != nil {
		panic(err)
	}
	byteHashedKey, _ := base64.StdEncoding.DecodeString(hashedKey)
	block, err := aes.NewCipher(byteHashedKey)
	if err != nil {
		panic(fmt.Errorf("failed to create AES cipher: %w", err))
	}		
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(fmt.Errorf("failed to create GCM: %w", err))
	}
	nonceSize := gcm.NonceSize()
	if len(encryptedText) < nonceSize + saltSize {
		panic(fmt.Errorf("ciphertext too short for GCM decryption"))
	}
	nonce, encryptedData := encryptedText[saltSize:saltSize + nonceSize], encryptedText[saltSize + nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		panic(fmt.Errorf("failed to decrypt: %w", err))
	}

	return string(plaintext)
}
