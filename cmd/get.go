package cmd

import (
	"fmt"

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

		if path, err := getPath(); err == nil {
			if secretVault, err := vault.FindVault(encodingkey, path); err == nil {
				fmt.Println(secretVault)
				if value, err := secretVault.Get(key); err == nil {
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
