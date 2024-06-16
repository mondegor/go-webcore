package main

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/mrlogbase"
	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"
)

type (
	// User - comment struct.
	User struct {
		Login string `validate:"required,min=3,max=16,login"`
	}
)

var regexpLogin = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-]+[a-zA-Z]$`)

func main() {
	logger := mrlogbase.New(mrlog.DebugLevel).With().Str("example", "validator").Logger()
	if err := mrlog.SetDefault(logger); err != nil {
		log.Fatal(err)
	}

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

// ValidateLogin - comment func.
func ValidateLogin(value string) bool {
	return regexpLogin.MatchString(value)
}
