package cryptosec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

var getNewCipher = func(key []byte) (cipher.Block, error) {
	return aes.NewCipher(key)
}

var readFull = func(r io.Reader, buf []byte) (n int, err error) {
	return io.ReadFull(r, buf)
}

var getNewCFBEncrypter = func(block cipher.Block, iv []byte) cipher.Stream {
	return cipher.NewCFBEncrypter(block, iv)
}

var xorKeyStream = func(stream cipher.Stream, dst, src []byte) {
	stream.XORKeyStream(dst, src)
}

var getNewCFBDecrypter = func(block cipher.Block, iv []byte) cipher.Stream {
	return cipher.NewCFBDecrypter(block, iv)
}

// Encrypt will take in a key and plaintext and return a hex representation
// of the encrypted value.
// This code is based on the standard library examples at:
//   - https://golang.org/pkg/crypto/cipher/#NewCFBEncrypter
func Encrypt(key, plaintext string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := readFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := getNewCFBEncrypter(block, iv)
	xorKeyStream(stream, ciphertext[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt will take in a key and a cipherHex (hex representation of
// the ciphertext) and decrypt it.
// This code is based on the standard library examples at:
//   - https://golang.org/pkg/crypto/cipher/#NewCFBDecrypter
func Decrypt(key, cipherHex string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("encrypt: cipher too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := getNewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	xorKeyStream(stream, ciphertext, ciphertext)
	return string(ciphertext), nil
}

func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()
	fmt.Fprint(hasher, key)
	cipherKey := hasher.Sum(nil)
	return getNewCipher(cipherKey)
}
