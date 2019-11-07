package cmd

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/manishjagtap/secret/vault"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// var dummyVault vault.Vault

// func init() {

// 	dummyVault.filepath = "fake/path/secret.txt"
// 	dummyVault.encodingkey = "key"
// 	dummyVault.data = map[string]string{
// 		"fake-key-1": "fake-value-1",
// 		"fake-key-2": "fake-value-2",
// 	}
// }

var tempFile = "./deleteMe.txt"

func TestGetFindVaultActualImplemenation(t *testing.T) {
	expectedErr := "no such file or directory"
	mockDirectoryStatus()

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1"}
	GetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Truef(t, strings.Contains(op, expectedErr), "Expected : error message should contain %v", expectedErr)

	resetDirectoryStatus()
}

func TestGetGetValueActualImplementation(t *testing.T) {
	expectedErr := "Key does not exist in the vault"
	mockDirectoryStatus()
	mockFindVault()

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1"}
	GetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Equalf(t, op, expectedErr, "Expected : %v", expectedErr)

	resetFindVault()
	resetDirectoryStatus()
}

func TestSetSetValueActualImplementation(t *testing.T) {
	expectedErr := "open : no such file or directory"
	mockDirectoryStatus()
	mockFindVault()

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1", "fake-value-1"}
	SetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Equalf(t, op, expectedErr, "Expected : %v", expectedErr)

	resetFindVault()
	resetDirectoryStatus()
}

func TestRootGetPathSuccess(t *testing.T) {
	mockGetWorkingDirectory()
	mockDirectoryStatus()
	_, err := getPath()

	assert.NoError(t, err, "Expected : No error")

	resetDirectoryStatus()
	resetGetWorkingDirectory()
}

func TestRootGetPathFileNotFoundError(t *testing.T) {
	expectedErr := "FileNotFoundError"
	initGetWorkingDirectoryError(expectedErr)
	_, err := getPath()

	if assert.Error(t, err) {
		assert.Equalf(t, err.Error(), expectedErr, "Expected error : %v", expectedErr)
	}
	resetGetWorkingDirectory()
}

func TestGetSuccess(t *testing.T) {
	expectedOutput := "{'fake-key-1':'fake-value-1'}"
	mockDirectoryStatus()
	mockFindVault()
	mockGetValue()

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1"}
	GetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Equalf(t, op, expectedOutput, "Expected : %v", expectedOutput)

	resetGetValue()
	resetFindVault()
	resetDirectoryStatus()
}

func TestGetGetPathError(t *testing.T) {
	expectedErr := "Invalid secret file path"
	initGetWorkingDirectoryError(expectedErr)

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1"}
	GetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Equalf(t, op, expectedErr, "Expected : %v", expectedErr)

	resetGetWorkingDirectory()
}

func TestGetFindVaultError(t *testing.T) {
	expectedErr := "Can not open vault"
	mockDirectoryStatus()
	errorFindVault(expectedErr)

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1"}
	GetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Equalf(t, op, expectedErr, "Expected : %v", expectedErr)

	resetFindVault()
	resetDirectoryStatus()
}

func TestGetGetValueError(t *testing.T) {
	expectedErr := "Error while fetching key from vault"
	mockDirectoryStatus()
	mockFindVault()
	errorGetValue(expectedErr)

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1"}
	GetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Equalf(t, op, expectedErr, "Expected : %v", expectedErr)

	resetGetValue()
	resetFindVault()
	resetDirectoryStatus()
}

func TestSetSuccess(t *testing.T) {
	expectedOutput := "`fake-key-1` successfully added to vault"
	mockDirectoryStatus()
	mockFindVault()
	mockSetValue()

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1", "fake-value-1"}
	SetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Equalf(t, op, expectedOutput, "Expected : %v", expectedOutput)

	resetSetValue()
	resetFindVault()
	resetDirectoryStatus()
}

func TestSetGetPathError(t *testing.T) {
	expectedErr := "Invalid secret file path"
	initGetWorkingDirectoryError(expectedErr)

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1", "fake-value-1"}
	SetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Equalf(t, op, expectedErr, "Expected : %v", expectedErr)

	resetGetWorkingDirectory()
}

func TestSetFindVaultError(t *testing.T) {
	expectedErr := "Can not open vault"
	mockDirectoryStatus()
	errorFindVault(expectedErr)

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1", "fake-value-1"}
	SetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Equalf(t, op, expectedErr, "Expected : %v", expectedErr)

	resetFindVault()
	resetDirectoryStatus()
}

func TestSetSetValueError(t *testing.T) {
	expectedErr := "Error while storing key into the vault"
	mockDirectoryStatus()
	mockFindVault()
	errorSetValue(expectedErr)

	file, old := setStdoutToFile()

	var tempCmd *cobra.Command
	args := []string{"fake-key-1", "fake-value-1"}
	SetCmd.Run(tempCmd, args)

	op := resetStdoutAndGetFileContent(file, old)

	assert.Equalf(t, op, expectedErr, "Expected : %v", expectedErr)

	resetSetValue()
	resetFindVault()
	resetDirectoryStatus()
}

func initGetWorkingDirectoryError(errMsg string) {
	getWorkingDirectory = func() (string, error) {
		return "", errors.New(errMsg)
	}
}

func resetGetWorkingDirectory() {
	getWorkingDirectory = func() (string, error) {
		return os.Getwd()
	}
}

func mockGetWorkingDirectory() {
	getWorkingDirectory = func() (string, error) {
		return "/fake/path", nil
	}
}

func mockDirectoryStatus() {
	directoryStatus = func(filePath string) (err error) {
		err = nil
		return
	}
}

func resetDirectoryStatus() {
	directoryStatus = func(filePath string) (err error) {
		_, err = os.Stat(filePath)
		return
	}
}

func mockFindVault() {
	findVault = func(encodingkey string, filepath string) (vault.Vault, error) {
		var dummyVault vault.Vault
		return dummyVault, nil
	}
}

func errorFindVault(errMsg string) {
	findVault = func(encodingkey string, filepath string) (vault.Vault, error) {
		var dummyVault vault.Vault
		return dummyVault, errors.New(errMsg)
	}
}

func resetFindVault() {
	findVault = func(encodingkey string, filepath string) (vault.Vault, error) {
		return vault.FindVault(encodingkey, filepath)
	}
}

func mockGetValue() {
	getValue = func(secretVault vault.Vault, key string) (string, error) {
		return "fake-value-1", nil
	}
}

func errorGetValue(errMsg string) {
	getValue = func(secretVault vault.Vault, key string) (string, error) {
		return "fake-value-1", errors.New(errMsg)
	}
}

func resetGetValue() {
	getValue = func(secretVault vault.Vault, key string) (string, error) {
		return secretVault.Get(key)
	}
}

func setStdoutToFile() (*os.File, *os.File) {
	var file *os.File

	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		file, _ = os.Create(tempFile)
	} else {
		file, _ = os.OpenFile(tempFile, os.O_WRONLY, 0777)
	}

	file.Truncate(0)
	file.Seek(0, 0)

	old := os.Stdout
	os.Stdout = file

	return file, old
}

func resetStdoutAndGetFileContent(file *os.File, old *os.File) string {
	//close old file
	file.Close()
	os.Stdout = old

	//open the same file in read mode
	file, _ = os.OpenFile(tempFile, os.O_RDONLY, 0777)
	reader := bufio.NewReader(file)
	op, _, _ := reader.ReadLine()

	return string(op)
}

func mockSetValue() {
	setValue = func(secretVault vault.Vault, key string, value string) error {
		return nil
	}
}

func errorSetValue(errMsg string) {
	setValue = func(secretVault vault.Vault, key string, value string) error {
		return errors.New(errMsg)
	}
}

func resetSetValue() {
	setValue = func(secretVault vault.Vault, key string, value string) error {
		return secretVault.Set(key, value)
	}
}
