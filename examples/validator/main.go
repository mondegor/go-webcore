package main

import (
	"context"
	"os"
	"regexp"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrlog/slog"

	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"
)

type (
	// User - представляет пользователя с валидируемым логином.
	User struct {
		Login string `validate:"required,min=3,max=16,login"`
	}
)

var regexpLogin = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-]+[a-zA-Z]$`)

func main() {
	l, _ := slog.NewLoggerAdapter(
		slog.WithWriter(os.Stdout),
		slog.WithLevel("DEBUG"),
	)
	logger := mrlog.WithAttrs(l, "example", "validator")

	validator := mrplayvalidator.New(logger)
	ctx := context.Background()

	if err := validator.Register("login", ValidateLogin); err != nil {
		mrlog.Error(logger, "critical error", "error", err)

		return
	}

	user1 := User{Login: "valid-login"}

	if err := validator.Validate(ctx, &user1); err != nil {
		mrlog.Info(logger, "USER MESSAGE1", "error", err)
	}

	user2 := User{Login: "not-valid-login!"}

	if err := validator.Validate(ctx, &user2); err != nil {
		mrlog.Info(logger, "USER MESSAGE2", "error", err)
	}

	user3 := User{Login: "really-long-login-len-24"}

	if err := validator.Validate(ctx, &user3); err != nil {
		mrlog.Info(logger, "USER MESSAGE3", "error", err)
	}
}

// ValidateLogin - проверяет, что логин соответствует допустимому формату.
func ValidateLogin(value string) bool {
	return regexpLogin.MatchString(value)
}
