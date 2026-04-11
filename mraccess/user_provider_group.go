package mraccess

import (
	"context"
	"errors"
)

type (
	// TypedUserProvider - провайдер пользователей, отмеченный каким-либо типом (db, jwt, etc).
	// Используется для группировки провайдеров по типу токена.
	TypedUserProvider struct {
		Type  string
		Value UserProvider
	}

	// userProviderGroup - внутренняя реализация агрегатора
	// нескольких провайдеров пользователей разных типов.
	userProviderGroup struct {
		type2provider   map[string]UserProvider
		typeByTokenFunc func(token string) string
	}
)

// NewUserProviderGroup - создаёт объект UserProvider, объединяющий несколько провайдеров разных типов.
// Функция typeByTokenFunc определяет тип провайдера по токену авторизации.
func NewUserProviderGroup(providers []TypedUserProvider, typeByTokenFunc func(value string) string) UserProvider {
	type2provider := make(map[string]UserProvider, len(providers))

	for _, provider := range providers {
		type2provider[provider.Type] = provider.Value
	}

	return &userProviderGroup{
		type2provider:   type2provider,
		typeByTokenFunc: typeByTokenFunc,
	}
}

// UserByToken - возвращает пользователя по указанному токену авторизации
// с привязанными к нему привилегиями и разрешениями.
// Определяет тип провайдера через typeByTokenFunc и делегирует запрос соответствующему провайдеру.
func (co *userProviderGroup) UserByToken(ctx context.Context, value string) (User, error) {
	if value == "" {
		return nil, errors.New("userProviderGroup: token value is empty")
	}

	if tp := co.typeByTokenFunc(value); tp != "" {
		if provider, ok := co.type2provider[tp]; ok {
			return provider.UserByToken(ctx, value)
		}
	}

	return nil, errors.New("userProviderGroup: provider not found for token")
}
