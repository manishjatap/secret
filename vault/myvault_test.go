package vault

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/manishjagtap/secret/cryptosec"
	"github.com/stretchr/testify/assert"
)

var dummyVault Vault

func init() {
	if path, err := os.Getwd(); err == nil {
		dummyVault.filepath = path + "/secret.txt"
		dummyVault.encodingkey = "key"
		dummyVault.data = map[string]string{
			"fake-key-1": "fake-value-1",
			"fake-key-2": "fake-value-2",
		}
	} else {
		log.Fatal(err)
		os.Exit(1)
	}
}

func TestSetSuccess(t *testing.T) {
	mockOpenFile()
	mockWriteDataToFile()
	err := dummyVault.Set("fake-key-3", "fake-value-3")
	assert.NoError(t, err, "Expected : No Error")
	resetWriteDataToFileError()
	resetOpenFileError()
}

func TestSetEncryptionError(t *testing.T) {
	expectedErr := "Encryption error"
	initEncryptDataError(expectedErr)
	err := dummyVault.Set("fake-key-3", "fake-value-3")
	if assert.Error(t, err, "Expected : Error") {
		assert.Equalf(t, err.Error(), expectedErr, "Expected error message : %v", expectedErr)
	}
	resetEncryptDataError()
}

func TestSetOpenFileError(t *testing.T) {
	expectedErr := "open /fake/file/path: no such file or directory"
	initOpenFileError(expectedErr)
	err := dummyVault.Set("fake-key-3", "fake-value-3")
	if assert.Error(t, err, "Expected : Error") {
		assert.Equalf(t, err.Error(), expectedErr, "Expected error message : %v", expectedErr)
	}
	resetOpenFileError()
}

func TestSetWriteDataToFileError(t *testing.T) {
	expectedErr := "Error in writing file"
	mockOpenFile()
	initWriteDataToFileError(expectedErr)
	err := dummyVault.Set("fake-key-3", "fake-value-3")
	if assert.Error(t, err, "Expected : Error") {
		assert.Equalf(t, err.Error(), expectedErr, "Expected error message : %v", expectedErr)
	}
	resetWriteDataToFileError()
	resetOpenFileError()
}

func TestGetSuccess(t *testing.T) {
	res, err := dummyVault.Get("fake-key-1")
	if assert.NoError(t, err, "Expected : No Error") {
		assert.Equal(t, res, "fake-value-1", "Expected : fake-value-1")
	}
}

func TestGetInvalidKey(t *testing.T) {
	expectedErr := "Key does not exist in the vault"
	_, err := dummyVault.Get("fake-key-Invalid")
	if assert.Error(t, err, "Expected : Error") {
		assert.Equalf(t, err.Error(), expectedErr, "Expected Error Message : %v", expectedErr)
	}
}

func TestFindVaultSuccess(t *testing.T) {
	mockOpenFile()
	mockGetNewScanner()
	myVault, _ := FindVault("key", "/fake/file/path")
	assert.Equal(t, myVault.data["google-api-key"], "GOOGLE.COM", "Invalid key-value pair")
	resetGetNewScanner()
	resetOpenFileError()
}

func TestFindVaultInvalidDescryptionKeyError(t *testing.T) {
	expectedErr := "Invalid encoding-key provided"
	mockOpenFile()
	mockGetNewScanner()
	_, err := FindVault("fake-key", "/fake/file/path")

	if assert.Error(t, err, "Expected : Error") {
		assert.Equalf(t, err.Error(), expectedErr, "Expected error : %v", expectedErr)
	}
	resetGetNewScanner()
	resetOpenFileError()
}

func TestFindVaultDescryptionError(t *testing.T) {
	expectedErr := "Descryption error"
	mockOpenFile()
	mockGetNewScanner()
	initDecryptDataError(expectedErr)
	_, err := FindVault("fake-key", "/fake/file/path")

	if assert.Error(t, err, "Expected : Error") {
		assert.Equalf(t, err.Error(), expectedErr, "Expected error : %v", expectedErr)
	}
	resetDecryptDataError()
	resetGetNewScanner()
	resetOpenFileError()
}

func TestFindVaultOpenFileError(t *testing.T) {
	path := "/fake/file/path"
	expectedErr := "open " + path + ": no such file or directory"
	_, err := FindVault("fake-key", path)
	if assert.Error(t, err, "Expected : Error") {
		assert.Equalf(t, err.Error(), expectedErr, "Expected error : %v", expectedErr)
	}
}

// func TestOpenFileFunction(t *testing.T) {
// 	//Created temp file to test
// 	filepath, _ := os.Getwd()
// 	filepath = filepath + "/deleteMe.txt"
// 	createFile, _ := os.Create(filepath)
// 	createFile.Close()

// 	if file, err := openFile(filepath, os.O_APPEND|os.O_WRONLY, 0777); err == nil {
// 		defer file.Close()
// 		if _, err := writeDataToFile(file, "file to be deleted"); err != nil {
// 			log.Fatal(err)
// 		}
// 	} else {
// 		log.Fatal(err)
// 	}
// }

func initEncryptDataError(errMsg string) {
	//init error
	encryptData = func(key, plaintext string) (string, error) {
		return "fake-ciphertext", errors.New(errMsg)
	}
}

func resetEncryptDataError() {
	//reset error
	encryptData = func(key, plaintext string) (string, error) {
		return cryptosec.Encrypt(key, plaintext)
	}
}

func initDecryptDataError(errMsg string) {
	//init error
	decryptData = func(key, ciphertext string) (string, error) {
		return "fake-plaintext", errors.New(errMsg)
	}
}

func resetDecryptDataError() {
	//reset error
	decryptData = func(key, ciphertext string) (string, error) {
		return cryptosec.Decrypt(key, ciphertext)
	}
}

func initWriteDataToFileError(errMsg string) {
	//init error
	writeDataToFile = func(f *os.File, data string) (int, error) {
		return 0, errors.New(errMsg)
	}
}

func resetWriteDataToFileError() {
	//reset error
	writeDataToFile = func(f *os.File, data string) (int, error) {
		return f.WriteString(data)
	}
}

func mockWriteDataToFile() {
	//init error
	writeDataToFile = func(f *os.File, data string) (int, error) {
		return 0, nil
	}
}

func initOpenFileError(errMsg string) {
	//init error
	openFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return new(os.File), errors.New(errMsg)
	}
}

func resetOpenFileError() {
	//reset error
	openFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return os.OpenFile(name, flag, perm)
	}
}

func mockOpenFile() {
	openFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return new(os.File), nil
	}
}

func mockGetNewScanner() {
	getNewScanner = func(f *os.File) *bufio.Scanner {
		buf := bytes.NewBufferString("91615023759679f0b30d6f8b956cca681d7f522815626a69d1bab37bf20554d351602b775f5aa6c7cd")
		return bufio.NewScanner(buf)
	}
}

func resetGetNewScanner() {
	getNewScanner = func(f *os.File) *bufio.Scanner {
		return bufio.NewScanner(f)
	}
}
