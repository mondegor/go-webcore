package mraccess

type (
	// TODO: возможно не нужен интерфейс для сущности User.

	// User - представляет пользователя с привязанными к нему
	// привилегиями и разрешениями.
	User interface {
		ID() [16]byte
		Group() string
		LangCode() string
		RightsChecker
	}

	// entryUser - внутренняя реализация интерфейса User.
	entryUser struct {
		id       [16]byte
		group    string
		langCode string
		rights   RightsChecker
	}
)

// NewUser - создаёт объект User.
func NewUser(id [16]byte, group, langCode string, rights RightsGetter) User {
	return &entryUser{
		id:       id,
		group:    group,
		langCode: langCode,
		rights:   rights.Rights(group),
	}
}

// ID - возвращает ID пользователя.
func (u *entryUser) ID() [16]byte {
	return u.id
}

// Group - возвращает группу пользователя (список ролей).
func (u *entryUser) Group() string {
	return u.group
}

// LangCode - возвращает язык пользователя.
func (u *entryUser) LangCode() string {
	return u.langCode
}

// HasPrivilege - сообщает о наличии указанной привилегии.
func (u *entryUser) HasPrivilege(name string) bool {
	return u.rights.HasPrivilege(name)
}

// HasPermission - сообщает о наличии указанного разрешения.
func (u *entryUser) HasPermission(name string) bool {
	return u.rights.HasPermission(name)
}
