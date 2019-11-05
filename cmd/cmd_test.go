package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootGetPathSuccess(t *testing.T) {
	mockGetWorkingDirectory()
	mockIsFileNotExist(false)
	_, err := getPath()

	assert.NoError(t, err, "Expected : No error")

	resetIsFileNotExist()
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

func TestGet(t *testing.T) {
	var tempCmd *cobra.Command
	var args []string

	get(tempCmd, args)
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

func mockIsFileNotExist(res bool) {
	isFileNotExist = func(err error) bool {
		return res
	}
}

func resetIsFileNotExist() {
	isFileNotExist = func(err error) bool {
		return os.IsNotExist(err)
	}
}
