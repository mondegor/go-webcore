# GoWebCore Changelog
Все изменения библиотеки GoWebCore будут документироваться на этой странице.

## 2023-11-13
### Added
- Добавлены новые типы в пакет `mrtype`: `FileInfo`, `File`, `SortParams`, `PageParams`, `NullableBool`;
- Добавлены парсеры для новых типов, такие как: `ParseRequiredBool`, `ParseNullableBool`, `ParseSortParams`, и т.д.;
- В логгер добавлены новые методы Caller, Warning и DisableFileLine (последний для отключения вывода информации о местоположения вызова лога);
- В `ServiceHelper` добавлен метод `Caller`;
- Добавлен интерфейс `BuilderPath` и его реализация для построения путей к файлам и URL к ресурсам;
- Добавлен интерфейс `ListSorter` для проверки существования сортируемого поля и получения поля для сортировки по умолчанию;

### Changed
- В логгере изменён `callerSkip` с 3 на 4, для того чтобы в логах выводить путь к родительской функции, откуда этот лог был вызван;
- Переименованы некоторые переменные и функции (типа Id -> ID) в соответствии с code style языка go;
- Переименовано: `mrcore.DefaultLogger()` -> `mrcore.Log()`
- Часть простых типов переехала из библиотеки `mrstorage` пакета `mrentity` в пакет `mrtype` ядра;
- В `RequestPath` переименован метод `GetInt` -> `GetInt64`;
- Удалены зависимости пакета `mrenum` от пакета `mrcore`;
- Доработана отправка файла, добавлены заголовки `Content-Length`, `Content-Disposition`;

## 2023-11-01
### Added
- Добавлена новая ошибка mrcore.FactoryErrHttpRequestParamMax;
- Добавлены короткие функции логирования mrcore.LogWarn, LogInfo, LogDebug;
- Добавлен парсинг строковых параметров из запроса: mrreq.ParseStr, mrreq.ParseStrList;

### Changed
- Оптимизирована работа с некоторыми структурами данных;
- Обновлены зависимости библиотеки;
- Добавлен пакет mrperms для работы с пользовательскими разрешениями и привилегиями. Основные сущности: ClientSection, ModulesAccess, RoleGroup;
- Переименованы следующие функции:
    - mrcore.FactoryErrHttpRequestParamLen -> mrcore.FactoryErrHttpRequestParamLenMax;
    - mrreq.CorrelationId -> ParseCorrelationId;
    - mrreq.Enum -> ParseEnum;
    - mrreq.EnumList -> ParseEnumList;
    - mrreq.Int64 -> ParseInt64;
    - mrreq.Int64List -> ParseInt64List;
- Перенесен mrenum.ItemStatus из библиотеки go-components;

### Removed
- Удалена сущность Platform, реализации достаточно в рамках проекта;

## 2023-10-08
### Added
- В пакет mrcrypto добавлены функции GenPassword и PasswordStrength;

### Changed
- Обновлены зависимости библиотеки;
- Обработка ошибок приведена к более компактному виду;

## 2023-09-20
### Added
- Добавлены ошибки FactoryErrHttpRequestParamEmpty и FactoryErrServiceEmptyInputData;
- Добавлен парсинг mrreq.Int64;

### Changed
- Переименовано FactoryErrServiceEntityTemporarilyUnavailable -> FactoryErrServiceTemporarilyUnavailable;
- Заменены tabs на пробелы в коде;
- При парсинге mrreq.Enum, если он пустой, возвращается ошибка;

## 2023-09-16
### Added
- В mrcore.ClientData добавился метод SendFile;
- Добавлен интерфейс Locker для захвата общих ресурсов в разделяемой памяти сервисов;
- Добавлена константа LockerDefaultExpiry;

### Changed
- Сообщение в логгере теперь формируется с помощью strings.Builder;
- Для некоторых структур добавлены автоматические проверки на реализацию ими необходимых интерфейсов; 

### Fixed
- Установка Content-Type теперь происходит непосредственно при отправке данных; 

## 2023-09-13
### Added
- Добавлены пакеты mrserver и mrview;
- Добавлен пример работы с logger;

### Changed
- Хелперы перенесены в пакет mrtool из mrcore;

## 2023-09-12
### Added
- Добавлен метод mrcore.LogErr;
- Добавлено общее описание библиотеки;

### Changed
- пакет mrenv разделён на mrctx и mrreq и соответственно удалены из названий постфиксы *FromContext, *FromRequest;
- Обновлены зависимости библиотеки;
- Обновлены списки часто используемых ошибок, некоторые названия ошибок переименованы;

### Fixed
- Исправлен баг: package mrenv -> package mrctx

## 2023-09-11
### Added
- Добавлены парсеры данных поступающих из запросов (enum, list и т.д.);

### Changed
- ExtractLogger -> LoggerFromContext;
- ExtractLocale -> LocaleFromContext;
- Изменён интерфейс логгера;

### Fixed
- Формат глобальных const, type, var приведён к общему виду;

## 2023-09-10
### Changed
- Обновлены зависимости библиотеки;

## 2023-09-03
### Added
- Подключен валидатор go-playground/validator через адаптер;

### Changed
- Переименован logger -> loggerAdapter;