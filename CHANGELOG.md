# GoWebCore Changelog
Все изменения библиотеки GoWebCore будут документироваться на этой странице.

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
  - `mrlib.MimeTypeByExt` -> `MimeTypeByFile`;
  - `ServiceHelper.IsNotFound` -> `IsNotFoundError`;
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
- Переименован метод `mrcore.ClientErrorWrapperFunc` -> `mrcore.ClientWrapErrorFunc`;
- Переименован метод `mrcore.FactoryErrServiceEntitySwitchStatusImpossible` -> `FactoryErrServiceSwitchStatusRejected`,
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
        - `mrserver.AppErrorListResponse` -> `ErrorListResponse`;
        - `mrserver.AppErrorAttribute` -> `ErrorAttribute`;
        - `mrserver.AppErrorAttributeNameSystem` -> `ErrorAttributeNameByDefault`;
        - `mrcore.FactoryErrServiceEntityNotUpdated` -> `FactoryErrServiceEntityNotStored`;
        - `mrcore.FactoryErrServiceIncorrectSwitchStatus` -> `FactoryErrServiceEntitySwitchStatusImpossible`,
          также изменился порядок и кол-во параметров;
        - `mrcore.FactoryErrInternalNoticeDataContainer` -> `FactoryErrStorageQueryDataContainer`
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
            - `WrapErrorForSelect` -> `WrapErrorEntityFetch` + добавлен параметр `entityData`;
            - `WrapErrorForUpdate` -> `WrapErrorEntityUpdate` + добавлен параметр `entityData`;
            - `WrapErrorForRemove` -> `WrapErrorEntityDelete` + добавлен параметр `entityData`;
            - `ReturnErrorIfItemNotFound` удалён, вместо него следует использовать `WrapErrorEntityFetch`;
    - изменено формирование идентификаторов валидаторов, которые проверяют пользовательские поля:
        - `errValidation{Name}` -> `validator_err_{name}` (пример: `errValidationGte` -> `validator_err_gte`);
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
    - `mrcore.ClientData` -> `mrcore.ClientContext`;
    - `ClientContext::ParseAndValidate` -> `ClientContext::Validate` (удалён `ClientContext::Parse`);
    - `mrcore.Log()` -> `mrcore.DefaultLogger()`;
    - `ModulesAccess::DefaultRole` -> `GuestRole`;
- Объединено: `MiddlewareFirst` + `MiddlewareAcceptLanguage` -> `MiddlewareFirst`, выделен объект `ClientTools`;
- В `ServiceHelper` доработаны `WrapErrorForUpdate` и `WrapErrorForRemove`, добавлен метод `WrapErrorForUpdateWithVersion`;
- Удалён интерфейс `RequestPath`, теперь значение из пути можно получать методом `ClientContext::ParamFromPath`;
- В `clientContext` преобразована функция преобразования ошибок `wrapErrorFunc` во внешнюю функцию,
  по умолчанию вызывается `DefaultErrorWrapperFunc`, но её можно переопределить на собственную;
- В некоторых местах оптимизирована конкатенация строк (`Sprintf` заменён на нативный "+");
- Обновлён `.editorconfig`;

### Fixed
- В методах пакета `mrcrypto`: `GenTokenBase64`, `GenTokenHex`, `GenTokenHexWithDelimiter`
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
- Переименовано: `mrcore.DefaultLogger()` -> `mrcore.Log()`
- Часть простых типов переехала из библиотеки `mrstorage` пакета `mrentity` в пакет `mrtype` ядра;
- В `RequestPath` переименован метод `GetInt` -> `GetInt64`;
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
    - `mrcore.FactoryErrHttpRequestParamLen` -> `mrcore.FactoryErrHttpRequestParamLenMax`;
    - `mrreq.CorrelationId` -> `ParseCorrelationId`;
    - `mrreq.Enum` -> `ParseEnum`;
    - `mrreq.EnumList` -> `ParseEnumList`;
    - `mrreq.Int64` -> `ParseInt64`;
    - `mrreq.Int64List` -> `ParseInt64List`;
- Перенесен `mrenum.ItemStatus` из библиотеки `go-components`;

### Removed
- Удалена сущность `Platform`, реализации достаточно в рамках проекта;

## 2023-10-08
### Added
- В пакет `mrcrypto` добавлены функции `GenPassword` и `PasswordStrength`;

### Changed
- Обновлены зависимости библиотеки;
- Обработка ошибок приведена к более компактному виду;

## 2023-09-20
### Added
- Добавлены ошибки `FactoryErrHttpRequestParamEmpty` и `FactoryErrServiceEmptyInputData`;
- Добавлен парсинг `mrreq.Int64`;

### Changed
- Переименовано `FactoryErrServiceEntityTemporarilyUnavailable` -> `FactoryErrServiceTemporarilyUnavailable`;
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
- Исправлен баг: package `mrenv` -> `package mrctx`

## 2023-09-11
### Added
- Добавлены парсеры данных поступающих из запросов (`enum`, `list` и т.д.);

### Changed
- `ExtractLogger` -> `LoggerFromContext`;
- `ExtractLocale` -> `LocaleFromContext`;
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
- Переименован `logger` -> `loggerAdapter`;