package parser

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-core/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// User - извлекает информацию о пользователе из HTTP-запроса.
	User struct {
		locations locationList
		logger    mrlog.Logger
	}

	// locationList - предоставляет доступ к предзагруженным часовым поясам приложения.
	// Если запрошенный пояс не зарегистрирован, то возвращается пояс
	// по умолчанию и ошибка. Исключение - nopLocationList: список, не заданный
	// приложением, ошибкой не считается, поэтому заглушка отдаёт time.UTC без неё.
	locationList interface {
		LocationByName(value string) (*time.Location, error)
	}
)

// NewUser - создаёт объект User.
// Если locations не указан, то часовой пояс запроса всегда определяется как time.UTC.
func NewUser(
	locations locationList,
	logger mrlog.Logger,
) *User {
	if locations == nil {
		locations = nopLocationList{}
	}

	return &User{
		locations: locations,
		logger:    logger,
	}
}

// UserID - извлекает ID пользователя из HTTP запроса.
func (p *User) UserID(r *http.Request) uuid.UUID {
	val, _, _ := strings.Cut(r.Header.Get(mrserver.HeaderKeyUserIDSlashGroup), "/") // отбрасывается группа пользователя
	if val == "" {
		return uuid.Nil
	}

	return p.userID(r.Context(), val)
}

// UserAndGroup - извлекает ID пользователя и группу из HTTP запроса.
func (p *User) UserAndGroup(r *http.Request) (userID uuid.UUID, group string) {
	val, group, ok := strings.Cut(r.Header.Get(mrserver.HeaderKeyUserIDSlashGroup), "/")
	if !ok && val == "" {
		return uuid.Nil, ""
	}

	if group == "" {
		p.logger.Warn(r.Context(), "user group is empty", "userID", val)
	}

	return p.userID(r.Context(), val), group
}

// SessionID - извлекает ID сессии пользователя из HTTP запроса.
func (p *User) SessionID(r *http.Request) string {
	return r.Header.Get(mrserver.HeaderKeySessionID)
}

// Location - возвращает часовой пояс пользователя из HTTP-запроса.
//
// Единственный источник - внутренний заголовок X-Internal-Time-Zone (IANA-имя), который
// CheckAccessHandler заполняет из профиля пользователя, а CheckAccessTokenHandler удаляет.
//
// Сам парсер источник значения не проверяет и доверяет заголовку как внутреннему,
// поэтому срезание внутренних заголовков на входе - предусловие, которое обеспечивает
// приложение. На маршрутах, проходящих через CheckAccessHandler, заголовок в любом случае
// приводится к согласованному состоянию: пояс из профиля записывается, а при пустом
// значении в профиле заголовок удаляется, - поэтому подставить его клиент не может.
// Предусловие остаётся значимым для маршрутов, которые эти middleware не проходят.
//
// Имя, не зарегистрированное в locationList, логируется и не применяется. Если заголовок
// отсутствует, пуст или содержит незарегистрированное имя - возвращается time.UTC.
// В отличие от локали, подбора ближайшего пояса нет. Исключение - парсер, созданный
// без locationList: он молча сводит любой пояс к time.UTC и в лог ничего не пишет.
//
// Набор доступных поясов ограничен списком, который приложение зарегистрировало при старте,
// поэтому обращения к базе часовых поясов в процессе обработки запроса не происходит.
func (p *User) Location(r *http.Request) *time.Location {
	if name := r.Header.Get(mrserver.HeaderKeyTimeZone); name != "" {
		loc, err := p.locations.LocationByName(name)
		if err == nil {
			return loc
		}

		p.logger.Warn(
			r.Context(),
			"TimeZone header is incorrect",
			"header", mrserver.HeaderKeyTimeZone,
			"timezone", name,
			"error", err,
		)
	}

	return time.UTC
}

func (p *User) userID(ctx context.Context, value string) uuid.UUID {
	userID, err := uuid.Parse(value)
	if err != nil {
		p.logger.Warn(ctx, "userID parse error", "userID", value, "error", err)

		return uuid.Nil
	}

	return userID
}

type (
	// nopLocationList - заглушка для случая, когда список часовых поясов не задан.
	// Приложение, которому часовые пояса не нужны, - штатная конфигурация, а не сбой,
	// поэтому заглушка молчит: ошибка не возвращается и в лог ничего не пишется,
	// а любой запрошенный пояс сводится к time.UTC.
	nopLocationList struct{}
)

// LocationByName - всегда возвращает time.UTC.
func (l nopLocationList) LocationByName(_ string) (*time.Location, error) {
	return time.UTC, nil
}
