package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/brandur/simplebox"
)

// Load loads secrets from the file
func Load(password *[simplebox.KeySize]byte) ([]Secret, error) {
	p, err := GetFilename()
	if err != nil {
		return nil, err
	}

	var secrets []Secret
	raw, err := ioutil.ReadFile(p)
	if err != nil {
		return secrets, nil
	}

	sb := simplebox.NewFromSecretKey(password)
	raw, err = sb.Decrypt(raw)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(raw, &secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}

// Save saves secrets to the file.
func Save(secrets []Secret, password *[simplebox.KeySize]byte) error {
	p, err := GetFilename()
	if err != nil {
		return err
	}

	data, _ := json.Marshal(secrets)
	err = os.MkdirAll(filepath.Dir(p), 0700)
	if err != nil {
		return err
	}

	sb := simplebox.NewFromSecretKey(password)
	data = sb.Encrypt(data)

	return ioutil.WriteFile(p, data, 0700)
}

// GetFilename gets the filename of the secrets file. Defaults to ~/.cmdotp/secrets
func GetFilename() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	p, err := filepath.Abs(usr.HomeDir)
	if err != nil {
		return "", err
	}
	p = filepath.Join(p, ".cmdotp", "secrets")
	return p, nil
}

// Exists returns true iff the file exists.
func Exists() (bool, error) {
	p, err := GetFilename()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(p)

	return !os.IsNotExist(err), nil
}
