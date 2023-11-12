package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/mondegor/go-webcore/mrview"
)

type (
	User struct {
		Login string `validate:"required,min=3,max=32,login"`
	}
)

var (
	regexpLogin = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
)

func main() {
	validator := mrview.NewValidator()

	if err := validator.Register("login", ValidateLogin); err != nil {
		fmt.Println(err)
		return
	}

	user := User{Login: "valid-login"}

	if err := validator.Validate(context.Background(), &user); err != nil {
		fmt.Println(err)
		return
	}

	user2 := User{Login: "not-valid-login!"}

	if err := validator.Validate(context.Background(), &user2); err != nil {
		fmt.Println(err)
		return
	}
}

func ValidateLogin(value string) bool {
	return regexpLogin.MatchString(value)
}
