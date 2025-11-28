package mraccess

import (
	"context"
)

type (
	// UserProvider - интерфейс источника пользователей.
	UserProvider interface {
		UserByToken(ctx context.Context, value string) (User, error)
	}

	oneUserProvider struct {
		user User
	}
)

// NewOneUserProvider - создаёт объект UserProvider.
func NewOneUserProvider(user User) UserProvider {
	return &oneUserProvider{
		user: user,
	}
}

// UserByToken - возвращает по указанному токену авторизации пользователя
// с привязанными к нему привилегиями и разрешениями.
func (p *oneUserProvider) UserByToken(_ context.Context, _ string) (User, error) {
	return p.user, nil
}
