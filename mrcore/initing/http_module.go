package initing

import (
	"fmt"

	"github.com/mondegor/go-core/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// HttpModule - описывает HTTP-модуль для создания и инициализации контроллеров.
	// Модуль - это логическая группа контроллеров с общими компонентами и разрешениями.
	HttpModule struct {
		// Caption - название модуля в свободной форме.
		Caption string

		// Permission - разрешение по умолчанию для всех контроллеров модуля.
		// Может быть переопределено на уровне контроллера.
		Permission string

		// InitSharedComponents - функция инициализации общих компонентов модуля.
		// Вызывается один раз перед созданием первого контроллера.
		// Позволяет создать и настроить ресурсы, которые будут доступны всем контроллерам модуля.
		InitSharedComponents func() (err error)

		// Controllers - список контроллеров, принадлежащих модулю.
		Controllers []HttpController
	}

	// HttpController - описывает HTTP-контроллер для создания и регистрации обработчиков.
	// Контроллер - это группа связанных HTTP-обработчиков (например: CRUD операции для сущности).
	HttpController struct {
		// Caption - название контроллера в свободной форме.
		Caption string

		// Permission - разрешение для всех обработчиков контроллера.
		// Если не указано, используется разрешение родительского модуля.
		Permission string

		// Create - функция создания экземпляра контроллера.
		Create func() (mrserver.HttpController, error)
	}
)

// CreateHttpControllers - создаёт и инициализирует все контроллеры для указанных модулей.
// Процесс инициализации:
//  1. Для каждого модуля вызывается InitSharedComponents (если определена);
//  2. Для каждого контроллера вызывается Create();
//  3. Если разрешение контроллера не указано, наследуется разрешение модуля;
//  4. К каждому обработчику применяются операции: WithPermission + operations.
func CreateHttpControllers(logger mrlog.Logger, modules []HttpModule, operations ...PrepareHandlerFunc) (list []mrserver.HttpController, err error) {
	for _, module := range modules {
		mrlog.Info(logger, "Create and init module", "module", module.Caption, "permission", module.Permission)

		if module.InitSharedComponents != nil {
			if err := module.InitSharedComponents(); err != nil {
				return nil, fmt.Errorf("init shared components: %w", err)
			}
		}

		for _, c := range module.Controllers {
			if c.Create == nil {
				return nil, fmt.Errorf("create controller for module '%s'", module.Caption)
			}

			controller, err := c.Create()
			if err != nil {
				return nil, err
			}

			if c.Caption != "" || c.Permission != "" {
				mrlog.Info(
					logger,
					"Create and init controller",
					"controller", c.Caption,
					"permission", c.Permission,
				)
			}

			// если разрешение контроллера не указано, то используется разрешение его модуля.
			if c.Permission == "" {
				c.Permission = module.Permission
			}

			list = append(
				list,
				PrepareHttpController(
					controller,
					append([]PrepareHandlerFunc{WithPermission(c.Permission)}, operations...)...,
				),
			)
		}
	}

	return list, nil
}
