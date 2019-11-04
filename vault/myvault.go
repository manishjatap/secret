package vault

import (
	"bufio"
	"errors"
	"os"
	"regexp"

	"github.com/manishjagtap/secret/encrypt"
)

//Vault : struct type
type Vault struct {
	filepath    string
	encodingkey string
	data        map[string]string
}

//Get : return the key
func (v Vault) Get(key string) (string, error) {
	if value := v.data[key]; value != "" {
		return value, nil
	}
	return "", errors.New("Key does not exist in the vault")
}

//Set : save the key
func (v Vault) Set(key, value string) error {
	//save the data in the vault (for faster use)
	v.data = make(map[string]string)
	v.data[key] = value

	//enrypt the data to be stored in the file
	cipherData, err := encrypt.Encrypt(v.encodingkey, key+"="+value)

	if err != nil {
		return err
	}

	err = updateVault(cipherData, v.filepath)

	if err != nil {
		return err
	}

	return nil
}

//FindVault : method to find vault
func FindVault(encodingkey, filepath string) (Vault, error) {
	var myVault Vault
	myVault.data = make(map[string]string)

	file, err := os.Open(filepath)

	if err != nil {
		return myVault, err
	}
	defer file.Close()

	myVault.filepath = filepath
	myVault.encodingkey = encodingkey

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		plaintext, err := encrypt.Decrypt(encodingkey, sc.Text())

		if err != nil {
			return myVault, err
		}

		regx := regexp.MustCompile(`(.*?)=`)

		key := regx.FindString(plaintext)

		if len(key) == 0 {
			return myVault, errors.New("Invalid encoding-key provided")
		}

		key = key[:len(key)-1] // removing '=' operator from string

		value := plaintext[len(key)+1:]

		myVault.data[key] = value
	}

	return myVault, nil
}

//updateVault : This is for internal use
func updateVault(data, filepath string) error {
	if file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0777); err == nil {
		if _, err := file.WriteString(data + "\n"); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
