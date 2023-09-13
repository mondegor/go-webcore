# GoWebCore Changelog
Все изменения библиотеки GoWebCore будут документироваться на этой странице.

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