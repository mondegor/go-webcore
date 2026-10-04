# Описание GoWebCore v0.29.2
Этот репозиторий содержит описание библиотеки GoWebCore.

## Статус библиотеки
Библиотека находится в стадии разработки.

## Описание библиотеки
Библиотека с базовой функциональностью для разработки web сервисов.
Общие интерфейсы (логгер, контроль доступа, воркеры, идемпотентность, ошибки и т.д.)
находятся в библиотеке [`go-core`](https://github.com/mondegor/go-core),
а в GoWebCore располагаются их адаптеры и http-инфраструктура. В библиотеку входят:
- интерфейсы http роутера (`mrserver.HttpRouter`), локализатора (`mrcore.Localizer`)
  и валидатора (`mrview.Validator`), которые могут быть реализованы уже в конкретных проектах;
- адаптер стандартного http сервера (`mrserver/httpserver`);
- адаптеры http роутеров:
    - `go-chi/chi/v5` (`mrserver/mrchi`);
    - `julienschmidt/httprouter` (`mrserver/mrjulienrouter`);
- middleware: проверка доступа и токена доступа, идемпотентность, сбор статистики запросов,
  перехват паник, генерация идентификатора запроса;
- адаптер cors (`rs/cors`);
- адаптер валидатора (`go-playground/validator/v10`);
- адаптеры для отправки ошибок в `sentry` (`mrclient/sentry`, handler для `slog` в `mrlog/slog/sentry`);
- клиенты для отправки писем по SMTP (`mrclient/mail`) и сообщений в `telegram` (`mrclient/telegram`);
- реализация метрик http запросов в `mrserver/mrprometheus.ObserveRequest`
  и метрик пула соединений с БД в `mrstorage/mrprometheus.DBCollector`;
- формирование http ответов (`mrserver/mrresp`, `mrserver/mrjson`);
- парсеры для некоторых типов данных, которые поступают из http запросов
  (в т.ч. определение IP клиента с учётом `X-Real-Ip` и `X-Forwarded-For`);
- парсеры для работы с файлами и изображениями;
- функции инициализации http модулей и метрик (`mrcore/initing`, `mrcore/mrinit`);
- отладочные утилиты (`mrdebug`) и хелперы для тестов (`mrtests/helpers`);
- готовые дашборды `grafana` для метрик (`grafana-dashboards`) и примеры использования (`examples`).

## Подключение библиотеки
`go get -u github.com/mondegor/go-webcore@v0.29.2`

## Установка библиотеки для её локальной разработки
- Выбрать рабочую директорию, где должна быть расположена библиотека
- `mkdir go-webcore && cd go-webcore` // создать и перейти в директорию проекта
- `git clone git@github.com:mondegor/go-webcore.git .`
- `cp .env.dist .env`
- `mrcmd go-dev deps` // загрузка зависимостей проекта
- Для работы утилит `gofumpt`, `goimports`, `gci`, `golangci-lint` необходимо запустить
  `mrcmd go-dev install-tools`. По умолчанию они устанавливаются последних версий;
  чтобы закрепить версию, раскомментируйте переменную `GO_DEV_TOOLS_INSTALL_*` в `.env`.
  `mockgen` и `gotext` в go-dev по умолчанию выключены и в библиотеке не используются

### Консольные команды используемые при разработке библиотеки

> Перед запуском консольных скриптов библиотеки необходимо скачать и установить утилиту Mrcmd.\
> Инструкция по её установке находится [здесь](https://github.com/mondegor/mrcmd#readme)

- `mrcmd go-dev help` // выводит список всех доступных go-dev команд;
- `mrcmd go-dev generate` // генерирует go файлы через встроенный механизм go:generate;
- `mrcmd go-dev gofumpt-fix` // исправляет форматирование кода (`gofumpt -l -w -extra ./`);
- `mrcmd go-dev goimports-fix` // исправляет imports, если это требуется (`goimports -l -w -local ${GO_DEV_IMPORTS_LOCAL_PREFIXES}` для всех go файлов, кроме сгенерированных);
- `mrcmd go-dev gci-fix` // упорядочивает imports (`gci`);
- `mrcmd go-dev lint` // запускает линтеры для проверки кода (на основе `.golangci.yaml`);
- `mrcmd go-dev test` // запускает тесты библиотеки;
- `mrcmd go-dev test-report` // запускает тесты библиотеки с формированием отчёта о покрытии кода (`test-coverage-full.html`);
- `mrcmd plantuml build-all` // генерирует файлы изображений из `.puml` [подробнее](https://github.com/mondegor/mrcmd-plugins/blob/master/plantuml/README.md#%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%B0-%D1%81-%D0%B4%D0%BE%D0%BA%D1%83%D0%BC%D0%B5%D0%BD%D1%82%D0%B0%D1%86%D0%B8%D0%B5%D0%B9-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0-markdown--plantuml);

#### Короткий вариант выше приведённых команд (Makefile)
- `make deps` // аналог `mrcmd go-dev deps`
- `make deps-upgrade` // аналог `mrcmd go-dev get -u ./...` + `mrcmd go-dev tidy`
- `make generate` // аналог `mrcmd go-dev generate`
- `make lint` // аналог `mrcmd go-dev gofumpt-fix` + `goimports-fix` + `gci-fix` + `lint`
- `make test` // аналог `mrcmd go-dev test`
- `make test-report` // аналог `mrcmd go-dev test-report`
- `make plantuml` // аналог `mrcmd plantuml build-all`

> Чтобы расширить список команд, необходимо создать Makefile.mk и добавить
> туда дополнительные команды, все они будут добавлены в единый список команд make утилиты.
