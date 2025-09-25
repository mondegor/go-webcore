package mraccess

import (
	"context"

	"github.com/google/uuid"
)

type (
	// Member - сущность с наделёнными ей правами доступа к системе.
	Member interface {
		ID() uuid.UUID
		Group() string
		LangCode() string
	}

	// MemberProvider - возвращает сущность с наделёнными ей правами доступа к системе.
	MemberProvider interface {
		MemberByToken(ctx context.Context, value string) (Member, error)
	}

	// Section - comment interface.
	Section interface {
		Name() string
		BuildPath(routePath string) string
		Privilege() string
	}
)
