package cmd

import (
	"github.com/spf13/cobra"
)

//RootCmd : This variable is used for secret command
var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "This command is used to handle secret vault",
}

var encodingkey string

func init() {
	// StringVarP is like StringVar, but accepts a shorthand letter that can be used after a single dash.
	// StringVarP(p *string, name, shorthand string, value string, usage string)
	RootCmd.PersistentFlags().StringVarP(&encodingkey, "encodingkey", "k", "", "this key is used for encoding and decoding secrets")
}
