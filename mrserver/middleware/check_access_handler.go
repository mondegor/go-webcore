package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mraccess"
	"github.com/mondegor/go-core/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"
)

// CheckAccessHandler - middleware для проверки доступа пользователя к конечному обработчику запроса.
//
// Логика работы:
//  1. Извлекает access token из запроса;
//  2. Получает данные пользователя через userProvider;
//  3. Проверяет привилегию (Privilege) и разрешение (Permission) пользователя;
//  4. Устанавливает внутренние заголовки: UserID/Group, SessionID, а также
//     X-Internal-Lang-Code и X-Internal-Time-Zone из профиля пользователя;
//  5. Вызывает следующий обработчик в цепочке.
func CheckAccessHandler(
	logger mrlog.Logger,
	action mraccess.Action,
	userProvider mraccess.UserProvider,
) func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
	return func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			logger.Debug(
				ctx,
				"CheckAccessHandler",
				"handler", action.Name,
				"privilege", action.Privilege,
				"permission", action.Permission,
			)

			accessToken := request.AccessToken(r)
			if accessToken == "" {
				return errors.ErrHttpClientUnauthorized
			}

			currentUser, err := userProvider.UserByToken(ctx, accessToken)
			if err != nil {
				// if errors.Is(err, errors.ErrAccessForbidden) {
				// 	return errors.ErrHttpAccessForbidden
				// }
				return err
			}

			logger.Debug(ctx, "current user", "userId", uuid.UUID(currentUser.ID()).String())

			if action.Privilege != mraccess.PrivilegePublic && !currentUser.Has(action.Privilege) {
				return errors.ErrHttpAccessForbidden
			}

			if !currentUser.Has(action.Permission) {
				return errors.ErrHttpAccessForbidden
			}

			r.Header.Set(mrserver.HeaderKeyInternalUserIDSlashGroup, uuid.UUID(currentUser.ID()).String()+"/"+currentUser.Group()) // userId/group

			sessionID := currentUser.SessionID()
			if sessionID == "" {
				logger.Error(ctx, "session id is empty", "userId", uuid.UUID(currentUser.ID()).String())
			}

			r.Header.Set(mrserver.HeaderKeyInternalSessionID, sessionID)

			// язык и часовой пояс из профиля кладутся во внутренние заголовки, а клиентские
			// Accept-Language и X-Accept-Time-Zone не трогаются: сервер не перекрывает клиента,
			// а добавляет источник более высокого приоритета
			//
			// ВНИМАНИЕ: язык и пояс запроса этим не фиксируются - выше внутренних заголовков
			// ParserLocale и ParserTimeZone ставят query-параметры (?lang, ?tz), которые
			// приходят от клиента и здесь не срезаются: разовое переопределение
			// остаётся за клиентом
			r.Header.Set(mrserver.HeaderKeyInternalLangCode, currentUser.LangCode())
			r.Header.Set(mrserver.HeaderKeyInternalTimeZone, currentUser.TimeZone())

			if err = next(w, r); err != nil {
				// if errors.Is(err, errors.ErrAccessForbidden) {
				// 	return errors.ErrHttpAccessForbidden
				// }
				// если ошибка обработчика не связана с доступом к ресурсу
				return err
			}

			return nil
		}
	}
}
