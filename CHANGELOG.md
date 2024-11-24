# GoWebCore Changelog
Все изменения библиотеки GoWebCore будут документироваться на этой странице.

## 2024-11-24
### Added
- Добавлены:
  - Интерфейсы `mrsender.MailProvider`, `mrsender.MessageProvider`; 
  - `mail.Message`, `smtp.MailClient` для формирования и отправки электронных писем;
  - `telegrambot.MessageClient` для отправки сообщений в `telegram`;
- Добавлена ошибка `mrcore.ErrUseCaseRequiredDataIsEmpty`;

## 2024-11-17
### Added
- Добавлен `decorator.SourceEmitter` для добавления источника данных к уже имеющемуся `mrsender.EventEmitter`;

### Removed
- Удалён метод `mrsender.EmitWithSource`, вместо него нужно использовать `decorator.SourceEmitter`;

## 2024-11-16
### Added
- Добавлен процесс `signal.Interception` предназначенный для перехвата системных событий.
  При удалён `mrrun.AddSignalHandler`, выполнявший ранее эту функцию;
- Добавлен процесс `onstartup.Process` для выполнения работы в момент старта приложения;
- В интерфейс логгера добавлен метод `Int64()`;
- Добавлена структура `mrrun.StartingProcess` - для передачи информации о процессе находящемся в момент запуска;
- В `mrrun.Process` добавлен метод ReadyTimeout(), чтобы контролировать подвисший в момент запуска процесс;
- Добавлены константы для атрибутов: `mrapp.KeyProcessID`, `mrapp.KeyRequestID` и др.;

### Changed
- Изменена логика управления ошибками с учётом обновлений библиотеки `go-sysmess`.
  При этом был удалён `mrinit.ErrorManager` и добавлена вспомогательная структура `mrinit.ErrorSettings`
  и функции `AllEnabled()`, `WithCaller()`, `WithOnCreated()`, `AllDisabled()`, `CreateErrorOptionsMap()`
  которые помогают переопределять опции у конкретных Proto ошибок;
- Переименовано `ErrInternalProcessIsStoppedByTimeout` на `ErrInternalTimeoutPeriodHasExpired`;
- Теперь в реквестах внешний параметр `CorrelationID` выводится в логгере отдельным атрибутом;
- Доработан `MessageProcessor`, добавлены таймауты и отслеживание непредвиденного завершения работы воркеров;

## 2024-11-05
### Fixed
- Поправлен `README.md`;

## 2024-10-27
### Added
- Добавлен многопоточный сервис обработки сообщений `MessageProcessor` на основе консьюмера и обработчика;
- Добавлены функции `mrapp.WithProcessContext` и `ProcessCtx` для сохранения ID процесса в контексте;
- Добавлен `StorageErrorWrapper` для оборачивания ошибок инфраструктурного слоя приложения;
- Добавлены следующие типы ошибок:
    - `mrcore.ErrInternalCaughtPanic`;
    - `mrcore.ErrInternalProcessIsStoppedByTimeout`;
    - `mrcore.ErrInternalUnexpectedEOF`;
- Добавлены `FileOption` и `ImageOption` для более удобного задания опций у `mrparser.File` и `mrparser.Image`;
- Добавлена функция `mrcore.CastSliceToAnySlice`;

### Changed
- В функцию `CastToAppError` добавлен параметр `defFunc` для возможности возвращения ошибки по умолчанию;
- Переименованы:
    - `MimeTypeList.NewListByExts` -> `MimeTypesByExts`;
    - `Scheduler` -> `TaskScheduler`;
    - `TaskShell` -> `JobWrapper`;

## 2024-10-11
### Added
- Подключён и настроен линтер `gci`;

### Changed
- Тип некоторых полей исправлен на `uint64`; 

## 2024-10-09
### Changed
- Все проверки соответствия структур интерфейсам перенесены в тесты;
- Проведена профилактика по глобальным переменным, включен линтер `gochecknoglobals`;
- В пакете `mrserver/mrreq` в большинстве функций заменён входящий параметр `r *http.Request` на интерфейс `valueGetter`;
- Переименованы некоторые функции:
    - `BoolToInt64` -> `CastBoolToNumber`;
    - `BoolToPointer` -> `CastBoolToPointer`;
    - `Int32ToPointer` -> `CastNumberToPointer`;
    - `StringToPointer` -> `CastTimeToPointer`;
    - `TimeToPointer` -> `CastTimeToPointer`;
    - `TimePointerCopy` -> `CopyTimePointer`;
    - `ManagedError` -> `EnrichedError`;
- В ошибку `mrcore.ErrHttpRequestParseData` добавлен параметр для возможности задания подробностей ошибки;
- Переработан интерфейс для `ErrorHandler`, добавлены ему методы `Perform()` и PerformWithCommit();
- Доработан `EventEmitter`, добавлены метрики для подсчёта событий поступающих из разных источников;

## 2024-09-29
### Changed
- Добавлена директория `.cache` в `.gitignore`, скорректирована документация;

## 2024-09-28
### Added
- Добавлена ошибка `mrcore.ErrStorageConnectionIsBusy`;
- В `mrresp.SystemInfoConfig` добавлен массив `Processes`, который заполняется запущенными процессами приложения с их статусами;
- Добавлены объекты `mrrun.AppHealth`, `HealthProbe` для создания обработчиков проверяющие работоспособность системы;
- Добавлен пакет `mrtests/helpers` используемый в тестировании системы;

### Changed
- Доработан `AppRunner`, добавлена функция `onStartup` в метод `AppRunner.Run`,
  которая запускается только после запуска всех процессов;
- Переименовано `Usecase` -> `UseCase`;
- Доработан `mrprometheus.ObserveRequest`, добавлен дополнительный label: `location`;
- Поправлен `.editorconfig`, добавлены `*.proto`, `*.mk`;
- Удалена нулевая ёмкость у map;

## 2024-09-15
### Added
- Добавлен пакет `mrrun` для централизованного запуска процессов (сервисов);
- Добавлены комментарии к методам пакета `mrparser`;
- Добавлен логгер к `mrlib.MimeTypeList`, `mrparser.File`, `mrparser.Image`;
- Добавлены короткие команды в Makefile, обновлена инструкция команд.

### Changed
- Переименован `mrhttpserver` -> `mrhttp` и доработан адаптер http сервера в этом пакете;
- Доработан `mrschedule.Scheduler`;

### Removed
- Удалены методы `PrepareToStart` (вместо этого используется пакет `mrrun`);

## 2024-09-10
### Changed
- Доработан адаптер http сервера и перенесён в пакет `mrhttpserver`;
- Переименовано `mrserver.PrepareAppToStart` -> `mrapp.PrepareToStart`;

## 2024-09-08
### Added
- Добавлена функция `mrlog.HasCtx()` для того чтобы определить содержится ли логгер в указанном контексте.
- При возвращении ошибки `mrcore.ErrHttpRequestParseData` назначается её статус `422`;
- В адаптер `zerolog` добавлена опция `Options.Stdout` для возможности переопределения потока вывода.

### Changed
- Библиотека `mrcrypt` сделана в виде объекта с возможностью задания логгера.
- Поправлены `.env` переменные под новую версию `mrcmd`;

### Removed
- Удалена функция установки логгера по умолчанию `mrlog.SetDefault()`,
  необходимо всегда явно устанавливать логгер через объект или передавать его в контексте.

## 2024-08-03
### Added
- Добавлен интерфейс парсера `mrserver.RequestParserFloat64` и его реализация;
- Добавлены комментарии к некоторым сущностям;
- Добавлены функции:
    - `mrlib.CutBefore()`;
    - `mrlib.EqualFloat()`;

### Changed
- Изменены название метрик в `NewObserveRequest`, добавлен `namespace` в качестве параметра;
- При валидации `mrplayvalidator.Validate()` формируется более полное имя поля, где произошла ошибка;

### Fixed
- Доработана заглушка логгера `noplog.LoggerAdapter`;

## 2024-07-20
### Changed
- Изменены название метрик в `NewObserveRequest`, добавлен `namespace` в качестве параметра;

## 2024-07-14
### Changed
- В `MimeType` расширение можно указывать как с точкой так и без;
- Уточнено определение версии приложения;

## 2024-07-12
### Added
- Добавлена валидация `mrview.ValidateTripleSize`;
- Добавлена заглушка логгера `mrlog.noplog`; 
- Добавлены `mrlib.RoundFloat`, `RoundFloat2` `RoundFloat4` `RoundFloat8` для обнуления незначимых знаков после запятой;

## 2024-07-06
### Changed
- Строковые значения `LogLevel` приведены к единому стандарту `enum` (теперь они в верхнем регистре);
- Обновлен `github.workflows`;
- Обновлена документация по командам используемых при разработке;

## 2024-06-30
### Added
- Добавлена функция `Version()` для автоматического определения версии приложения; 
- Добавлено описание по локальной разработке библиотеки;
- Добавлен `ValidateDoubleSize`;

### Fixed
- Исправлено неправильное использование `isDebug`;

## 2024-06-24
### Added
- Добавлены `prometheus` метрики `ObserveRequest` для сбора статистики http запросов;

### Changed
- Заменено: `IP2intMust` -> `IP2int`;

### Removed
- Удалена поддержка соединения http сервера по сокету,
  также удалены `ListenTypeSock`, `ListenTypePort`;

## 2024-06-22
### Added
- Добавлен планировщик задач `mrworker`.`mrschedule`;
- Добавлена переменная `Environment` для задания рабочего окружения;
- Добавлен адаптер `mrsentry.Adapter` для отправки ошибок в `sentry`;

### Changed
- Изменён формат вывода ошибок в `MiddlewareRecoverHandler`;
- Теперь используется структура `ProtoExtra` для дополнительных настроек при создании ошибок;

## 2024-06-16
### Changed
- Настроен линтер `reviev` (`.golangci.yaml`);
- Добавлены комментарии к некоторым сущностям;

### Removed
- Удалена глобальная переменная `IsDebug`;

## 2024-06-15
### Changed
- Обновлена система формирования ошибок в связи с внедрением новой версии библиотеки `go-sysmess`:
    - изменён формат создания новых ошибок;
- Подключены линтеры с их настройками (`.golangci.yaml`);
- Добавлены комментарии для публичных объектов и методов;
- `MimeTypeList` теперь задаётся из `config.yaml`; 
- Сделаны следующие замены:
    - `FactoryErrHttp* -> FactoryErrHTTP*`;
    - `HandlerGetSystemInfoAsJson -> HandlerGetSystemInfoAsJSON`;
    - `HandlerGetStatusOKAsJson -> HandlerGetStatusOkAsJSON`;
    - `HandlerGetStructAsJson -> HandlerGetStructAsJSON`;
    - `HandlerGetStatInfoAsJson -> HandlerGetStatInfoAsJSON`;
    - `HandlerGetNotFoundAsJson -> HandlerGetNotFoundAsJSON`;
    - `HandlerGetMethodNotAllowedAsJson -> HandlerGetMethodNotAllowedAsJSON`;
    - `HandlerGetFatalErrorAsJson  HandlerGetFatalErrorAsJSON`;
    - `ParseUserIp -> ParseUserIP`;
    - `contentTypeJson -> contentTypeJSON`;
    - `contentTypeProblemJson -> contentTypeProblemJSON`;
- Обновлены многие пакеты: `mrperms`, `mrcrypt`, `mrresp`,
  `mrparser`, `mrlog`, `mridempotency`, `mrserver` и др.;

## 2024-04-10
### Added
- Добавлена поддержка роутера `chi`;
- Добавлен `MiddlewareRecoverHandler` для того чтобы логировать вызовы panic
  и возвращать клиенту результат в виде internal ошибки (`mrresp.HandlerGetFatalErrorAsJson`);
- Добавлена константа `mrserver.VarRestOfURL` для использования её в качестве
  параметра URL в которую роутер будет сохранять необработанную часть `URL` адреса;
- Добавлены `mrserver.StatRequest`, `mrresp.HandlerGetStatInfoAsJson` для фиксации и
  отображения статистики по запросам;

### Changed
- Приведено форматирование `var` к стандарту библиотеки;
- Добавлен флаг `mrzerolog.EventAdapter.isAutoCallerEnabled` чтобы `auto caller`
  срабатывал только при вызове `Err` метода, при этом ошибка не должна содержать свой `CallStack`;
- Изменён формат параметров в URL для роутера: `/v1/sample/:id -> /v1/sample/{id}`,
  такого вида параметры совместимы с `chi` роутером, а также их легче привести
  к виду URL используемым в `julienschmidt`;
- Переименовано:
  - `mrresponse` -> `mrresp`;
  - `mrjulienrouter.PathParam -> mrjulienrouter.URLPathParam`;

### Removed
- Удалён статус `ItemStatusRemoved`, теперь удаление контролируется через отдельное поле `deleted_at`;

## 2024-03-23
### Added
- Добавлена ошибка `FactoryErrUseCaseEntityNotAvailable`;

### Changed
- Рефакторинг сборки компонентов системы (factory):
    - добавлены `mrfactory.PrepareEachController()`, `mrfactory.PrepareController()`;
    - стандартизованы `mrfactory.WithPermission()` и `mrfactory.WithMiddlewareCheckAccess()`;
    - заменены `mrserver.HttpMiddleware`, `HttpMiddlewareFunc`, `HttpHandlerAdapterFunc`
      на нативные варианты;
    - удален метод `mrserver.HttpRouter.HttpHandlerFunc()`;
    - `mrrscors.New() + mrrscors.CorsAdapter.Middleware() -> mrrscors.Middleware()`;
    - стандартизованы:
        - `NewMiddlewareHandlerAdapter() -> MiddlewareHandlerAdapter()`;
        - `mrserver.MiddlewareCheckAccess`;
        - `mrserver.MiddlewareIdempotency`;

### Fixed
- Поправлена передача correlationID в контекст логгера;
- Исправлен баг в `mrparser.Int64.PathParamInt64`, `mrreq.ParseInt64` заменён на `strconv.ParseInt`;
- Исправлен баг в `mrserver.MiddlewareCheckAccess` - проверка доступа к обработчику; 

## 2024-03-20
### Changed
- Обновлены вспомогательные функции для перевода значений в указатели,
  а также добавлен `required` параметр:
    - `mrtype.BoolPointer -> BoolToPointer`;
    - `StringPointer -> StringToPointer`;
    - `TimePointer -> TimeToPointer`;
    - добавлены: `Int32ToPointer`, `Int64ToPointer`;

## 2024-03-19
### Added
- Добавлена функция `ValidateRewriteName`;

### Changed
- Поправлено форматирование документации;

## 2024-03-18
### Changed
- Внедрена новая версия библиотеки `go-sysmess`, в связи с этим:
    - всем ошибкам с типами `ErrorKindInternal` и `ErrorKindSystem` конструктор заменён на `NewFactoryWithCaller()`;
    - константа `ErrorKindInternalNotice` переименована в `ErrorKindInternal`;
    - централизовано логируются ошибки, у которых `appError.HasCallStack() == true`;
- В `mrzerolog.CallerWithSkipFrame` добавлена логика отключающая автоматический вызов `CallStack()`;
- Поправлены `skipFrame` в некоторых вызовах `Caller()`;

## 2024-03-17
### Changed
- Для `mrparser.ItemStatus` добавлена возможность указания значения по умолчанию,
  которое устанавливается, если значение не было передано из вне;
- Также функция `ParseItemStatusList` была оптимизирована и стала частью `mrparser.ItemStatus`;

## 2024-03-15
### Fixed
- Исправления в обработке фильтра `ParseRangeInt64`;

## 2024-03-14
### Added
- Добавлены стандартные обработчики ошибок приложения (404 и 405);
    - `HandlerGetMethodNotAllowedAsJson()`;
    - `HandlerGetNotFoundAsJson()`;
- Добавлен обработчик для вывода системной информации:
    - `HandlerGetSystemInfoAsJson()`;
- Добавлена проверка общего размера загружаемых файлов;
    - добавлена ошибка `FactoryErrHttpRequestFileTotalSizeMax`;
- Добавлена обработка нескольких загруженных файлов `mrreq.FormFiles`;
- Добавлена ошибка `FactoryErrInternalFailedToOpen`;
- Добавлен агрегатор парсеров `mrparser.Parser`;
- В `mrserver.RequestParserValidate` добавлен метод `ValidateContent`;
- Добавлены `mrtype.FileHeader` и `mrtype.ImageHeader`;
- В `mrminio.DownloadFile` добавлена проверка на успешную загрузку файла;
- Добавление обработки `http.ErrMissingBoundary` для отображения пользовательской ошибки,
  что файл не был загружен на сервер;
- Добавлена функциональность связанная с идемпотентностью операций:
    - добавлены `mridempotency.Provider` и `mridempotency.Response`;
    - добавлен метод `ResponseSender.SendBytes`;
    - добавлен прокси `NewCacheableResponseWriter` для возможности кэширования ответов запросов;
    - добавлен `MiddlewareIdempotency`;

### Changed
- Рефакторинг:
    - переименование `FactoryErrService* -> FactoryErrUseCase*`, `errService* -> errUseCase*`;
    - переименование интерфейсов `*Service -> *UseCase`;
    - переименование `X-Correlation-ID -> X-Correlation-Id`;
    - переименовано (`lastModified`, `ModifiedAt`) -> `UpdatedAt`;
- Настройки `PageSizeMax` и `PageSizeDefault` вынесены в общие настройки модулей `ModulesSettings.General`;
- Парсер `SortPage` разделён на два: `ListSorter`, `ListPager`;
- `mrparser.FormFileContents -> mrparser.FormFiles`, теперь возвращается список загруженных файлов без их открытия;
- В `FileResponseSender` изменился интерфейс методов `SendFile` и `SendAttachmentFile`;
- Переработан механизм загрузки файлов;
- Изменился интерфейс `mrserver.RequestDecoder`, `*http.Request` заменён `content io.Reader` и добавлен контекст;
- `mrserver.NewStatResponseWriter` вынесен в отдельный файл;
- Доработан адаптер `mrzerolog.LoggerAdapter`, добавлена поддержка `Auto Caller`;
- В `mrzerolog.Caller` поправлен skip параметр;
- `Error().Caller().Err(err) -> Error().Err(err)`;

### Fixed
- При обработке файлов исправлено `WrapImageError` на `WrapFileError`;

## 2024-02-05
### Added
- Добавлена новая пользовательская ошибка `mrcore.FactoryErrHttpFileUpload`;

## 2024-02-01
### Fixed
- Добавлен забытый вызов `.Send()` при записи ошибок в лог в некоторых методах,
  без которого информация не записывалась в журнал лога;

## 2024-01-30
### Added
- В `mrlib` добавлены `CallEachFunc`, `CloseFunc`, `Close` для группового закрытия ресурсов;

### Changed
- Для многих методов добавлен параметр `ctx context.Context`;
- Логгер перенесён из `mrcore` в `mrlog`, переработан его интерфейс и добавлен адаптер для `zerolog`.
  Все текущие вызовы логгера исправлены под новый интерфейс (например заменено `.Info(` на `.Info().Msg(` и т.д.);
- `EventBox` перенесён из `mrcore` в `mrsender`, переименован в `EventEmitter` и доработан его интерфейс;
- `MiddlewareFirst` переименован в `MiddlewareGeneral` и доработан, добавлена трассировка запросов,
  для более подробной статистики добавлен декоратор `StatResponseWriter`;
- Из `mrcore` перенесены и доработаны интерфейсы связанные с разрешениями в `mrperms`;
- Переработан `ServerAdapter` добавлен `PrepareToStart` для более гибкого управления его запуском и остановкой;
- Переименовано:
    - `mrserver.HandlerAdapter -> NewMiddlewareHttpHandlerAdapter`;
    - `mrcore.Debug -> mrcore.IsDebug`;
    - `ServiceHelper -> UsecaseHelper` (для устранения неоднозначности слова `Service`);

## 2024-01-25
### Added
- Добавлено: 
  `mrtype.BoolPointerCopy()`;
  `mrtype.TimePointerCopy()`;
  `mrlib.MimeTypeByExt()`;
- Разработаны парсеры для получения файлов и изображений из `multipart` формы.
  Они реализуют следующие интерфейсы `mrserver.RequestParserFile`, `mrserver.RequestParserImage`.
  Также для них были добавлены:
    - функция `mrreq.FormFile()` для извлечения файла из `multipart` формы;
    - типы `mrtype.Image`,`mrtype.ImageContent` подобные `mrtype.File`;
    - функции для работы с изображениями: `mrlib.DecodeImageConfig()`, `mrlib.CheckImage()`, `DecodeImage()`;
    - новые типы пользовательских ошибок `FactoryErrHttpRequestFile*`, `FactoryErrHttpRequestImage*`;

### Changed
- Переименовано `ErrorAttributeNameByDefault -> ErrorAttributeIDByDefault`;
- Интерфейс `mrserver.RequestParser` и его реализации были разложены на следующие интерфейсы:
    - `mrserver.RequestParserString`;
    - `mrserver.RequestParserInt64`;
    - `mrserver.RequestParserBool`;
    - `mrserver.RequestParserDateTime`;

### Removed
- `mrserver.RequestParserPath` удалён вместо него используется `mrserver.RequestParserParamFunc`;

## 2024-01-22
### Added
- Добавлена поддержка google/uuid (используется парсером запросов);

### Changed
- Расформирован объект `ClientContext` и его одноименный интерфейс, в результате:
    - изменена сигнатура обработчиков с `func(c mrcore.ClientContext)` на `func(w http.ResponseWriter, r *http.Request) error`;
    - с помощью интерфейсов `RequestDecoder`, `ResponseEncoder` можно задавать различные форматы
      принимаемых и отправляемых данных (сейчас реализован только формат `JSON`);
    - запросы обрабатываются встраиваемыми в обработчики объектов `mrparser.*` через интерфейсы:
      `mrserver.RequestParserPath`, `RequestParser`, `RequestParserItemStatus`, `RequestParserKeyInt32`,
      `RequestParserSortPage`, `RequestParserUUID`, `RequestParserValidate`;
    - ответы отправляются встраиваемыми в обработчики объекты `mrresponse.*` через интерфейсы:
      `mrserver.ResponseSender`, `FileResponseSender`, `ErrorResponseSender`;
    - вместо метода `Validate(structRequest any)` используется объект `mrparser.Validator`;
- Произведены следующие замены:
    - `HttpController.AddHandlers -> Handlers() []HttpHandler`.
      Убрана зависимость контроллера от роутера и секции.
      Для установки стандартных разрешений добавлены следующие методы
      `mrfactory.WithPermission`, `mrfactory.WithMiddlewareCheckAccess`; 
    - `ModulesAccess -> AccessControl` и добавлен интерфейс `mrcore.AccessControl`;
    - `AccessObject -> AccessRights`;
    - `ClientSection` -> AppSection удалена зависимость от `AccessControl`;
    - `DefaultWrapErrorFunc -> DefaultHttpErrorOverrideFunc`;
    - `IsAuthorized -> IsGuestAccess` (с инверсией флага);
- Перенесены следующие библиотеки:
    - `rs/cors -> mrserver/mrrscors`;
    - `julienschmidt/httprouter -> mrserver/mrjulienrouter`;
    - `go-playground/validator -> mrview/mrplayvalidator`;
    - `mrreq -> mrserver/mrreq`;
- Загрузка ролей (`loadRoleConfig`) происходит через библиотеку `yaml`, удалена зависимость от `ilyakaznacheev/cleanenv`;
- Расформирован объект `ClientContext` и его одноименный интерфейс, в результате;
- При внедрении новой версии библиотеки `go-sysmess` было заменено:
    - `ErrorInternalID -> ErrorCodeInternal`;
    - `mrerr.FieldError -> CustomError`;
    - `mrerr.FieldErrorList -> CustomErrorList`;
    - `FactoryErrInternalWithData` было удалено, вместо неё используется `FactoryErrInternal.WithAttr(...)`;

### Removed
- Удалёно `FactoryErrHttpResponseSendData`;

## 2024-01-19
### Changed
- Обновлены зависимости библиотеки;

## 2024-01-18
### Added
- Добавлен тип `mrtype.FileContent`, отличается от `mrtype.File` тем, что содержит тело файла в виде байтов;
- Добавлены функции `mrreq.File` и `mrreq.FileContent` для получения файла из реквеста;

### Changed
- Тип `mrenum.ItemStatus` теперь считывает из БД и сохраняет в БД целочисленные значения;

## 2024-01-16
### Added
- Добавлена ошибка `FactoryErrServiceInvalidFile`;
- Добавлен новый тип `mrtype.ImageInfo` и доработан тип `mrtype.FileInfo`;
- Добавлены функции `mrlib.MimeType`, `mrtype.StringPointer`, `mrtype.TimePointer`;
- Добавлены системные обработчики `HandlerGetHealth`, `HandlerGetStatusOKAsJson`, `HandlerGetStructAsJson`;

### Changed
- Тип `mrtype.NullableBool` заменён на `*bool`;
- Переименовано:
    - `mrlib.MimeTypeByExt -> MimeTypeByFile`;
    - `ServiceHelper.IsNotFound -> IsNotFoundError`;
- Доработки для поддержки `mrerr.CallerOptions`;

### Removed
- Удалён метод `ServiceHelper.WrapErrorEntity`;

## 2023-12-13
### Added
- В `AppHelper.Close()` добавлено логирование при закрытии соединения;
- Добавлен `mrlib.NewLockerMutex()`, который возвращает объект реализующий интерфейс `mrcore.Locker`,
  блокировку осуществляет с помощью Mutex;

### Changed
- В конструктор `ClientSection` теперь передаётся `ClientSectionOptions`;
  
## 2023-12-11
### Added
- Добавлен `mrdebug.NewLockerStub()`, который возвращает объект-заглушку реализующий интерфейс `mrcore.Locker`;

## 2023-12-10
### Removed
- Удалены бесполезные вызовы Caller() в `mrtool.AppHelper`;

## 2023-12-09
### Changed
- В конструктор `Logger` теперь передаётся `LoggerOptions`, через который:
    - теперь можно управлять отображением стека вызовов;
    - можно передать `CallerEnabledFunc`, которая в зависимости от ошибки может автоматически отключать
      отображение линии вызова в методах Err и Warn (это полезно, если в самой ошибке эта информация уже содержится);
- Доработана логика копирования объектов в `ClientContext.WithContext`, `Logger.With`, `Logger.Caller`, `ServiceHelper.Caller`;

### Removed
- Удалён метод `Logger.DisableFileLine()`, теперь это управляется через `LoggerOptions`;

## 2023-12-06
### Added
- Добавлены следующие типы ошибок:
    - `mrcore.FactoryErrInternalNotice` - используется для оборачивания ошибок без активации callstack;
    - `mrcore.FactoryErrInternalWithData` - используется как контейнер с данными с активацией callstack;
    - `mrcore.FactoryErrWithData` - используется как контейнер с данными без активации callstack;
    - `mrcore.FactoryErrServiceEntityVersionInvalid` - сообщает, что версия объекта не валидна
      (при сохранении, с использованием механизма версий);
    - `mrcore.FactoryErrServiceOperationFailed` - используется при любом неуспешном запросе (`API`, `Storage`);

### Changed
- Изменился порядок параметров в `mrcore.FactoryErrHttpRequestParseParam`;
- Удалён параметр из `mrcore.FactoryErrServiceTemporarilyUnavailable`;
- Теперь поле `ErrorDetailsResponse.ErrorTraceID` отображается пользователю только тогда, когда связанное
  с ней ошибка записывается в лог;
- Переименован метод `mrcore.ClientErrorWrapperFunc -> mrcore.ClientWrapErrorFunc`;
- Переименован метод `mrcore.FactoryErrServiceEntitySwitchStatusImpossible -> FactoryErrServiceSwitchStatusRejected`,
  а также изменено в нём кол-во параметров;
- переработан `mrtool.ServiceHelper`:
    - добавлен метод `IsNotFound` - для определения, что указанная ошибка связана с тем, что запись не найдена;
    - добавлен метод `WrapErrorEntity` для обёртывания ошибок при запросах с возможностью указания ошибки,
      которая будет создана, если присутствует ошибка, о том, что запись не найдена;
    - добавлен метод `WrapErrorEntityNotFoundOrFailed` для обёртывания ошибок при запросах
      (например: при получении, сохранении, удалении записи);
    - добавлен метод `WrapErrorFailed` для обёртывания ошибок при запросах
      (например: при получении списка, при обращении к API, создании записи);
    - добавлен метод `WrapErrorEntityFailed` для обёртывания ошибок при запросах
      (например: при обращении к API для получения записи);
    - методы `WrapErrorEntityFetch`, `WrapErrorEntityUpdate`, `WrapErrorEntityDelete`,
      вместо них него следует использовать `WrapErrorEntityNotFoundOrFailed`;
    - метод `WrapErrorEntityInsert` удалён, вместо него следует использовать `WrapErrorEntityFailed`;

### Removed
- Удалены следующие неиспользуемые типы ошибок:
    - `mrcore.FactoryErrHttpResponseSystemTemporarilyUnableToProcess`;
    - `mrcore.FactoryErrInternalInvalidData`;
    - `mrcore.FactoryErrInternalMapValueNotFound`;
    - `mrcore.FactoryErrServiceEmptyInputData`;
    - `mrcore.FactoryErrServiceEntityNotCreated`;
    - `mrcore.FactoryErrServiceEntityNotStored`;
    - `mrcore.FactoryErrServiceEntityNotRemoved`;
    - `mrcore.FactoryErrStorageQueryDataContainer`;
    - `mrcore.FactoryErrStorageFetchedInvalidData`;

## 2023-12-04
### Added
- Добавлена глобальная переменная отладочного режима и два метода для управления им:
  `mrcore.SetDebug()` и `mrcore.Debug()`;
- В интерфейс логгера добавлен метод `Level()` которы удобно использовать совместно с `mrcore.Debug()`;
- Во включенном debug режиме в `ErrorDetailsResponse.Details` и в `ErrorAttribute.DebugInfo`
  добавляется полная информация об ошибке предназначенная для разработчиков и тестировщиков;

### Changed
- Переработан механизм отправки ответа в sendErrorResponse в ClientContext;
- Доработан механизм обработки ошибок, в результате:
    - обновлены и переименованы следующие сущности:
        - `mrserver.AppErrorListResponse -> ErrorListResponse`;
        - `mrserver.AppErrorAttribute -> ErrorAttribute`;
        - `mrserver.AppErrorAttributeNameSystem -> ErrorAttributeNameByDefault`;
        - `mrcore.FactoryErrServiceEntityNotUpdated -> FactoryErrServiceEntityNotStored`;
        - `mrcore.FactoryErrServiceIncorrectSwitchStatus -> FactoryErrServiceEntitySwitchStatusImpossible`,
          также изменился порядок и кол-во параметров;
        - `mrcore.FactoryErrInternalNoticeDataContainer -> FactoryErrStorageQueryDataContainer`
          (и перенесён в `errors_storage.go`);
        - `mrcore.FactoryErrInternalParseData` удалён;
        - `mrcore.FactoryErrServiceEntityVersionIsIncorrect` удалён, вместо него следует
          использовать -> `FactoryErrServiceEntityNotFound`;
    - Internal и System ошибки для пользователя отображались со стандартным сообщением
      "Внутренняя ошибка сервера".
      Теперь можно для каждой такой ошибки, в yaml файлах на разных языках, описать более
      подробную причину понятную пользователю с инструкциями, что делать в данной ситуации;
    - переработан `mrtool.ServiceHelper`, в котором:
        - добавлен метод `WrapErrorEntityInsert`;
        - обновлены и переименованы следующие методы:
            - `WrapErrorForSelect -> WrapErrorEntityFetch` + добавлен параметр `entityData`;
            - `WrapErrorForUpdate -> WrapErrorEntityUpdate` + добавлен параметр `entityData`;
            - `WrapErrorForRemove -> WrapErrorEntityDelete` + добавлен параметр `entityData`;
            - `ReturnErrorIfItemNotFound` удалён, вместо него следует использовать `WrapErrorEntityFetch`;
    - изменено формирование идентификаторов валидаторов, которые проверяют пользовательские поля:
        - `errValidation{Name} -> validator_err_{name}` (пример: `errValidationGte -> validator_err_gte`);
        - новые идентификаторы также участвуют для описания пользовательских ошибок (в yaml файлах на разных языках);
- В `mrserver.corsAdapter` вместо передачи параметра `Debug` теперь передаётся `Logger` системы,
  и в зависимости от его настроек решается, включать `Debug` режим или нет;

## 2023-11-26
### Added
- Добавлено логирования ошибок при парсинге `json` данных поступивших в запросе;
- Добавлена ошибка `FactoryErrHttpMultipartFormFile` используемая, когда возникли
  проблемы при загрузке файла на сервер;
- При описании пользовательских ошибок (в yaml файлах на разных языках), можно теперь
  использовать следующие переменные: `name`, `type`, `value`, `param`;
- Добавлена функция `mrreq.ParseDateTime`;
- Добавлены функции для отладки: `MultipartForm`, `MultipartFileHeader`;

## 2023-11-23
### Added
- Добавлено логирование при старте текущих `Cors.AllowedOrigins`;

### Fixed
- Исправлен парсинг параметров сортировки списка;

## 2023-11-20
### Added
- Добавлены новые виды ошибок:
    - `FactoryErrHttpClientUnauthorized`;
    - `FactoryErrHttpAccessForbidden`;
    - `FactoryErrServiceEntityVersionIsIncorrect`;
- В `RoleGroup` добавлен метод `IsAuthorized`;

### Changed
- Доработан `NewBuilderPath`, добавлена возможность указания собственного названия плейсхолдера для части пути;
- Переименовано:
    - `mrcore.ClientData -> mrcore.ClientContext`;
    - `ClientContext::ParseAndValidate -> ClientContext::Validate` (удалён `ClientContext::Parse`);
    - `mrcore.Log() -> mrcore.DefaultLogger()`;
    - `ModulesAccess::DefaultRole -> GuestRole`;
- Объединено: `MiddlewareFirst` + `MiddlewareAcceptLanguage -> MiddlewareFirst`, выделен объект `ClientTools`;
- В `ServiceHelper` доработаны `WrapErrorForUpdate` и `WrapErrorForRemove`, добавлен метод `WrapErrorForUpdateWithVersion`;
- Удалён интерфейс `RequestPath`, теперь значение из пути можно получать методом `ClientContext::ParamFromPath`;
- В `clientContext` преобразована функция преобразования ошибок `wrapErrorFunc` во внешнюю функцию,
  по умолчанию вызывается `DefaultErrorWrapperFunc`, но её можно переопределить на собственную;
- В некоторых местах оптимизирована конкатенация строк (`Sprintf` заменён на нативный "+");
- Обновлён `.editorconfig`;

### Fixed
- В методах пакета `mrcrypt`: `GenTokenBase64`, `GenTokenHex`, `GenTokenHexWithDelimiter`
  выдавалось больше символов, чем указывалось в параметре `length`;

### Removed
- Удалены функции `WithLogger`, `WithCorrelationID`, `WithLocale` из пакета `mrctx`;

## 2023-11-13
### Added
- Добавлены новые типы в пакет `mrtype`: `FileInfo`, `File`, `SortParams`, `PageParams`, `NullableBool`;
- Добавлены парсеры для новых типов, такие как: `ParseRequiredBool`, `ParseNullableBool`, `ParseSortParams`, и т.д.;
- В логгер добавлены новые методы Caller, Warning и DisableFileLine
  (последний для отключения вывода информации о местоположения вызова лога);
- В `ServiceHelper` добавлен метод `Caller`;
- Добавлен интерфейс `BuilderPath` и его реализация для построения путей к файлам и URL к ресурсам;
- Добавлен интерфейс `ListSorter` для проверки существования сортируемого поля
  и получения поля для сортировки по умолчанию;

### Changed
- В логгере изменён `callerSkip` с 3 на 4, для того чтобы в логах выводить путь
  к родительской функции, откуда этот лог был вызван;
- Переименованы некоторые переменные и функции (типа Id -> ID) в соответствии с code style языка go;
- Переименовано: `mrcore.DefaultLogger() -> mrcore.Log()`
- Часть простых типов переехала из библиотеки `mrstorage` пакета `mrentity` в пакет `mrtype` ядра;
- В `RequestPath` переименован метод `GetInt -> GetInt64`;
- Удалены зависимости пакета `mrenum` от пакета `mrcore`;
- Доработана отправка файла, добавлены заголовки `Content-Length`, `Content-Disposition`;
- Обновлены зависимости библиотеки;
- Все файлы библиотеки были пропущены через `gofmt`;

## 2023-11-01
### Added
- Добавлена новая ошибка `mrcore.FactoryErrHttpRequestParamMax`;
- Добавлены короткие функции логирования `mrcore.LogWarn`, `LogInfo`, `LogDebug`;
- Добавлен парсинг строковых параметров из запроса: `mrreq.ParseStr`, `mrreq.ParseStrList`;

### Changed
- Оптимизирована работа с некоторыми структурами данных;
- Обновлены зависимости библиотеки;
- Добавлен пакет `mrperms` для работы с пользовательскими разрешениями и привилегиями.
  Основные сущности: `ClientSection`, `ModulesAccess`, `RoleGroup`;
- Переименованы следующие функции:
    - `mrcore.FactoryErrHttpRequestParamLen -> mrcore.FactoryErrHttpRequestParamLenMax`;
    - `mrreq.CorrelationId -> ParseCorrelationId`;
    - `mrreq.Enum -> ParseEnum`;
    - `mrreq.EnumList -> ParseEnumList`;
    - `mrreq.Int64 -> ParseInt64`;
    - `mrreq.Int64List -> ParseInt64List`;
- Перенесен `mrenum.ItemStatus` из библиотеки `go-components`;

### Removed
- Удалена сущность `Platform`, реализации достаточно в рамках проекта;

## 2023-10-08
### Added
- В пакет `mrcrypt` добавлены функции `GenPassword` и `PasswordStrength`;

### Changed
- Обновлены зависимости библиотеки;
- Обработка ошибок приведена к более компактному виду;

## 2023-09-20
### Added
- Добавлены ошибки `FactoryErrHttpRequestParamEmpty` и `FactoryErrServiceEmptyInputData`;
- Добавлен парсинг `mrreq.Int64`;

### Changed
- Переименовано `FactoryErrServiceEntityTemporarilyUnavailable -> FactoryErrServiceTemporarilyUnavailable`;
- Заменены tabs на пробелы в коде;
- При парсинге `mrreq.Enum`, если он пустой, возвращается ошибка;

## 2023-09-16
### Added
- В `mrcore.ClientData` добавился метод `SendFile`;
- Добавлен интерфейс `Locker` для захвата общих ресурсов в разделяемой памяти сервисов;
- Добавлена константа `LockerDefaultExpiry`;

### Changed
- Сообщение в логгере теперь формируется с помощью `strings.Builder`;
- Для некоторых структур добавлены автоматические проверки на реализацию ими необходимых интерфейсов; 

### Fixed
- Установка `Content-Type` теперь происходит непосредственно при отправке данных; 

## 2023-09-13
### Added
- Добавлены пакеты `mrserver` и `mrview`;
- Добавлен пример работы с `logger`;

### Changed
- Хелперы перенесены в пакет `mrtool` из `mrcore`;

## 2023-09-12
### Added
- Добавлен метод `mrcore.LogErr`;
- Добавлено общее описание библиотеки;

### Changed
- пакет `mrenv` разделён на `mrctx` и `mrreq` и соответственно удалены из названий
  постфиксы `*FromContext`, `*FromRequest`;
- Обновлены зависимости библиотеки;
- Обновлены списки часто используемых ошибок, некоторые названия ошибок переименованы;

### Fixed
- Исправлен баг: package `mrenv -> package mrctx`

## 2023-09-11
### Added
- Добавлены парсеры данных поступающих из запросов (`enum`, `list` и т.д.);

### Changed
- `ExtractLogger -> LoggerFromContext`;
- `ExtractLocale -> LocaleFromContext`;
- Изменён интерфейс логгера;

### Fixed
- Формат глобальных `const`, `type`, `var` приведён к общему виду;

## 2023-09-10
### Changed
- Обновлены зависимости библиотеки;

## 2023-09-03
### Added
- Подключен валидатор `go-playground/validator` через адаптер;

### Changed
- Переименован `logger -> loggerAdapter`;