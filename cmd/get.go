package cmd

import (
	"fmt"
	"os"

	"github.com/manishjagtap/secret/vault"
	"github.com/spf13/cobra"
)

//GetCmd : Fetch the key from a secret vault
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch the key from a secret vault",
	Long:  "usage: secret get <key-name> -k <encoding-key>",
	Run: func(c *cobra.Command, args []string) {

		key := args[0]

		if path, err := os.Getwd(); err == nil {
			path = string(path + "/secret.txt")
			if secretVault, err := vault.FindVault(encodingkey, path); err == nil {
				if value, err := secretVault.Get(key); err == nil {
					fmt.Printf("Key: %v, Value: %v", key, value)
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
