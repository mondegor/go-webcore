package user

import (
	"github.com/google/uuid"
)

type (
	// User - comment struct.
	User struct {
		id       uuid.UUID
		group    string
		langCode string
	}
)

// New - создаёт объект User.
func New(id uuid.UUID, group, langCode string) *User {
	return &User{
		id:       id,
		group:    group,
		langCode: langCode,
	}
}

// ID - comment method.
func (u *User) ID() uuid.UUID {
	return u.id
}

// Group - comment method.
func (u *User) Group() string {
	return u.group
}

// LangCode - comment method.
func (u *User) LangCode() string {
	return u.langCode
}
