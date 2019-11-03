package main

import (
	"fmt"
	"log"

	"github.com/manishj/secret/encrypt"
)

func main() {
	fmt.Println("hello")
	ciphertext, err := encrypt.Encrypt("test", "hello")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ciphertext)

	plaintext, err := encrypt.Decrypt("test", ciphertext)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(plaintext)

}
