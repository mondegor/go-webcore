# GoWebCore Changelog
Все изменения библиотеки GoWebCore будут документироваться на этой странице.

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