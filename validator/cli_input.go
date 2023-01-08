package validator

import (
	"errors"
	"strconv"
	"strings"
)

func HostnameValidator(input string) error {
	if len(strings.Split(input, "//")) == 2 {
		return nil
	}
	return errors.New("please provide https:// or http://")
}

func IsNumber(error_message string) func(input string) error {
	return func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New(error_message)
		}
		return nil
	}
}

func IsNotEmpty(input_space string) func(input string) error {
	return func(input string) error {
		if len(input) != 0 {
			return nil
		}
		return errors.New("please provide " + input_space)
	}
}

func PasswordValidator(input string) error {
	if len(input) < 6 {
		return errors.New("password must have more than 6 characters")
	}
	return nil
}
