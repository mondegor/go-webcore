package mrparser

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// User - comment struct.
	User struct {
		logger mrlog.Logger
	}
)

// NewUser - создаёт объект User.
func NewUser(logger mrlog.Logger) *User {
	return &User{
		logger: logger,
	}
}

// UserID - comment method.
func (p *User) UserID(r *http.Request) uuid.UUID {
	val, _, _ := strings.Cut(r.Header.Get(mrreq.HeaderKeyUserIDSlashGroup), "/") // отбрасывается группа пользователя
	if val == "" {
		return uuid.Nil
	}

	return p.userID(r.Context(), val)
}

// UserAndGroup - comment method.
func (p *User) UserAndGroup(r *http.Request) (userID uuid.UUID, group string) {
	val, group, ok := strings.Cut(r.Header.Get(mrreq.HeaderKeyUserIDSlashGroup), "/")
	if !ok && val == "" {
		return uuid.Nil, ""
	}

	if group == "" {
		p.logger.Warn(r.Context(), "user group is empty", "userID", val)
	}

	return p.userID(r.Context(), val), group
}

func (p *User) userID(ctx context.Context, value string) uuid.UUID {
	userID, err := uuid.Parse(value)
	if err != nil {
		p.logger.Warn(ctx, "userID parse error", "userID", value, "error", err)

		return uuid.Nil
	}

	return userID
}
