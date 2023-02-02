package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func EncryptAES(key []byte, plaintext []byte) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func DecryptAES(key []byte, ciphertext []byte) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		panic(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	ret, _ := gcm.Open(nil, nonce, ciphertext, nil)
	return ret
}
