package main

import (
	"context"
	"regexp"

	"github.com/mondegor/go-webcore/mrlog"

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
	logger := mrlog.New(mrlog.DebugLevel).With().Str("example", "validator").Logger()
	ctx := mrlog.WithContext(context.Background(), logger)

	validator := mrplayvalidator.New()

	if err := validator.Register("login", ValidateLogin); err != nil {
		logger.Info().Err(err).Msg("critical error")
		return
	}

	user1 := User{Login: "valid-login"}

	if err := validator.Validate(ctx, &user1); err != nil {
		logger.Info().Err(err).Msg("USER MESSAGE")
	}

	user2 := User{Login: "not-valid-login!"}

	if err := validator.Validate(ctx, &user2); err != nil {
		logger.Info().Err(err).Msg("USER MESSAGE")
	}

	user3 := User{Login: "really-long-login-len-24"}

	if err := validator.Validate(ctx, &user3); err != nil {
		logger.Info().Err(err).Msg("USER MESSAGE")
	}
}

// ValidateLogin - comment func.
func ValidateLogin(value string) bool {
	return regexpLogin.MatchString(value)
}
