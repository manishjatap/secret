package cmd

import (
	"fmt"

	"github.com/manishjagtap/secret/vault"
	"github.com/spf13/cobra"
)

var getValue = func(secretVault vault.Vault, key string) (string, error) {
	return secretVault.Get(key)
}

//GetCmd : Fetch the key from a secret vault
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch the key from a secret vault",
	Long:  "usage: secret get <key-name> -k <encoding-key>",
	Run: func(c *cobra.Command, args []string) {

		key := args[0]

		if path, err := getPath(); err == nil {
			if secretVault, err := findVault(encodingkey, path); err == nil {
				if value, err := getValue(secretVault, key); err == nil {
					fmt.Printf("{'%v':'%v'}\n", key, value)
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(GetCmd)
}
