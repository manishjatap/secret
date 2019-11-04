package cmd

import (
	"fmt"
	"os"

	"github.com/manishjagtap/secret/vault"
	"github.com/spf13/cobra"
)

//SetCmd : Store the key into a secret vault
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Store the key into a secret vault",
	Long:  "usage: secret set <key-name> \"<key>\" -k <encoding-key>",
	Run: func(c *cobra.Command, args []string) {

		key := args[0]
		value := args[1]

		if path, err := os.Getwd(); err == nil {
			path = string(path + "/secret.txt")
			if secretVault, err := vault.FindVault(encodingkey, path); err == nil {
				if err := secretVault.Set(key, value); err == nil {
					fmt.Printf("`%v` successfully added to vault\n", key)
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
	RootCmd.AddCommand(SetCmd)
}
