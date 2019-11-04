package main

import (
	"fmt"
	"os"

	"github.com/manishjagtap/secret/vault"
)

func main() {
	const encodingkey string = "key"

	path, err := os.Getwd()

	v, err := vault.FindVault(encodingkey, string(path+"/secret.txt"))

	if err != nil {
		fmt.Println(err)
	}

	if err := v.Set("twitter_api_key", "my_secret_twitter_key"); err != nil {
		fmt.Println(err)
	}

	if err := v.Set("facebook_api_key", "my_secret_facebook_key"); err != nil {
		fmt.Println(err)
	}

	if value, err := v.Get("facebook_api_key"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("facebook_api_key", value)
	}

	if value, err := v.Get("twitter_api_key"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("twitter_api_key", value)
	}

}
