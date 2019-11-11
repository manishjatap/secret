package cryptosec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNewCipherActualImplementation(t *testing.T) {
	_, err := getNewCipher([]byte("fake-maxsize-key"))
	assert.NoError(t, err, "Not expecting any error")
}

func TestReadFullActualImplementation(t *testing.T) {
	_, err := readFull(rand.Reader, []byte("data"))
	assert.NoError(t, err, "Not expecting any error")
}

func TestGetNewCFBEncrypterActualImplementation(t *testing.T) {
	ciphertext := []byte("fake-maxsize-key")
	cipher, _ := getNewCipher([]byte("fake-maxsize-key"))
	stream := getNewCFBEncrypter(cipher, ciphertext[:aes.BlockSize])
	assert.NotNil(t, stream, "Not expecting value")
}

func TestGetNewCFBDecrypterActualImplementation(t *testing.T) {
	ciphertext := []byte("fake-maxsize-key")
	cipher, _ := getNewCipher([]byte("fake-maxsize-key"))
	stream := getNewCFBDecrypter(cipher, ciphertext[:aes.BlockSize])
	assert.NotNil(t, stream, "Not expecting nil value")
}

func TestXORKeyStreamActualImplementation(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotEqual(t, err, nil, "Expecting error")
	}()
	ciphertext := []byte("fake-maxsize-key")
	dst := ciphertext[aes.BlockSize:]
	cipher, _ := getNewCipher([]byte("fake-maxsize-key"))
	stream := getNewCFBDecrypter(cipher, ciphertext[:aes.BlockSize])
	xorKeyStream(stream, dst, ciphertext)
}

func TestEncrypt(t *testing.T) {

	mockGetNewCipher()
	mockReadFull()
	mockGetNewCFBEncrypter()
	mockXorKeyStream()

	_, err := Encrypt("fake-key", "fake-plaintext")

	assert.NoError(t, err, "Expected : No error")
}

func TestDescrypt(t *testing.T) {

	mockGetNewCipher()
	mockGetNewCFBDecrypter()
	mockXorKeyStream()

	_, err := Decrypt("fake-key", "000000000000000000000000000000000000000000000000000000000000") // Do not modify this test data

	assert.NoError(t, err, "Expected : No error")
}

func TestEncryptGetNewCipherError(t *testing.T) {
	expectedErr := "Error occured while creating new cipher"
	errorGetNewCipher(expectedErr)

	_, err := Encrypt("fake-key", "fake-plaintext")

	assert.Equalf(t, expectedErr, err.Error(), "Expected : %v", expectedErr)
}

func TestEncryptReadFullError(t *testing.T) {
	expectedErr := "Error occured in readFull function"
	mockGetNewCipher()
	errorReadFull(expectedErr)

	_, err := Encrypt("fake-key", "fake-plaintext")

	assert.Equalf(t, expectedErr, err.Error(), "Expected : %v", expectedErr)
}

func TestDecryptGetNewCipherError(t *testing.T) {
	expectedErr := "Error occured while creating new cipher"
	errorGetNewCipher(expectedErr)

	_, err := Decrypt("fake-key", "fake-plaintext")

	assert.Equalf(t, expectedErr, err.Error(), "Expected : %v", expectedErr)
}

func TestDecryptDecodeStringError(t *testing.T) {
	mockGetNewCipher()
	mockGetNewCFBDecrypter()
	mockXorKeyStream()

	_, err := Decrypt("fake-key", "Wrong-data") // Do not modify this test data

	assert.Error(t, err, "Expected : error")
}

func mockGetNewCipher() {
	getNewCipher = func(key []byte) (cipher.Block, error) {
		var temp cipher.Block
		return temp, nil
	}
}

func errorGetNewCipher(errMsg string) {
	getNewCipher = func(key []byte) (cipher.Block, error) {
		var temp cipher.Block
		return temp, errors.New(errMsg)
	}
}

func mockReadFull() {
	readFull = func(r io.Reader, buf []byte) (n int, err error) {
		return 0, nil
	}
}

func errorReadFull(errMsg string) {
	readFull = func(r io.Reader, buf []byte) (n int, err error) {
		return 0, errors.New(errMsg)
	}
}

func mockGetNewCFBEncrypter() {
	getNewCFBEncrypter = func(block cipher.Block, iv []byte) cipher.Stream {
		var temp cipher.Stream
		return temp
	}
}

func mockGetNewCFBDecrypter() {
	getNewCFBDecrypter = func(block cipher.Block, iv []byte) cipher.Stream {
		var temp cipher.Stream
		return temp
	}
}

func mockXorKeyStream() {
	xorKeyStream = func(stream cipher.Stream, dst, src []byte) {
		//do nothing
		return
	}
}
