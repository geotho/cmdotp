package main

import (
	"encoding/base32"
	"fmt"
	"os"
	"strings"

	"github.com/bgentry/speakeasy"
	"github.com/brandur/simplebox"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "GoTOTP"
	app.Action = PrintCodes
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "password",
			Usage: "password used to encrypt secrets file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a secret",
			Action:  AddSecret,
			Flags:   app.Flags,
		},
	}

	app.Run(os.Args)
}

func AddSecret(c *cli.Context) error {
	if c.NArg() != 2 {
		return fmt.Errorf("AddSecret expects 2 args, %d given: %v", c.NArg(), c.Args())
	}

	args := c.Args()

	if _, err := base32.StdEncoding.DecodeString(args[1]); err != nil {
		return fmt.Errorf("Secret is not valid base32")
	}

	if err := passwordPrompt(c); err != nil {
		return fmt.Errorf("Could not read password: %s", err)
	}

	key := getPassword(c)
	secrets, err := Load(key)
	if err != nil {
		if err.Error() == "Ciphertext could not be decrypted." {
			return fmt.Errorf("Could not load secrets: please check your password.")
		}
		return fmt.Errorf("Could not load secrets: %v", err)
	}

	secrets = append(secrets, Secret{
		Name:   args[0],
		Secret: args[1],
	})

	err = Save(secrets, key)
	if err != nil {
		return fmt.Errorf("Could not save secrets: %v", err)
	}
	return nil
}

func PrintCodes(c *cli.Context) error {
	if err := passwordPrompt(c); err != nil {
		return fmt.Errorf("Could not read password: %s", err)
	}

	key := getPassword(c)

	secrets, err := Load(key)
	if err != nil {
		if err.Error() == "Ciphertext could not be decrypted." {
			return fmt.Errorf("Could not load secrets: please check your password.")
		}
		return fmt.Errorf("Could not load secrets: %v", err)
	}

	for _, s := range secrets {
		fmt.Println(s.ToString())
	}
	return nil
}

// passwordPrompt asks the user for a password and then sets it in the context.
func passwordPrompt(c *cli.Context) error {
	pass, err := speakeasy.Ask("Password:\n")
	if err != nil {
		return err
	}

	return c.Set("password", pass)
}

// getPassword extracts password from context.
func getPassword(c *cli.Context) *[simplebox.KeySize]byte {
	pw := c.String("password")

	if len(pw) < 32 {
		pw = pw + strings.Repeat("0", 32-len(pw))
	}

	var key [simplebox.KeySize]byte
	copy(key[:], pw[:32])

	return &key
}
