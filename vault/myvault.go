package vault

import (
	"bufio"
	"errors"
	"os"
	"regexp"

	"github.com/manishjagtap/secret/cryptosec"
)

//Vault : struct type
type Vault struct {
	filepath    string
	encodingkey string
	data        map[string]string
}

var encryptData = func(key, plaintext string) (string, error) {
	return cryptosec.Encrypt(key, plaintext)
}

var decryptData = func(key, ciphertext string) (string, error) {
	return cryptosec.Decrypt(key, ciphertext)
}

var writeDataToFile = func(f *os.File, data string) (int, error) {
	return f.WriteString(data)
}

var openFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

var getNewScanner = func(f *os.File) *bufio.Scanner {
	return bufio.NewScanner(f)
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
	cipherData, err := encryptData(v.encodingkey, key+"="+value)

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

	file, err := openFile(filepath, os.O_RDONLY, 0777)

	if err != nil {
		return myVault, err
	}
	defer file.Close()

	myVault.filepath = filepath
	myVault.encodingkey = encodingkey

	sc := getNewScanner(file)
	for sc.Scan() {
		plaintext, err := decryptData(encodingkey, sc.Text())

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
	if file, err := openFile(filepath, os.O_APPEND|os.O_WRONLY, 0777); err == nil {
		if _, err := writeDataToFile(file, data+"\n"); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
