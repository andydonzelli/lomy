package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
)

// EncryptionSession using an AES cipher in CFB mode under the hood.
// Both the cipher's secret key and initialisation vector are produced
// from the sha256 hash of the user provided "secret"
type EncryptionSession struct {
	cfbEncrypter cipher.Stream
	cfbDecrypter cipher.Stream
}

func NewEncryptionSession(secret string) EncryptionSession {
	// Hash the secret string to form a 32 byte block
	hashBytes := hash(secret)

	block, err := aes.NewCipher(hashBytes)
	if err != nil {
		panic(err)
	}

	initialisationBytes := hash(secret + "for initialisation vector")[0:16]

	return EncryptionSession{
		cfbEncrypter: cipher.NewCFBEncrypter(block, initialisationBytes),
		cfbDecrypter: cipher.NewCFBDecrypter(block, initialisationBytes),
	}
}

// Encrypt the provided string into a base64-encoded string
func (ec *EncryptionSession) Encrypt(text string) string {
	plainText := []byte(text)
	cipherText := make([]byte, len(plainText))
	ec.cfbEncrypter.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText)
}

// Decrypt the base64-encoded string into a new string
func (ec *EncryptionSession) Decrypt(text string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	plainText := make([]byte, len(cipherText))
	ec.cfbDecrypter.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func hash(value string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	return hasher.Sum(nil)
}
