package main

import (
    "context"
    "fmt"
    "regexp"

    "github.com/mondegor/go-core/mrview"
)

var regexpLogin = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)

type User struct {
    Login string `validate:"required,min=3,max=32,login"`
}

func main() {
    validator := mrview.NewValidator()
    err := validator.Register("login", ValidateLogin)

    if err != nil {
        fmt.Println(err)
        return
    }

    user := User{Login: "valid-login"}
    err = validator.Validate(context.Background(), &user)

    if err != nil {
        fmt.Println(err)
        return
    }

    user2 := User{Login: "not-valid-login!"}
    err = validator.Validate(context.Background(), &user2)

    if err != nil {
        fmt.Println(err)
        return
    }
}

func ValidateLogin(value string) bool {
    return regexpLogin.MatchString(value)
}
