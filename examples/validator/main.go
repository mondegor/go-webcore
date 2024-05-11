package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"
)

type (
	User struct {
		Login string `validate:"required,min=3,max=16,login"`
	}
)

var (
	regexpLogin = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-]+[a-zA-Z]$`)
)

func main() {
	validator := mrplayvalidator.New()

	if err := validator.Register("login", ValidateLogin); err != nil {
		fmt.Println(err)
		return
	}

	user1 := User{Login: "valid-login"}

	if err := validator.Validate(context.Background(), &user1); err != nil {
		fmt.Println(err)
	}

	user2 := User{Login: "not-valid-login!"}

	if err := validator.Validate(context.Background(), &user2); err != nil {
		fmt.Println(err)
	}

	user3 := User{Login: "really-long-login-len-24"}

	if err := validator.Validate(context.Background(), &user3); err != nil {
		fmt.Println(err)
	}
}

func ValidateLogin(value string) bool {
	return regexpLogin.MatchString(value)
}
