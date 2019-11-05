package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

//RootCmd : This variable is used for secret command
var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "This command is used to handle secret vault",
}

var encodingkey string

var getWorkingDirectory = func() (string, error) {
	return os.Getwd()
}

// var directoryStatus(filePath string) (*os.FileInfo, errors) {
// 	return os.Stat(dir)
// }

var isFileNotExist = func(err error) bool {
	return os.IsNotExist(err)
}

func init() {
	// StringVarP is like StringVar, but accepts a shorthand letter that can be used after a single dash.
	// StringVarP(p *string, name, shorthand string, value string, usage string)
	RootCmd.PersistentFlags().StringVarP(&encodingkey, "encodingkey", "k", "", "this key is used for encoding and decoding secrets")
	getPath()
}

func getPath() (string, error) {
	if dir, err := getWorkingDirectory(); err == nil {
		dir = string(dir + "/secret.txt")

		if _, err := os.Stat(dir); isFileNotExist(err) {
			return "", err
		}
		return dir, err
	} else {
		return "", err
	}
}
