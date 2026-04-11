package mraccess

import (
	"context"
)

type (
	// UserProvider - предоставляет метод для получения пользователя по токену авторизации.
	// Возвращает ошибку, если пользователь не найден или токен недействителен.
	UserProvider interface {
		UserByToken(ctx context.Context, value string) (User, error)
	}

	// oneUserProvider - внутренняя реализация провайдера, возвращающего
	// всегда одного и того же пользователя независимо от токена.
	oneUserProvider struct {
		user User
	}
)

// NewOneUserProvider - создаёт объект UserProvider, который всегда возвращает одного и того же пользователя.
// Полезно для тестирования или конфигураций с единственным фиксированным пользователем.
func NewOneUserProvider(user User) UserProvider {
	return &oneUserProvider{
		user: user,
	}
}

// UserByToken - возвращает заранее заданного пользователя независимо от токена.
func (p *oneUserProvider) UserByToken(_ context.Context, _ string) (User, error) {
	return p.user, nil
}
