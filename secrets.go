package main

import (
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

type Secret struct {
	Name   string
	Secret string
}

func (s Secret) GenCode() string {
	str, err := totp.GenerateCode(s.Secret, time.Now())
	if err != nil {
		panic(err)
	}
	return str
}

func (s Secret) ToString() string {
	return fmt.Sprintf("%s\t\t%s", s.Name, s.GenCode())
}
